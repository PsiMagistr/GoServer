import { utils } from "../utils/utils_functions.js";

export const ui = {
    renderPlayerList(worlds, players){
    const uiElements = utils.getElementsBySelectors({ list: "#players-list" });
        if (!uiElements.list) return;
        uiElements.list.innerHTML = "";
        for(let player of players){
            this.addPlayerToUI(player);
        } 
        console.log("54554545")
        console.log(worlds)   
    },
    addPlayerToUI(player){
            const uiElements = utils.getElementsBySelectors({ list: "#players-list" });
            if (!uiElements.list) return;           
            if (document.getElementById(`player-${player.id}`)) return;
            const li = document.createElement('li');
            li.id = `player-${player.id}`;
            li.className = 'player-link';
            li.innerText = `${player.name}`
            uiElements.list.appendChild(li);    
    },
    removePlayerFromUI(playerID){               
         const uiElements = utils.getElementsBySelectors({item: `#player-${playerID}`});
         if (!uiElements.item) return;        
         uiElements.item.remove()    
    },
}