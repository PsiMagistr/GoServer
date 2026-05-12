let gameLoopId = null;
import { network } from "./network.js";
import { CreateCharacterTemplate } from "../templates/create_character.js";
import { gameTemplate } from "../templates/game.js";
import { utils } from "../utils/utils_functions.js";


export function initGame(char) {
    if(gameLoopId){
        cancelAnimationFrame(gameLoopId)
        console.log("Старый цикл остановлен.");       
    }
    const app = document.getElementById('app');
    network.connect(); 
    app.innerHTML = gameTemplate(char);  
    const UISchema = {
        canvas:"#gameCanvas",
        chatInput:"#chat-input",
        chatBtn:"#chat-send-btn",       

    }    
    const UiElements = utils.getElementsBySelectors(UISchema);
    const ctx = UiElements.canvas.getContext("2d");
    const sendMessage = ()=>{
        const text = UiElements.chatInput.value.trim();
        if (text != "" && network.socket){
            const packet = {
                type:"chat_msg",
                text:text,
            }
            network.send(packet);
            UiElements.chatInput.value = "";
        }    
    }
    UiElements.chatBtn.onclick = sendMessage;       
    // Простейший цикл отрисовки
    function draw() {
        ctx.fillStyle = 'black';
        ctx.fillRect(0, 0, UiElements.canvas.width, UiElements.canvas.height);
        ctx.fillStyle = 'lime';
        ctx.fillRect(280, 180, 40, 40); // Наш "игрок"
        gameLoopId = requestAnimationFrame(draw);
    }
    draw();
    console.log("Наш персонаж.")
    console.log(char)
}

 