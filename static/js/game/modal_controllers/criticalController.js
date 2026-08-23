import { modalManager } from "../modalManager.js";
import { modalCriticalTemplate } from "../../templates/critical_modal.js";
export const criticalController = {
    open(data){
        modalManager.show(modalCriticalTemplate, data, { closable: true }, this);        
    },
    hide(){
        modalManager.hide();    
    }
    
}