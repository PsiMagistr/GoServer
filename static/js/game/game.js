let gameLoopId = null;
import { network } from "./network.js";
import { engine } from "./engine.js";
import { CreateCharacterTemplate } from "../templates/create_character.js";
import { gameTemplate } from "../templates/game.js";
import { utils } from "../utils/utils_functions.js";
import { interaction } from "./interaction.js";
export const gameState = {
    world:null,
    player:null,
    isInitialized:false,
    hoveredNodeId:null,
    
}


export function initGame(char) {    
    const app = document.getElementById('app');    
    app.innerHTML = gameTemplate(char); 
    const UISchema = {               
        canvas:"#gameCanvas",      
    }  
    const UiElements = utils.getElementsBySelectors(UISchema);    
    interaction.init(app);
    engine.init(UiElements.canvas);    
    network.connect();
}

 