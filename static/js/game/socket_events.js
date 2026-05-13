import { ui } from "./ui.js";
export const socket_events = {
    room_presence(msg){
         console.log("Игроки в нашей комнате:");         
         ui.renderPlayerList(msg.players)          
    },
    self_load(msg){
         console.log(`Личные данные:`);
         console.log(msg.player);         
    },
    player_joined(msg){
        console.log(`К нам присоединился(лась) ${msg.player.name}`);
        ui.addPlayerToUI(msg.player)
    },
    player_left(msg){
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