import { ui } from "./ui.js";
import { gameState } from "./game.js";
import { engine } from "./engine.js";
export const socket_events = {
    room_presence(msg) {
        console.log("Игроки в нашей комнате:");
        //ui.renderPlayerList(msg.worlds, msg.players)
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
    async self_load(msg) {
        if (!gameState.isInitialized) {
            console.log("*******")            
            const location = document.querySelector("#location_label");           
            gameState.world = msg.world;                       
            gameState.world = msg.world;
            gameState.player = msg.player;
            location.textContent = `Локация: ${gameState.world.points[gameState.player.location_id].name}`
            gameState.isInitialized = true;
            const assetsToLoad = {
                map: `/assets/maps/${gameState.player.world_id}.png`,
                hero: `/assets/avatars/${gameState.player.gender}/${gameState.player.avatar_id}.png`
            }
            try{
                 engine.images = await engine.loaderAssets(assetsToLoad);
            }
            catch(error){
                 console.error("Критическая ошибка:", error);
            }           
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
        gameState.isMoving = false;
        gameState.player.location_id = msg.location_id;
        const location = document.querySelector("#location_label");
        location.textContent = `Локация: ${gameState.world.points[gameState.player.location_id].name}`;
        // Скрываем оверлей
        const overlay = document.getElementById('move-overlay');
        if (overlay) overlay.style.display = 'none';        
        console.log("Вы прибыли в", msg.location_id);
    }
    

}