let gameLoopId = null;
import { CreateCharacterTemplate } from "./templates/create_character.js";
import { gameTemplate } from "./templates/game.js";

export function initGame(char) {
    if(gameLoopId){
        cancelAnimationFrame(gameLoopId)
        console.log("Старый цикл остановлен.");       
    }
    const app = document.getElementById('app');
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onopen = () => {
        console.log("Связь с сервером установлена!");
        // Теперь сервер знает, что "Виктория" в сети
    };

    socket.onmessage = (event) => {
        // Сюда будут прилетать сообщения от сервера
        const msg = JSON.parse(event.data);
        if(msg.type === "room_presence"){
            console.log("Старые игроки в комнате:");            
            for(const player of msg.players){
                console.log(player.name);
            }
        }
        else if(msg.type === "player_joined"){
            console.log(`К нам присоединился(лась) ${msg.player.name}`);
        }
        else if(msg.type === "player_left"){
           console.log(`Нас покинул(ла) ${msg.player.name}`);
        }
        //console.log("Получено от сервера:", msg);
    };

    socket.onclose = (event) => {
        console.log("Соединение разорвано", event.reason);        
    };

    socket.onerror = (error) => {
        console.error("Ошибка сокета:", error);
    };
    // Сохраним сокет в глобальную переменную или передадим в движок,
    // чтобы можно было отправлять сообщения позже (например, при ударе)
    window.gameSocket = socket; 
    app.innerHTML = gameTemplate(char);
    const canvas = document.getElementById('gameCanvas');
    const ctx = canvas.getContext('2d');
    // Простейший цикл отрисовки
    function draw() {
        ctx.fillStyle = 'black';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = 'lime';
        ctx.fillRect(280, 180, 40, 40); // Наш "игрок"
        gameLoopId = requestAnimationFrame(draw);
    }
    draw();
    console.log("Наш персонаж.")
    console.log(char)
}

 