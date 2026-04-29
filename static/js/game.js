let gameLoopId = null;
import { CreateCharacterTemplate } from "./templates/create_character.js";
import { gameTemplate } from "./templates/game.js";

export function initGame(username) {
    if(gameLoopId){
        cancelAnimationFrame(gameLoopId)
        console.log("Старый цикл остановлен.");       
    }
    const app = document.getElementById('app');
    app.innerHTML = gameTemplate(username);

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
}

 