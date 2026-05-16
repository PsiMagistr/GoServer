import { ui } from "./ui.js";
import { gameState } from "./game.js";
import { engine } from "./engine.js";
export const socket_events = {
    room_presence(msg) {
        console.log("Игроки в нашей комнате:");
        ui.renderPlayerList(msg.players)
    },
    async self_load(msg) {
        if (!gameState.isInitialized) {
            gameState.world = msg.world;
            console.log("++++++");
            console.log(msg);           
            gameState.world = msg.world;
            gameState.player = msg.player;
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
        console.log(`К нам присоединился(лась) ${msg.player.name}`);
        ui.addPlayerToUI(msg.player)
    },
    player_left(msg) {
        console.log("+++++++++")
        console.log(msg.player)
        console.log(`Нас покинул(ла) ${msg.player.name}`);
        ui.removePlayerFromUI(msg.player.id);
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
    }

}