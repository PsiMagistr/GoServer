import { modalManager } from "../modalManager.js";
import { statsModalTemplate } from "../../templates/stats_modal.js";
import { gameState } from "../game.js";
import { gameActions } from "../actions.js";

export const statsController = {
    draft:null,
    open(){
        console.log("Создаю черновик из этого:", gameState.player.strength);
        this.draft = JSON.parse(JSON.stringify(gameState.player));
        this.draft.hasChanges = false;
        this.draft.isLoading = false;
        this.draft.error = "";
        this.show();
    },
    increment(statName){
        if(this.draft.isLoading) return;
        if(this.draft.free_points > 0){
            this.draft[statName]++;
            this.draft.free_points--;
            this.draft.hasChanges = true;
            this.draft.error = "";
            this.refresh();            
        }       
               
    },
    show(){                       
        modalManager.show(statsModalTemplate, this.draft,{closable:true});
    },
    refresh(){
        modalManager.refresh(this.draft)
    },
    hide(){
        modalManager.hide();
    },
    commit(){
        if(!this.draft.hasChanges || this.draft.isLoading) return
        this.draft.isLoading = true;
        this.refresh();
        gameActions.sendPlayerStats(this.draft);               
    },
    onServerSuccess(newPlayerData){
        if(!this.draft) return
        this.draft = JSON.parse(JSON.stringify(newPlayerData));
        this.draft.isLoading = false;
        this.draft.hasChanges = false;               
        // 5. Перерисовываем окно (кнопка Сохранить станет серой, а [+] разблокируются)
        this.refresh();        
    },
    onServerError(error){                     
         this.draft.error = error.Message;
         this.draft.isLoading = false;
         this.refresh();   
    }
}