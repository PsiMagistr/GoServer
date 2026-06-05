import { modalManager } from "../modalManager.js";
import { statsModalTemplate } from "../../templates/stats_modal.js";
import { gameState } from "../game.js";
import { network } from "../network.js";

export const statsController = {
    draft:null,
    open(){
        this.draft = JSON.parse(JSON.stringify(gameState.player));
        this.draft.hasChanges = false;
        this.render();
    },
    increment(statName){
        if(this.draft.free_points > 0){
            this.draft[statName]++;
            this.draft.free_points--;
            this.draft.hasChanges = true;
            this.render()            
        }       
               
    },
    render(){
        const html = statsModalTemplate(this.draft)        
        modalManager.show(html);
    },
    hide(){
        modalManager.hide();
    },
    commit(){
        this.hide();         
    }
}