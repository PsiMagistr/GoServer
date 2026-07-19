import { ui } from "./ui.js";
import { gameState } from "./game.js";
import { engine } from "./engine.js";
import { statsController } from "./modal_controllers/statsController.js";
import { battleController } from "./modal_controllers/battleController.js";
import { utils } from '../utils/utils_functions.js'

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

const changeLabel = (msg) => {
    const locLabel = document.querySelector("#location_label");
    if (locLabel && gameState.world) {
        const node = gameState.world.points[msg.location_id];
        locLabel.textContent = `Мир: ${gameState.world.name} Локация: ${node ? node.name : msg.location_id}`;
    }
}

export const socket_events = {
    async world_sync(msg) {
        console.log("Глобальная синхронизация мира...");
        console.log(msg.challenges)
        console.log("Бой")
        console.log(msg.battle_info)
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
            <span class="p-name">${p.name}</span>`

        );
        ui.renderList("#worlds-list", msg.worlds, "world", "world-link", (w) => w.name);
        if (msg.challenges) {
            const preparedChallenges = msg.challenges.map(ch => ({
                ...ch,
                id: ch.sender_id // Клонируем sender_id в поле id
            }));
            ui.renderList(
                "#challenges-list",
                preparedChallenges,
                "invite",
                "challenge-row",
                (ch) => `
                <span class="challenge-name">${ch.sender_name}</span>
                <div class="challenge-actions">
                    <button class="btn-battle accept" id="accept-${ch.sender_id}">
                        ⚔
                        <span class="mini-timer" id="timer-${ch.sender_id}">${ch.time_left}</span>
                    </button>
                    <button class="btn-battle decline" id="decline-${ch.sender_id}">✕</button>
                </div>
            `
            );
            const ch = msg.challenges;
            console.log("+++++++")
            for (const ch of msg.challenges) {
                const timerStop = utils.createTimer(
                    ch.time_left,
                    (sec) => {
                        const el = document.getElementById(`timer-${ch.sender_id}`);
                        if (el) el.innerText = `${sec}`;
                    },
                    () => {
                        ui.removeItemFromUI("invite", ch.sender_id);
                    }
                )
                gameState.challengeTimers[ch.sender_id] = timerStop
            }
        }
        // 3. ЛОГИКА ОВЕРЛЕЯ (Показываем или скрываем сразу, не дожидаясь загрузки картинок)
        const overlay = document.getElementById('move-overlay');
        if (overlay) {
            overlay.style.display = 'none';
            changeLabel(msg)
        }
        console.log("world_sync")
        console.log(msg.player)
        console.log("Заклинания")
        console.log(msg.player.spells)        
        const assetsToLoad = {
            map: `./assets/maps/${msg.world_id}.png`,
            hero: `./assets/avatars/${msg.player.gender}/${msg.player.avatar_id}.png`
        };
        try {
            engine.images = await engine.loaderAssets(assetsToLoad);
            engine.init(document.getElementById('gameCanvas'));
            gameState.isInitialized = true;           
            // Здесь мы НЕ прячем оверлей, так как его должен спрятать только move_complete
            if (gameState.stopMovingTimer != null) gameState.stopMovingTimer()
            const state = msg.player.state;
            if (state === 1 && msg.duration > 0) {
                if (!overlay) return
                overlay.style.display = 'flex';
                overlay.querySelector('.target-name').innerText = `${msg.world_name}, ${msg.location_name}`;
                let timeLeft = msg.duration;
                const timerEl = overlay.querySelector('.timer-count');
                timerEl.innerText = timeLeft;
                gameState.stopMovingTimer = utils.createTimer(
                    timeLeft,
                    (sec) => {
                        timerEl.innerText = sec;
                    }, () => {
                        console.log("Прыгнули через портал");                        
                    })
            }
            else if (state == 2) {
                battleController.open(msg.battle_info)
            }
        } catch (e) {
            console.error("Ошибка синхронизации ассетов:", e);
        }
    },


    // Вызывается при старте любого перемещения (и при реконнекте после world_sync)
    move_starting(msg) {
        gameState.player.state = msg.state;
        const overlay = document.getElementById('move-overlay');
        if (!overlay) return;
        overlay.style.display = 'flex';
        overlay.querySelector('.target-name').innerText = `${msg.world_name}, ${msg.location_name}`;
        let secondsLeft = Math.floor(msg.duration);
        const timerEl = overlay.querySelector('.timer-count');
        timerEl.innerText = secondsLeft;
        if (gameState.stopMovingTimer) gameState.stopMovingTimer()
        gameState.stopMovingTimer = utils.createTimer(secondsLeft, (sec) => {
            timerEl.innerText = sec;
        },
            () => {
                timerEl.innerText = "Прибытие...";
            });

    },
    move_complete(msg) {
        if (gameState.stopMovingTimer) gameState.stopMovingTimer()
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
    sys_msg(msg) {
        console.log(msg);
        const chatContainer = document.getElementById('chat-messages');
        if (!chatContainer) return;
        const div = document.createElement('div');
        div.className = 'chat-line';
        div.innerHTML = `<span class="sys-msg">[СИСТЕМА]:</span><span class="chat-text"> ${msg.text}</span>`;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    },
    battle_log(msg){        
        const battleChatContainer = document.querySelector("#battle-log");
        if (!battleChatContainer) return;       
        const div = document.createElement('div');
        div.innerHTML = `<span class="sys-msg">[СИСТЕМА]:</span><span class="chat-text"> ${msg.text}</span>`;
        battleChatContainer.appendChild(div);
        battleChatContainer.scrollTop = battleChatContainer.scrollHeight;
    },
    whisper_received(msg) {
        const chatContainer = document.getElementById('chat-messages');
        if (!chatContainer) return;
        const div = document.createElement('div');
        div.className = 'chat-line';
        div.innerHTML = `<span class="chat-receiver">[От ${msg.from}]:</span> <span class="chat-text">${msg.text}</span>`;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    },
    whisper_sent(msg) {
        const chatContainer = document.getElementById('chat-messages');
        if (!chatContainer) return;
        const div = document.createElement('div');
        div.className = 'chat-line';
        div.innerHTML = `<span class="chat-sender">[К ${msg.to}]:</span> <span class="chat-text">${msg.text}</span>`;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    },
    player_update(msg) {
        gameState.player = msg.player;
        statsController.onServerSuccess(gameState.player);

    },
    error_msg(msg) {
        statsController.onServerError(msg.error);
    },
    new_challenge(msg) {
        const ch = msg.challenge;
        ui.addItemToList(
            "#challenges-list", // ID твоего <ul> для заявок
            { id: ch.sender_id, name: ch.sender_name },
            "invite",
            "challenge-row",
            (item) => `
                <span class="challenge-name">${ch.sender_name}</span>
                <div class="challenge-actions">
                    <button class="btn-battle accept" id="accept-${ch.sender_id}">
                        ⚔
                        <span class="mini-timer" id="timer-${item.id}">${ch.time_left}</span>
                    </button>
                    <button class="btn-battle decline" id="decline-${item.id}">✕</button>
                </div>
            `
        )
        const timerStop = utils.createTimer(
            ch.time_left,
            (sec) => {
                const el = document.getElementById(`timer-${ch.sender_id}`);
                if (el) el.innerText = `${sec}`;
            },
            () => {
                ui.removeItemFromUI("invite", ch.sender_id);
            }
        )
        gameState.challengeTimers[ch.sender_id] = timerStop;
    },
    battle_start(msg) {       
        battleController.open(msg.battle_info)       
    },
    battle_update(msg){
        battleController.open(msg.battle_info)        
    },
    battle_end(msg){
        alert(msg.reason)
    }
}

