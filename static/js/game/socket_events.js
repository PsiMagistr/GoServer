import { ui } from "./ui.js";
import { gameState } from "./game.js";
import { engine } from "./engine.js";
export const socket_events = {
    room_presence(msg) {
        console.log("Игроки в нашей комнате:");        
        ui.renderList(
            '#players-list',
            msg.players,
            "player",
            'player-link',
            (p)=>`${p.name}`,
        );
        ui.renderList(
            "#worlds-list",
            msg.worlds,
            "world",
            "world-link",
            (w) => `Портал в ${w.name}`,
        )
    },
    async world_sync(msg) {       
        console.log("Глобальная синхронизация мира...");        
        gameState.isMoving = false;
        gameState.player = msg.player;
        gameState.world = msg.world;        
        // Сбрасываем флаг, чтобы форсировать загрузку новой карты
        gameState.isInitialized = false;
        // 2. Сразу обновляем список игроков (передаем массив соседей)
        ui.renderList(
            '#players-list',
            msg.players,
            "player",
            'player-link',
            (p)=>`${p.name}`,
        );
        ui.renderList(
            "#worlds-list",
             msg.worlds,
            "world",
            "world-link",
            (w) => `Портал в ${w.name}`,
        )
        // 3. Запускаем загрузку ресурсов
        const assetsToLoad = {
            map: `/assets/maps/${msg.world_id}.png`,
            hero: `/assets/avatars/${msg.player.gender}/${msg.player.avatar_id}.png`
        };

        try {
            engine.images = await engine.loaderAssets(assetsToLoad);
            engine.init(document.getElementById('gameCanvas'));
            gameState.isInitialized = true;            
            // Скрываем оверлей только КОГДА ВСЁ ЗАГРУЗИЛОСЬ
            document.getElementById('move-overlay').style.display = 'none';
        } catch (e) {
            console.error("Ошибка синхронизации:", e);
        }
    },
        
    player_joined(msg) {
        ui.addItemToList(
            "#players-list", 
            msg.player, 
            "player", 
            "player-link", 
            (p) => `${p.name}`
        );
    },
    player_left(msg) {
        console.log("+++++++++")
        console.log(msg.player)
        console.log(`Нас покинул(ла) ${msg.player.name}`);
        ui.removeItemFromUI("player", msg.player.id);
    },
    chat_msg(msg) {
        const chatContainer = document.getElementById('chat-messages');
        if (!chatContainer) return;
        // Создаем элемент сообщения
        const div = document.createElement('div');
        div.className = 'chat-line';
        div.innerHTML = `<span class="chat-sender">${msg.sender}:</span> <span class="chat-text">${msg.text}</span>`;
        chatContainer.appendChild(div);
        // Авто-прокрутка чата вниз
        chatContainer.scrollTop = chatContainer.scrollHeight;
    },

    move_starting(msg){        
        gameState.isMoving = true;        
        // Показываем оверлей (создадим его в шаблоне)
        const overlay = document.getElementById('move-overlay');
        if (overlay) {
            overlay.style.display = 'flex';
            overlay.querySelector('.target-name').innerText = msg.target_name;
            
            let secondsLeft = msg.duration;
            const timerEl = overlay.querySelector('.timer-count');
            timerEl.innerText = secondsLeft;

            // Локальный таймер для плавности
            //clearInterval(moveInterval);
            let moveInterval = setInterval(() => {
                secondsLeft--;
                if (secondsLeft >= 0) timerEl.innerText = secondsLeft;
                if (secondsLeft <= 0) clearInterval(moveInterval);
            }, 1000);           
        }
    },

    move_complete(msg){
        console.log("move_complit") 
        if (msg.world_id && msg.world_id !== gameState.player.world_id) {
            console.log(`Смена мира: ${gameState.player.world_id} -> ${msg.world_id}`);        
            gameState.player.world_id = msg.world_id;     
            gameState.isInitialized = false;           
        }
        gameState.player.location_id = msg.location_id;
        const location = document.querySelector("#location_label");
        location.textContent = `Локация: ${gameState.world.points[gameState.player.location_id].name}`;
        // Скрываем оверлей
        const overlay = document.getElementById('move-overlay');
        if (overlay) overlay.style.display = 'none';        
        console.log("Вы прибыли в", msg.location_id);
        gameState.isMoving = false;
    },

}