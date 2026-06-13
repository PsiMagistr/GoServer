import { ui } from "./ui.js";
import { gameState } from "./game.js";
import { engine } from "./engine.js";
import { statsController } from "./modal_controllers/statsController.js";

// Глобальная переменная для управления таймером внутри этого файла
let moveInterval = null;

const stopTimer = () => {
    if (moveInterval) {
        clearInterval(moveInterval);
        moveInterval = null;
    }
}

 

const clearMoveTimer = () => {
    if (moveInterval) {
        clearInterval(moveInterval);
        moveInterval = null;
    }
};

const changeLabel = (msg)=>{
    const locLabel = document.querySelector("#location_label");
    if (locLabel && gameState.world) {
        const node = gameState.world.points[msg.location_id];
        locLabel.textContent = `Мир: ${gameState.world.name} Локация: ${node ? node.name : msg.location_id}`;
    }
}

export const socket_events = {   
    async world_sync(msg) {
        console.log("Глобальная синхронизация мира...");        
        // 1. Синхронизируем состояние
        gameState.player = msg.player;        
        gameState.world = msg.world;
        gameState.isInitialized = false;
        // 2. Обновляем списки
        ui.renderList(
            '#players-list',
             msg.players,
            "player",
            'player-link',
            (p) => `
            <span class="p-name">${p.name}</span>
            <!--<button class="challenge-btn" id="challenge-${p.id}">⚔️</button>-->
            `           
        );
        ui.renderList("#worlds-list", msg.worlds, "world", "world-link", (w) => w.name);

        // 3. ЛОГИКА ОВЕРЛЕЯ (Показываем или скрываем сразу, не дожидаясь загрузки картинок)
        const overlay = document.getElementById('move-overlay');
        console.log("world_sync")
        console.log(msg.player)
        const state = msg.player.state
        if (state === 1 && msg.duration > 0) {
            /*gameState.isMoving = true;*/
            stopTimer()
            if (overlay) {
                overlay.style.display = 'flex';
                // Подставляем красивые названия, которые Go прислал в world_sync
                overlay.querySelector('.target-name').innerText = `${msg.world_name}, ${msg.location_name}`;                
                let timeLeft = msg.duration;
                const timerEl = overlay.querySelector('.timer-count');
                timerEl.innerText = timeLeft;

                // Запускаем локальный отсчет для плавности
                moveInterval = setInterval(() => {
                    timeLeft--;
                    if (timeLeft >= 0) {
                        timerEl.innerText = timeLeft;
                    } else {
                        timerEl.innerText = "Прибытие...";
                        stopTimer();
                    }
                }, 1000);
            }
        } else {
            // Если не движемся — гарантированно прячем
            if (overlay) overlay.style.display = 'none';
            gameState.isMoving = false;            
            stopTimer(); 
            /////////////////////////          
            changeLabel(msg);
        }

        // 4. ЗАГРУЗКА РЕСУРСОВ (в фоне)
        const assetsToLoad = {
            map: `/assets/maps/${msg.world_id}.png`,
            hero: `/assets/avatars/${msg.player.gender}/${msg.player.avatar_id}.png`
        };

        try {
            engine.images = await engine.loaderAssets(assetsToLoad);
            engine.init(document.getElementById('gameCanvas'));
            gameState.isInitialized = true;
            // Здесь мы НЕ прячем оверлей, так как его должен спрятать только move_complete
        } catch (e) {
            console.error("Ошибка синхронизации ассетов:", e);
        }
    },


    // Вызывается при старте любого перемещения (и при реконнекте после world_sync)
    move_starting(msg) {
        gameState.player.state = msg.state;        
        stopTimer();

        const overlay = document.getElementById('move-overlay');
        if (!overlay) return;

        overlay.style.display = 'flex';
        overlay.querySelector('.target-name').innerText = `${msg.world_name}, ${msg.location_name}`;
        
        let secondsLeft = Math.floor(msg.duration);
        const timerEl = overlay.querySelector('.timer-count');
        timerEl.innerText = secondsLeft;

        moveInterval = setInterval(() => {
            secondsLeft--;
            if (secondsLeft >= 0) {
                timerEl.innerText = secondsLeft;
            } else {
                timerEl.innerText = "Прибытие...";
                stopTimer();
            }
        }, 1000);
    },

    move_complete(msg) {
        stopTimer();       
        gameState.player.state = msg.state;              
        // Скрываем оверлей
        const overlay = document.getElementById('move-overlay');
        if (overlay) overlay.style.display = 'none';
        // Если это был прыжок между мирами, world_sync уже прилетел или прилетит,
        // но локацию обновим здесь для надежности
        gameState.player.location_id = msg.location_id;       
        
        // Обновляем заголовок локации в UI
       changeLabel(msg)
       ui.renderList('#players-list', msg.players, "player", 'player-link', (p) => p.name);
       ui.renderList("#worlds-list", msg.worlds, "world", "world-link", (w) => w.name);        
    },

    player_left(msg) {
        ui.removeItemFromUI("player", msg.player.id);
    },

    player_joined(msg) {
        ui.addItemToList("#players-list", msg.player, "player", "player-link", (p) => p.name);
    },

    chat_msg(msg) {
        const chatContainer = document.getElementById('chat-messages');
        if (!chatContainer) return;
        const div = document.createElement('div');
        div.className = 'chat-line';
        div.innerHTML = `<span class="chat-sender">${msg.sender}:</span> <span class="chat-text">${msg.text}</span>`;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    },
    sys_msg(msg){
        console.log(msg);
        const chatContainer = document.getElementById('chat-messages');
        if (!chatContainer) return;
        const div = document.createElement('div');       
        div.className = 'chat-line';
        div.innerHTML = `<span class="sys-msg">[СИСТЕМА]:</span><span class="chat-text"> ${msg.text}</span>`;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    },
    whisper_received(msg){
        const chatContainer = document.getElementById('chat-messages');
        if (!chatContainer) return;
        const div = document.createElement('div');       
        div.className = 'chat-line';
        div.innerHTML = `<span class="chat-receiver">[От ${msg.from}]:</span> <span class="chat-text">${msg.text}</span>`;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    },
    whisper_sent(msg){
        const chatContainer = document.getElementById('chat-messages');
        if (!chatContainer) return;
        const div = document.createElement('div');       
        div.className = 'chat-line';
        div.innerHTML = `<span class="chat-sender">[К ${msg.to}]:</span> <span class="chat-text">${msg.text}</span>`;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    },    
    player_update(msg){               
        gameState.player = msg.player;
        statsController.onServerSuccess(gameState.player);         
       
    },
    error_msg(msg){
        statsController.onServerError(msg.error);
    },
    new_challenge(msg){
        console.log("Приглашение на бой")
        console.log(msg)
    },


}

