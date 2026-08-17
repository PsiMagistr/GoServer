package game

import (
	"fmt"
	"math"
	"sync"
	"time"

	"GoServer/internal/config"
	"GoServer/internal/database"
	"GoServer/internal/models"
)

// Структура для комнатных сообщений
type RoomMessage struct {
	WorldID    string
	LocationID string
	Payload    interface{}
}

type MoveData struct {
	DestinationID      string
	ArrivalTime        time.Time
	TargetWorldName    string
	TargetLocationName string
}

type BattleChallenge struct { // Заявка на бой
	SenderID   int64     `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	TargetID   int64     `json:"target_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	TimeLeft   int       `json:"time_left"`
}

type Hub struct {
	mu             sync.RWMutex
	Clients        map[int64]*Client
	movingPlayers  map[int64]*MoveData
	challenges     map[int64]map[int64]*BattleChallenge
	activeBattles  map[int64]*Battle
	playerToBattle map[int64]int64
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan interface{}
	RoomBroadcast  chan RoomMessage
	lastBattleID   int64 // Добавляем поле-счетчик
}

func NewHub() *Hub {
	return &Hub{
		Clients:        make(map[int64]*Client),
		movingPlayers:  make(map[int64]*MoveData),
		challenges:     make(map[int64]map[int64]*BattleChallenge),
		activeBattles:  make(map[int64]*Battle),
		playerToBattle: make(map[int64]int64),
		Register:       make(chan *Client, 64),
		Unregister:     make(chan *Client, 64),
		Broadcast:      make(chan interface{}, 256),
		RoomBroadcast:  make(chan RoomMessage, 256),
		lastBattleID:   time.Now().Unix(),
	}
}

func (h *Hub) Run() {
	regenTicker := time.NewTicker(5 * time.Second)
	cleanupTicker := time.NewTicker(1 * time.Minute)
	defer regenTicker.Stop()
	for {
		select {
		case client := <-h.Register: // Регистрация.
			h.handleRegister(client)
		case client := <-h.Unregister:
			h.handleUnregister(client)
		case globalMessage := <-h.Broadcast:
			h.BroadcastToAll(globalMessage)
		case roomMessage := <-h.RoomBroadcast:
			h.BroadcastToRoom(roomMessage.WorldID, roomMessage.LocationID, roomMessage.Payload)
		case <-regenTicker.C:
			h.handleRegeniration()
		case <-cleanupTicker.C:
			h.cleanupChallenges()
		}
	}
}

/*Снапшот боя*/
func (h *Hub) getBattleSnapshot(battleID int64, forPlayerID int64) *BattleSnapshot {
	b, ok := h.activeBattles[battleID]
	if !ok {
		return nil
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	var you, opponent models.Character
	if b.AttackerData.ID == forPlayerID {
		you = b.AttackerData
		opponent = b.DefenderData
	} else {
		you = b.DefenderData
		opponent = b.AttackerData
	}
	return &BattleSnapshot{
		BattleID: b.ID,
		Round:    b.Round,
		TimeLeft: int(math.Ceil(float64(time.Until(b.ExpiresAt).Seconds()))),
		You: BattleFighterDTO{
			ID:       you.ID,
			Name:     you.Name,
			Level:    you.Level,
			HP:       you.HP, // Из боя
			MaxHP:    you.MaxHP,
			Mana:     you.Mana, // Из боя
			MaxMana:  you.MaxMana,
			AvatarID: you.AvatarID,
			Gender:   you.Gender,
		},
		Opponent: BattleFighterDTO{
			ID:       opponent.ID,
			Name:     opponent.Name,
			Level:    opponent.Level,
			HP:       opponent.HP, // Из боя
			MaxHP:    opponent.MaxHP,
			Mana:     opponent.Mana, // Из боя
			MaxMana:  opponent.MaxMana,
			AvatarID: opponent.AvatarID,
			Gender:   opponent.Gender,
		},
		Log: b.Log,
	}
}

/*Подготовка пакета syncWorld*/
func (h *Hub) prepareSyncState(client *Client) map[string]interface{} {
	neighbors := h.getNeighbors(client.Character.WorldID, client.Character.LocationID)
	currentWorld := Universe[client.Character.WorldID]
	currentNode := currentWorld.Points[client.Character.LocationID]
	challenges := h.GetChallenges(client.Character.ID)
	var timeLeft int
	worldName := currentWorld.Name
	locationName := currentNode.Name
	moveInfo, isMoving := h.movingPlayers[client.Character.ID]
	if isMoving {
		client.Character.State = models.StatusMoving // Синхронизируем стейт
		secondsLeft := time.Until(moveInfo.ArrivalTime).Seconds()
		timeLeft = int(math.Ceil(secondsLeft))
		worldName = moveInfo.TargetWorldName
		locationName = moveInfo.TargetLocationName
	}
	battleID, hasNum := h.playerToBattle[client.Character.ID]
	var battleInfo *BattleSnapshot
	if hasNum {
		client.Character.State = models.StatusBattle
		battleInfo = h.getBattleSnapshot(battleID, client.Character.ID)

	}
	return map[string]interface{}{
		"type":          "world_sync",
		"player":        client.Character,
		"world":         currentWorld,
		"players":       neighbors,
		"challenges":    challenges,
		"worlds":        currentNode.Worlds,
		"is_moving":     isMoving,
		"duration":      timeLeft,
		"world_id":      client.Character.WorldID, // Например "main_city"
		"location_id":   client.Character.LocationID,
		"world_name":    worldName,
		"location_name": locationName,
		"battle_info":   battleInfo,
	}
}

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	// 1. Безопасность
	if oldClient, ok := h.Clients[client.Character.ID]; ok {
		h.Send(oldClient, map[string]interface{}{
			"type": "connection_replaced",
		})
		// Принудительно рвем старую трубу
		fmt.Printf("Персонаж %s вытеснен новой сессией.\n", client.Character.Name)
		oldClient.Conn.Close()
	}
	// 2. Регистрация
	h.Clients[client.Character.ID] = client
	// 3. Подготовка данных (Snapshot)
	syncData := h.prepareSyncState(client)
	h.mu.Unlock()
	// 4. Отправка (уже без блокировки всего сервера!)
	h.Send(client, syncData)
	// 5. Уведомление окружающих
	h.BroadcastToRoomExcept(client.Character.WorldID, client.Character.LocationID, client.Character.ID, map[string]interface{}{
		"type": "player_joined",
		"player": map[string]interface{}{
			"id":        client.Character.ID,
			"name":      client.Character.Name,
			"avatar_id": client.Character.AvatarID,
			"level":     client.Character.Level,
			"gender":    client.Character.Gender,
			"state":     client.Character.State,
		},
	})
	fmt.Printf("Персонаж %s онлайн.\n", client.Character.Name)
}

func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	currentInMap, ok := h.Clients[client.Character.ID]
	if ok && currentInMap == client {
		type Params = struct {
			charID      int64
			charName    string
			charWorldID string
			charLocID   string
			charHP      float64
			charMana    float64
		}
		param := Params{
			charID:      client.Character.ID,
			charWorldID: client.Character.WorldID,
			charLocID:   client.Character.LocationID,
			charName:    client.Character.Name,
			charHP:      client.Character.HP,
			charMana:    client.Character.Mana,
		}
		fmt.Printf("Персонаж %s не в сети. \n", param.charName)
		go func(p Params) {
			err := database.UpdateCharacterHpMana(p.charID, p.charHP, p.charMana)
			if err != nil {
				fmt.Printf("ОШИБКА сохранения персонажа %s (ID %d): %v", p.charName, p.charID, err)
			}
		}(param)
		delete(h.Clients, client.Character.ID)
		close(client.Send)
		h.mu.Unlock()
		h.BroadcastToRoom(param.charWorldID, param.charLocID, map[string]interface{}{
			"type": "player_left",
			"player": map[string]interface{}{
				"id":   param.charID,
				"name": param.charName,
			},
		})
	} else {
		// Если это "призрак" старой сессии, просто отпускаем замок и ничего не делаем
		h.mu.Unlock()
		// log.Println("Игнорируем попытку удаления устаревшей сессии")
	}
}

func (h *Hub) BroadcastToRoom(worldID string, locationID string, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if client.Character.WorldID == worldID && client.Character.LocationID == locationID {
			select {
			case client.Send <- message:
			default:
				client.Conn.Close()
			}
		}
	}
}

func (h *Hub) BroadcastToRoomExcept(worldID string, locationID string, exeptID int64, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if client.Character.WorldID == worldID && client.Character.LocationID == locationID && client.Character.ID != exeptID {
			select {
			case client.Send <- message:
			default:
				client.Conn.Close()
			}
		}
	}
}

func (h *Hub) BroadcastToAll(message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			client.Conn.Close()
		}
	}
}

func (h *Hub) BroadcastPrivateMessage(charID int64, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.Clients[charID]
	if !ok {
		return
	}
	select {
	case client.Send <- message:
	default:
		client.Conn.Close()
	}
}

func (h *Hub) Send(client *Client, message interface{}) {
	select {
	case client.Send <- message:
	default:
		client.Conn.Close()
	}
}

func (h *Hub) getNeighbors(worldID string, locationID string) []map[string]interface{} {
	neighbors := make([]map[string]interface{}, 0)
	for _, other := range h.Clients {
		if other.Character.WorldID == worldID && other.Character.LocationID == locationID {
			neighbors = append(neighbors, map[string]interface{}{
				"id":        other.Character.ID,
				"name":      other.Character.Name,
				"avatar_id": other.Character.AvatarID,
				"level":     other.Character.Level,
				"gender":    other.Character.Gender,
				"state":     other.Character.State,
			})
		}
	}
	return neighbors
}

func (h *Hub) GetClientByName(name string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if client.Character.Name == name {
			return client
		}
	}
	return nil
}

// Регенерация жизни и маны.
func (h *Hub) handleRegeniration() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if h.GetFullStatus(client.Character.ID) != models.StatusFree {
			// fmt.Println("Вы не можете восстанавливаться")
			continue
		}
		hpChanged := client.AddHP(2)
		mpChanged := client.AddMana(5)
		if hpChanged || mpChanged {
			h.Send(client, map[string]interface{}{
				"type":   "player_update",
				"player": client.Character,
			})
		}
	}
}

func (h *Hub) SystemMsg(c *Client, text string) {
	h.Send(c, map[string]interface{}{
		"type": "sys_msg",
		"text": text,
	})
}

func (h *Hub) BattleMsg(c *Client, text string) {
	h.Send(c, map[string]interface{}{
		"type": "battle_log",
		"text": text,
	})
}

func (h *Hub) GetFullStatus(charID int64) models.PlayerStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, moving := h.movingPlayers[charID]
	if moving {
		return models.StatusMoving
	}
	client, ok := h.Clients[charID]
	if ok {
		return client.Character.State
	}
	return models.StatusFree
}

func (h *Hub) GetActiveClient(charID int64) (*Client, bool) { // Не использовать в местах где мьютекс уже взят.
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.Clients[charID]
	return client, ok
}

func (h *Hub) GetChallenges(RecipientID int64) []*BattleChallenge {
	var myChallenges []*BattleChallenge
	pending, exists := h.challenges[RecipientID]
	if exists {
		for _, challenge := range pending {
			duration := time.Until(challenge.ExpiresAt)
			if duration > 0 {
				challenge.TimeLeft = int(math.Ceil(duration.Seconds()))
				myChallenges = append(myChallenges, challenge)
			}
		}
	}
	return myChallenges
}

func (h *Hub) cleanupChallenges() {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now()
	count := 0
	for recipientID, invites := range h.challenges {
		for senderID, challenge := range invites {
			if now.After(challenge.ExpiresAt) {
				delete(invites, senderID)
				count++
			}
			if len(invites) == 0 {
				delete(h.challenges, recipientID)
			}
		}
	}
	if count > 0 {
		fmt.Printf("[CLEANUP] Удалено просроченных заявок: %d\n", count)
	}
}

func (h *Hub) GetInviteFromSpecificPlayer(recipientID int64, senderID int64) *BattleChallenge {
	invites, ok := h.challenges[recipientID]
	if ok {
		challange, find := invites[senderID]
		if find {
			return challange
		}

	}
	return nil
}

func (h *Hub) validateBattleTurn(c *Client, selectedIDs []int) error {
	if len(selectedIDs) != 5 {
		return fmt.Errorf("Нужно выбрать ровно 5 заклинаний")
	}
	attackCount := 0
	shieldCount := 0
	var totalMana float64
	for _, id := range selectedIDs {
		spell, exists := database.AllSpells[id]
		if !exists {
			return fmt.Errorf("заклинание ID:%d не существует", id)
		}
		hasIt := false
		for _, owned := range c.Character.Spells {
			if owned.ID == id {
				hasIt = true
				break
			}
		}
		if !hasIt {
			return fmt.Errorf("вы не владеете магией '%s'", spell.Name)
		}
		if spell.Type == "attack" {
			attackCount++
		} else {
			shieldCount++
		}
		totalMana += spell.ManaCost

	}
	if shieldCount != 2 || attackCount != 3 {
		return fmt.Errorf("комбинация должна быть: 2 щита и 3 атаки")
	}
	// 4. Проверка ресурсов
	if totalMana > c.Character.Mana {
		return fmt.Errorf("недостаточно маны (нужно %.1f, есть %.1f)", totalMana, c.Character.Mana)
	}
	return nil
}

// Метод надзирателя
func (h *Hub) battleTimerGuard(battleID int64, round int) {
	h.mu.RLock()
	b, exists := h.activeBattles[battleID]
	h.mu.RUnlock()
	if !exists {
		return
	}
	time.Sleep(time.Until(b.ExpiresAt.Add(2 * time.Second)))
	b.mu.Lock()
	if b.Finished || b.Round != round {
		b.mu.Unlock()
		return
	}
	fmt.Println("Гуард выход по таймауту")
	b.mu.Unlock()
	h.resolveBattleRound(b, true)
}

func (h *Hub) resolveBattleRound(b *Battle, isTimeout bool) {
	fmt.Println("resolveBattleRound")
	b.mu.Lock()
	if b.Finished {
		fmt.Println("Финишировано resolveBattleRound")
		b.mu.Unlock()
		return
	}
	p1Moved := b.AttackerTurn != nil
	p2Moved := b.DefenderTurn != nil

	if isTimeout {
		fmt.Println("Выход по таймауту")
		if !p1Moved && !p2Moved {
			fmt.Println("Оба бойца пропустили ход")
			b.mu.Unlock() // Отпускаем ПЕРЕД вызовом finishBattle
			h.finishBattle(b, "Ничья: оба бойца пропустили ход.")
			return
		}
		if !p1Moved {
			b.mu.Unlock()
			h.finishBattle(b, fmt.Sprintf("Победа %s: оппонент не явился.", b.DefenderData.Name))
			return
		}
		if !p2Moved {
			b.mu.Unlock()
			h.finishBattle(b, fmt.Sprintf("Победа %s: оппонент не явился.", b.AttackerData.Name))
			return
		}
	}

	// ЛОГИКА НОВОГО РАУНДА
	// b.Log = append(b.Log, fmt.Sprintf("Раунд %d завершен. Расчет произведен.", b.Round))
	combatResult := h.calculateCombat(b)
	// ПРИМЕНЯЕМ ЦИФРЫ
	b.AttackerData.HP -= combatResult.DamageToAttacker
	b.DefenderData.HP -= combatResult.DamageToDefender
	b.AttackerData.HP += combatResult.HPToAttacker
	b.DefenderData.HP += combatResult.HPToDefender
	b.AttackerData.Mana += combatResult.ManaToAttacker
	b.DefenderData.Mana += combatResult.ManaToDefender
	b.AttackerData.ClampStats()
	b.DefenderData.ClampStats()
	// Добавляем логи раунда в общий лог боя
	b.Log = append(b.Log, combatResult.RoundLog...)
	// ПРОВЕРКА ПОБЕДЫ
	if b.AttackerData.HP <= 0 || b.DefenderData.HP <= 0 {
		reason := "Бой завершен."
		if b.AttackerData.HP <= 0 && b.DefenderData.HP <= 0 {
			reason = "Оба бойца пали. Ничья!"
		} else if b.AttackerData.HP <= 0 {
			reason = fmt.Sprintf("Победа за %s!", b.DefenderData.Name)
		} else {
			reason = fmt.Sprintf("Победа за %s!", b.AttackerData.Name)
		}
		b.mu.Unlock()
		h.sendBattleUpdates(b)
		h.finishBattle(b, reason)
		return
	}
	b.Round++
	b.AttackerTurn = nil
	b.DefenderTurn = nil

	roundDuration := time.Duration(config.Get().GAME.ROUNDTIME) * time.Second
	b.ExpiresAt = time.Now().Add(roundDuration)

	b.mu.Unlock() // Закончили менять данные боя

	// Рассылаем снапшоты (уже без лока)
	// h.sendBattleUpdates(b)

	// Запускаем таймер на следующий раунд
	h.sendBattleUpdates(b)
	go h.battleTimerGuard(b.ID, b.Round)
}

func (h *Hub) finishBattle(b *Battle, reason string) {
	b.mu.Lock() // Сначала защищаем сам объект боя
	if b.Finished {
		b.mu.Unlock()
		return
	}
	h.mu.Lock()
	b.Finished = true
	delete(h.activeBattles, b.ID)
	delete(h.playerToBattle, b.AttackerData.ID)
	delete(h.playerToBattle, b.DefenderData.ID)
	h.mu.Unlock()
	attacker, aOk := h.GetActiveClient(b.AttackerData.ID)
	if aOk {
		attacker.Character.State = models.StatusFree
		attacker.Character.HP = b.AttackerData.HP
		attacker.Character.Mana = b.AttackerData.Mana

		// СОХРАНЕНИЕ: пишем в базу полную структуру
		go database.UpdateCharacter(attacker.Character)
	}

	defender, dOk := h.GetActiveClient(b.DefenderData.ID)
	if dOk {
		defender.Character.State = models.StatusFree
		defender.Character.HP = b.DefenderData.HP
		defender.Character.Mana = b.DefenderData.HP
		go database.UpdateCharacter(defender.Character)
	}

	// Уведомляем игроков
	endMsg := map[string]interface{}{
		"type":   "battle_end",
		"reason": reason,
	}
	h.Send(attacker, endMsg)
	h.Send(defender, endMsg)
	h.broadcastStateUpdate(attacker.Character)
	h.broadcastStateUpdate(defender.Character)
}

func (h *Hub) calculateCombat(b *Battle) RoundCalculationResult {
	res := RoundCalculationResult{
		RoundLog: make([]string, 0),
	}

	// Коэффициент регенерации (например, 5% от макс)
	regenRate := 0.01

	for i := 0; i < 5; i++ {
		s1 := database.AllSpells[b.AttackerTurn[i]] // Заклинание Агрессора
		s2 := database.AllSpells[b.DefenderTurn[i]] // Заклинание Защитника
		idx1 := elementToIndex[s1.Element]
		idx2 := elementToIndex[s2.Element] // СЛУЧАЙ 1: Атака против Щита
		res.ManaToAttacker -= s1.ManaCost
		res.ManaToDefender -= s2.ManaCost
		if s1.Type == "attack" && s2.Type == "shield" {
			mult := atkVsShield[idx1][idx2]
			// dmg := (float64(b.AttackerData.Strength) * s1.Power) * mult
			dmg := s1.Power * mult
			res.DamageToDefender += dmg
			res.RoundLog = append(res.RoundLog, fmt.Sprintf("Слот %d: %s бьет в %s. Урон: %.1f", i+1, s1.Name, s2.Name, dmg))
		} else if s2.Type == "attack" && s1.Type == "shield" {
			mult := atkVsShield[idx2][idx1]
			// dmg := (float64(b.DefenderData.Strength) * s2.Power) * mult
			dmg := s2.Power * mult
			res.DamageToAttacker += dmg
			res.RoundLog = append(res.RoundLog, fmt.Sprintf("Слот %d: %s отражает %s. Урон вам: %.1f", i+1, s1.Name, s2.Name, dmg))

			// СЛУЧАЙ 2: Атака против Атаки
		} else if s1.Type == "attack" && s2.Type == "attack" {
			mult1 := atkVsAtk[idx1][idx2]
			mult2 := atkVsAtk[idx2][idx1]
			// dmgToDef := (float64(b.AttackerData.Strength) * s1.Power) * mult1
			// dmgToAtk := (float64(b.DefenderData.Strength) * s2.Power) * mult2
			dmgToDef := s1.Power * mult1
			dmgToAtk := s2.Power * mult2
			res.DamageToDefender += dmgToDef
			res.DamageToAttacker += dmgToAtk
			res.RoundLog = append(res.RoundLog, fmt.Sprintf("Слот %d: Встречная атака! %s vs %s. Вы нанесли %.1f, получили %.1f", i+1, s1.Name, s2.Name, dmgToDef, dmgToAtk))

			// СЛУЧАЙ 3: Щит против Щита
		} else if s1.Type == "shield" && s2.Type == "shield" {
			if s1.Element == s2.Element {
				// Одинаковые щиты - реген маны
				mGain1 := b.AttackerData.MaxMana * regenRate
				mGain2 := b.DefenderData.MaxMana * regenRate
				res.ManaToAttacker += mGain1
				res.ManaToDefender += mGain2
				res.RoundLog = append(res.RoundLog, fmt.Sprintf("Слот %d: Резонанс щитов %s. +Мана.", i+1, s1.Element))
			} else {
				// Разные щиты - реген HP
				hGain1 := b.AttackerData.MaxHP * regenRate
				hGain2 := b.DefenderData.MaxHP * regenRate
				res.HPToAttacker += hGain1
				res.HPToDefender += hGain2
				res.RoundLog = append(res.RoundLog, fmt.Sprintf("Слот %d: Сплетение защит %s и %s. +Жизнь.", i+1, s1.Name, s2.Name))
			}
		}
	}
	return res
}

func (h *Hub) sendBattleUpdates(b *Battle) {
	atkInfo := h.getBattleSnapshot(b.ID, b.AttackerData.ID)
	defInfo := h.getBattleSnapshot(b.ID, b.DefenderData.ID)
	attacker, ok := h.GetActiveClient(b.AttackerData.ID)
	if ok {
		h.Send(attacker, map[string]interface{}{
			"type":        "battle_update",
			"battle_info": atkInfo,
		})
	}
	defender, ok := h.GetActiveClient(b.DefenderData.ID)
	if ok {
		h.Send(defender, map[string]interface{}{
			"type":        "battle_update",
			"battle_info": defInfo,
		})
	}
}

func (h *Hub) broadcastStateUpdate(char *models.Character) {
	h.BroadcastToRoomExcept(char.WorldID, char.LocationID, char.ID, map[string]interface{}{
		"type":  "player_state_update",
		"id":    char.ID,
		"state": char.State,
	})
}
