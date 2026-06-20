import { modalManager } from "../modalManager.js";
import {battleModalTemplate } from "../../templates/battle_modal.js";
export const battleController = {
    open(data){
        this.data = data;
        this.show()
    },
    show(){
        modalManager.show(battleModalTemplate, this.data);   
        
    }
      

}