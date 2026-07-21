package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"GoServer/internal/config"
	"GoServer/internal/database"
	"GoServer/internal/models"
)

type CommandHandler func(c *Client, h *Hub, data map[string]interface{})

type CommitStatsRequest struct {
	Strength  int `json:"strength"`
	Agility   int `json:"agility"`
	Intuition int `json:"intuition"`
	Wisdom    int `json:"wisdom"`
	Charm     int `json:"charm"`
	Vitality  int `json:"vitality"`
}

var commands = map[string]CommandHandler{
	"chat_msg":         handleChat,
	"move":             handleMoveRequest,
	"portal_request":   handlePortalMoveRequest,
	"private_chat":     handleWhisperRequest,
	"commit_stats":     handleStatsCommitRequest,
	"battle_challenge": handleBattleChallenge,
	"battle_accept":    handleBattleAccept,
	"battle_turn":      handleBattleTurn,
	"battle_decline":   handleBattleDecline,
}

func handleChat(c *Client, h *Hub, data map[string]interface{}) {
	rawText, _ := data["text"].(string)
	text := strings.TrimSpace(rawText)
	if text == "" {
		return
	}
	if len([]rune(text)) > 150 {
		text = string([]rune(text)[:150])
	}
	chatPacket := map[string]interface{}{
		"type":   "chat_msg",
		"sender": c.Character.Name,
		"text":   text,
	}
	h.RoomBroadcast <- RoomMessage{
		WorldID:    c.Character.WorldID,
		LocationID: c.Character.LocationID,
		Payload:    chatPacket,
	}
}

func handleMoveRequest(c *Client, h *Hub, data map[string]interface{}) {
	clientStatus := h.GetFullStatus(c.Character.ID)
	if clientStatus != models.StatusFree {
		h.SystemMsg(c, "Вы не можете двигаться.")
		return
	}
	targetID, ok := data["target_id"].(string)
	if !ok {
		h.SystemMsg(c, "Неверный вормат данных.")
		return
	}

	world, exists := Universe[c.Character.WorldID]

	targetNode, exists := world.Points[targetID]
	sourceNode := world.Points[c.Character.LocationID]

	if !exists || targetID == c.Character.LocationID {
		h.SystemMsg(c, "Локации не существует, либо Вы уже находитесь там.")
		return
	}
	dx := float64(targetNode.X - sourceNode.X)
	dy := float64(targetNode.Y - sourceNode.Y)
	result := math.Ceil(math.Sqrt(dx*dx+dy*dy) / 10)

	duration := time.Duration(result) * time.Second
	charID := c.Character.ID
	worldID := c.Character.WorldID
	h.mu.Lock()
	h.movingPlayers[c.Character.ID] = &MoveData{
		DestinationID:      targetID,
		ArrivalTime:        time.Now().Add(duration),
		TargetWorldName:    world.Name,
		TargetLocationName: targetNode.Name,
	}
	h.mu.Unlock()

	h.Send(c, map[string]interface{}{
		"type":          "move_starting",
		"world_name":    world.Name,
		"location_name": targetNode.Name,
		"duration":      duration.Seconds(),
		"state":         h.GetFullStatus(charID),
	})

	go func() {
		time.Sleep(duration)
		h.mu.Lock()
		delete(h.movingPlayers, charID)
		// 1. Обновляем базу (это можно делать без мьютекса)

		activeClient, online := h.Clients[charID]

		if online {
			// Сначала официально "переставляем" игрока в новую комнату в памяти
			oldLockID := activeClient.Character.LocationID
			activeClient.Character.LocationID = targetID
			activeClient.Character.State = models.StatusFree
			// ТЕПЕРЬ собираем соседей. Теперь игрок сам попадет в этот список!
			newNeighbors := h.getNeighbors(activeClient.Character.WorldID, targetID)

			// Получаем данные комнаты для списка порталов
			currentWorld := Universe[activeClient.Character.WorldID]
			currentNode := currentWorld.Points[targetID]
			h.mu.Unlock() // Все операции с данными закончены, отпускаем
			// 3. Отправляем пакет прибытия
			h.Send(activeClient, map[string]interface{}{
				"type":          "move_complete",
				"location_id":   targetID,
				"location_name": currentNode.Name,
				"players":       newNeighbors,
				"worlds":        currentNode.Worlds,
				"state":         h.GetFullStatus(activeClient.Character.ID),
			})

			// 4. Оповещаем остальных
			h.BroadcastToRoomExcept(worldID, oldLockID, charID, map[string]interface{}{
				"type": "player_left",
				"player": map[string]interface{}{
					"id":   charID,
					"name": activeClient.Character.Name,
				},
			})
			h.BroadcastToRoomExcept(worldID, targetID, charID, map[string]interface{}{
				"type":   "player_joined",
				"player": activeClient.Character,
			})
		} else {
			h.mu.Unlock()
		}
		_ = database.UpdateCharacterLocation(charID, targetID)
	}()
}

func handlePortalMoveRequest(c *Client, h *Hub, data map[string]interface{}) {
	clientFullStatus := h.GetFullStatus(c.Character.ID)
	if clientFullStatus != models.StatusFree {
		h.SystemMsg(c, "Вы не можете двигаться.")
		return
	}
	targetWorldID, ok := data["world_id"].(string)
	if !ok {
		return
	}

	targetWorld, exists := Universe[targetWorldID]
	if !exists {
		return
	}

	currentWorld := Universe[c.Character.WorldID]
	currentNode := currentWorld.Points[c.Character.LocationID]

	// Проверяем, есть ли портал в этот мир
	canTeleport := false
	for _, el := range currentNode.Worlds {
		if el.ID == targetWorldID {
			canTeleport = true
			break
		}
	}

	if !canTeleport {
		fmt.Printf("Игрок %s: попытка незаконной телепортации\n", c.Character.Name)
		return
	}

	const portalDuration = 10 * time.Second // Твои 200 секунд
	charID := c.Character.ID
	oldWorldID := c.Character.WorldID
	oldLocID := c.Character.LocationID

	h.mu.Lock()
	h.movingPlayers[charID] = &MoveData{
		DestinationID:      "portal",
		TargetWorldName:    targetWorld.Name,
		TargetLocationName: Universe[targetWorldID].Points["portal"].Name,
		ArrivalTime:        time.Now().Add(portalDuration),
	}
	h.mu.Unlock()

	// Сообщаем о начале долгого перехода
	h.Send(c, map[string]interface{}{
		"type":          "move_starting",
		"world_name":    targetWorld.Name,
		"location_name": Universe[targetWorldID].Points["portal"].Name,
		"duration":      int(portalDuration.Seconds()),
		"state":         h.GetFullStatus(c.Character.ID),
	})

	go func() {
		time.Sleep(portalDuration)
		h.mu.Lock()
		delete(h.movingPlayers, charID)
		// Обновляем БД
		//_ = database.UpdateCharacterWorld(charID, targetWorldID, "portal")
		activeClient, online := h.Clients[charID]

		if online {
			activeClient.Character.WorldID = targetWorldID
			activeClient.Character.LocationID = "portal"
			activeClient.Character.State = models.StatusFree
			h.mu.Unlock()
			// Полная синхронизация для прыгнувшего
			h.Send(activeClient, map[string]interface{}{
				"type":        "world_sync",
				"location_id": "portal",
				"world_id":    targetWorldID,
				"player":      activeClient.Character,
				"world":       Universe[targetWorldID],
				"players":     h.getNeighbors(targetWorldID, "portal"),
				"worlds":      Universe[targetWorldID].Points["portal"].Worlds,
			})

			// Оповещаем старый мир
			h.BroadcastToRoomExcept(oldWorldID, oldLocID, charID, map[string]interface{}{
				"type":   "player_left",
				"player": map[string]interface{}{"id": charID, "name": activeClient.Character.Name},
			})
			// Оповещаем новый мир
			h.BroadcastToRoomExcept(targetWorldID, "portal", charID, map[string]interface{}{
				"type":   "player_joined",
				"player": activeClient.Character,
			})
		} else {
			h.mu.Unlock()
		}
		_ = database.UpdateCharacterWorld(charID, targetWorldID, "portal")
	}()
}

// Шепот для приват-чата.
func handleWhisperRequest(c *Client, h *Hub, data map[string]interface{}) {
	targetName, _ := data["target_name"].(string)
	text := data["text"].(string)

	if text == "" || targetName == "" {
		return
	}
	if targetName == c.Character.Name {
		h.SystemMsg(c, "Вы пытаетесь отправить сообщение самому себе!")
		return
	}
	if len([]rune(text)) > 150 {
		text = string([]rune(text)[:150])
	}
	targetClient := h.GetClientByName(targetName)
	if targetClient == nil {
		h.SystemMsg(c, "Персонаж "+targetName+" не в сети.")
		return
	}
	h.Send(targetClient, map[string]interface{}{
		"type": "whisper_received",
		"from": c.Character.Name,
		"text": text,
	})
	h.Send(c, map[string]interface{}{
		"type": "whisper_sent",
		"to":   targetName,
		"text": text,
	})
}

func handleStatsCommitRequest(c *Client, h *Hub, data map[string]interface{}) {
	statsData, ok := data["stats"]
	if !ok {
		log.Println("Ошибка нет ключа stats")
		return
	}
	var req CommitStatsRequest
	dataBytes, _ := json.Marshal(statsData)
	err := json.Unmarshal(dataBytes, &req)
	if err != nil {
		log.Println("Ошибка демаршалинга статов")
		return
	}
	h.mu.Lock()
	diffStr := req.Strength - c.Character.Strength
	diffAgi := req.Agility - c.Character.Agility
	diffInt := req.Intuition - c.Character.Intuition
	diffVit := req.Vitality - c.Character.Vitality
	diffCharm := req.Charm - c.Character.Charm
	diffWisdom := req.Wisdom - c.Character.Wisdom
	totalSpent := diffStr + diffAgi + diffInt + diffVit + diffCharm + diffWisdom
	if diffStr < 0 || diffAgi < 0 || diffInt < 0 || diffVit < 0 || diffCharm < 0 || diffWisdom < 0 {
		log.Printf("Игрок %s пытался уменьшить характеристики!", c.Character.Name)
		h.mu.Unlock()
		return
	}
	if totalSpent > c.Character.FreePoints {
		log.Printf("Игрок %s пытался потратить %d очков, имея %d",
			c.Character.Name, totalSpent, c.Character.FreePoints)
		h.mu.Unlock()
		return
	}
	updatedChar := *c.Character
	updatedChar.Strength = req.Strength
	updatedChar.Agility = req.Agility
	updatedChar.Intuition = req.Intuition
	updatedChar.Vitality = req.Vitality
	updatedChar.Charm = req.Charm
	updatedChar.Wisdom = req.Wisdom
	updatedChar.MaxHP = 150 + float64(updatedChar.Vitality*2)
	updatedChar.MaxMana = 100 + float64(updatedChar.Wisdom*2)
	updatedChar.FreePoints -= totalSpent
	h.mu.Unlock()
	go func(char models.Character) {
		err := database.UpdateCharacter(&char)
		h.mu.Lock()
		defer h.mu.Unlock()
		activeClient, online := h.Clients[char.ID]
		if online {
			if err != nil {
				log.Printf("Ошибка записи статов в БД: %v", err)
				h.Send(activeClient, map[string]interface{}{
					"type":  "error_msg",
					"error": err,
				})
				return
			}
			// Сверяем, что клиент тот же самый (реконнект)
			*activeClient.Character = char
			h.Send(activeClient, map[string]interface{}{
				"type":   "player_update",
				"player": activeClient.Character,
			})
		}
	}(updatedChar)
}

func handleBattleChallenge(c *Client, h *Hub, data map[string]interface{}) { // Отправляем заявку на бой
	targetIDRaw, ok := data["target_id"]
	if !ok {
		h.SystemMsg(c, "Неверный формат данных.")
		return
	}
	targetIDFloat, ok := targetIDRaw.(float64)
	if !ok {
		// Если вдруг с фронтенда пришла строка или что-то еще - сервер не упадет
		h.SystemMsg(c, "Ошибка: ID персонажа должен быть числом.")
		return
	}
	targetID := int64(targetIDFloat)
	targetClient, online := h.GetActiveClient(targetID)
	if !online {
		h.SystemMsg(c, "Заявка не была подана. Персонаж не в сети.")
		return
	}
	if h.GetFullStatus(targetID) != models.StatusFree || h.GetFullStatus(c.Character.ID) != models.StatusFree {
		h.SystemMsg(c, "Заявка не была подана. Персонаж "+targetClient.Character.Name+" или Вы сами заняты.")
		return
	}
	if c.Character.LocationID != targetClient.Character.LocationID || c.Character.WorldID != targetClient.Character.WorldID {
		h.SystemMsg(c, "Заявка не была подана. Персонаж "+targetClient.Character.Name+" в другом мире или локации.")
		return
	}
	/**/
	h.mu.Lock()
	backChallenge := h.GetInviteFromSpecificPlayer(c.Character.ID, targetID)
	if backChallenge != nil {
		if time.Now().Before(backChallenge.ExpiresAt) {
			h.mu.Unlock()                          // Отпускаем Хаб, так как запуск боя возьмет свои локи
			executeBattleStart(h, targetClient, c) // (Агрессор - тот кто вызвал ПЕРВЫМ)
			return
		}
	}
	////
	if h.challenges[targetID] == nil {
		h.challenges[targetID] = make(map[int64]*BattleChallenge)
	}

	existing := h.GetInviteFromSpecificPlayer(targetID, c.Character.ID)
	if existing != nil {
		if time.Now().Before(existing.ExpiresAt) {
			h.mu.Unlock()
			h.SystemMsg(c, "Вы уже отправили вызов этому персонажу. Дождитесь ответа.")
			return
		}
	}
	expires := time.Now().Add(time.Second * 60)
	timeLeft := int(math.Ceil(time.Until(expires).Seconds()))
	challenge := &BattleChallenge{
		SenderID:   c.Character.ID,
		SenderName: c.Character.Name,
		TargetID:   targetID,
		ExpiresAt:  expires,
		TimeLeft:   timeLeft,
	}
	h.challenges[targetID][c.Character.ID] = challenge
	h.mu.Unlock()
	h.Send(targetClient, map[string]interface{}{
		"type":      "new_challenge",
		"challenge": challenge,
	})
	h.SystemMsg(c, "Вы вызвали "+targetClient.Character.Name+" на бой.")
}

// Принятие боя
func handleBattleAccept(c *Client, h *Hub, data map[string]interface{}) {
	senderIDRaw, ok := data["sender_id"]
	if !ok {
		return
	}
	senderIDFloat, ok := senderIDRaw.(float64)
	if !ok {
		h.SystemMsg(c, "Ошибка: ID персонажа должен быть числом.")
		return
	}
	senderID := int64(senderIDFloat)
	h.mu.Lock()
	myInvites, hasInvites := h.challenges[c.Character.ID]
	if !hasInvites {
		h.mu.Unlock()
		h.SystemMsg(c, "У вас нет активных вызовов.")
		return
	}
	challenge, exists := myInvites[senderID]
	if !exists {
		h.mu.Unlock()
		h.SystemMsg(c, "Вызов от этого игрока не найден.")
		return
	}

	if time.Now().After(challenge.ExpiresAt) {
		delete(h.challenges[c.Character.ID], senderID)
		h.mu.Unlock()
		h.SystemMsg(c, "Срок действия вызова истек.")
		return
	}

	attacker, online := h.Clients[senderID]
	if !online || attacker.Character.State != models.StatusFree {
		h.mu.Unlock()
		h.SystemMsg(c, "Противник уже не в сети или занят.")
		return
	}

	if c.Character.State != models.StatusFree {
		h.mu.Unlock()
		h.SystemMsg(c, "Вы заняты и не можете принять заявку.")
		return
	}

	h.mu.Unlock()
	executeBattleStart(h, attacker, c)
}

func executeBattleStart(h *Hub, attacker *Client, defender *Client) {
	h.mu.Lock()
	if attacker.Character.State != models.StatusFree || defender.Character.State != models.StatusFree {
		h.mu.Unlock()
		h.SystemMsg(attacker, "Заявка не была отправлена. Кто-то из игроков занят")
		return
	}
	battleID := atomic.AddInt64(&h.lastBattleID, 1)
	newBattle := &Battle{
		ID:           battleID,
		AttackerData: *attacker.Character,
		DefenderData: *defender.Character,
		Round:        1,
		ExpiresAt:    time.Now().Add(time.Duration(config.Get().GAME.ROUNDTIME) * time.Second),
	}
	h.activeBattles[battleID] = newBattle
	h.playerToBattle[attacker.Character.ID] = battleID
	h.playerToBattle[defender.Character.ID] = battleID
	attacker.Character.State = models.StatusBattle
	defender.Character.State = models.StatusBattle
	delete(h.challenges, attacker.Character.ID)
	delete(h.challenges, defender.Character.ID)
	atkInfo := h.getBattleSnapshot(battleID, attacker.Character.ID)
	defInfo := h.getBattleSnapshot(battleID, defender.Character.ID)
	h.mu.Unlock()
	h.Send(attacker, map[string]interface{}{"type": "battle_start", "battle_info": atkInfo})
	h.Send(defender, map[string]interface{}{"type": "battle_start", "battle_info": defInfo})
	go h.battleTimerGuard(newBattle.ID, newBattle.Round)
}

func handleBattleTurn(c *Client, h *Hub, data map[string]interface{}) {
	var req BattleTurnRequest
	dataBytes, _ := json.Marshal(data)
	if err := json.Unmarshal(dataBytes, &req); err != nil {
		h.BattleMsg(c, "Ошибка: Неверный формат хода.")
		return
	}

	h.mu.RLock()
	battle, exists := h.activeBattles[req.BattleID]
	h.mu.RUnlock()

	if !exists {
		h.BattleMsg(c, "Ошибка: Бой не найден")
		return
	}

	// 1. БЕРЕМ ЗАМОК ВРУЧНУЮ
	battle.mu.Lock()

	// 2. ПРОВЕРКИ (внутри лока)
	if battle.Finished {
		battle.mu.Unlock() // ВСЕГДА разблокируем перед return!
		h.BattleMsg(c, "Бой уже завершен")
		return
	}

	if battle.Round != req.Round {
		battle.mu.Unlock()
		h.BattleMsg(c, "Ошибка: Раунд уже завершен")
		return
	}

	// 3. ВАЛИДАЦИЯ (тоже внутри лока)
	if err := h.validateBattleTurn(c, req.Spells); err != nil {
		battle.mu.Unlock()
		h.BattleMsg(c, err.Error())
		return
	}

	// 4. ЗАПИСЫВАЕМ ХОД
	if c.Character.ID == battle.AttackerData.ID {
		if battle.AttackerTurn != nil {
			battle.mu.Unlock()
			h.BattleMsg(c, "Вы уже походили.")
			return
		}
		battle.AttackerTurn = req.Spells
	} else if c.Character.ID == battle.DefenderData.ID {
		if battle.DefenderTurn != nil {
			battle.mu.Unlock()
			h.BattleMsg(c, "Вы уже походили.")
			return
		}
		battle.DefenderTurn = req.Spells
	} else {
		battle.mu.Unlock()
		h.BattleMsg(c, "Вы не в этом бою.")
		return
	}

	// 5. ПРОВЕРКА ЗАВЕРШЕНИЯ ХОДОВ
	if battle.AttackerTurn != nil && battle.DefenderTurn != nil {
		// ВАЖНО: Мы закончили менять данные боя в ЭТОЙ функции.
		// Отпускаем замок ПЕРЕД вызовом тяжелой логики.
		battle.mu.Unlock()

		fmt.Println("Все походили, запускаем расчет для боя:", battle.ID)
		h.resolveBattleRound(battle, false)
	} else {
		// Если второй игрок еще не походил
		battle.mu.Unlock()
		h.BattleMsg(c, "Ход принят. Ожидание противника...")
	}
}

// Удаление заявок.
func handleBattleDecline(c *Client, h *Hub, data map[string]interface{}) {
	senderIDRaw, ok := data["sender_id"]
	if !ok {
		return
	}
	senderID := int64(senderIDRaw.(float64))
	if !ok {
		h.SystemMsg(c, "Ошибка: ID персонажа должен быть числом.")
		return
	}
	h.mu.Lock()
	invites, exists := h.challenges[c.Character.ID]
	if exists {
		delete(invites, senderID)
	}
	if len(invites) == 0 {
		delete(h.challenges, c.Character.ID)
	}
	h.mu.Unlock()
	attacker, online := h.GetActiveClient(senderID)
	if online {
		h.SystemMsg(attacker, "Персонаж "+c.Character.Name+" отклонил ваш вызов.")
	}
}
