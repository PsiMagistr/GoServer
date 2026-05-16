let gameLoopId = null;
import { network } from "./network.js";
import { engine } from "./engine.js";
import { CreateCharacterTemplate } from "../templates/create_character.js";
import { gameTemplate } from "../templates/game.js";
import { utils } from "../utils/utils_functions.js";
export const gameState = {
    world:null,
    player:null,
    isInitialized:false,
}


export function initGame(char) {    
    const app = document.getElementById('app');    
    app.innerHTML = gameTemplate(char);  
    const UISchema = {        
        canvas:"#gameCanvas",
        chatInput:"#chat-input",
        chatBtn:"#chat-send-btn",

    }    
    const UiElements = utils.getElementsBySelectors(UISchema);
    engine.init(UiElements.canvas);
    network.connect(); 
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
    //console.log("Наш персонаж.")
    //console.log(char)
}

 