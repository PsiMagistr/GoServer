import { modalManager } from "../modalManager.js";
import { battleModalTemplate } from "../../templates/battle_modal.js";
import { engine } from "../engine.js";
export const battleController = {
    async open(data) {
        this.data = data;       
        const you = data.you;
        const opponent = data.opponent;
        const assetsToLoad = {
            you: `../../../assets/avatars/${you.gender}/${you.avatar_id}.png`,
            opponent: `../../../assets/avatars/${opponent.gender}/${opponent.avatar_id}.png`,
        };
        try {            
            const images = await engine.loaderAssets(assetsToLoad);
            this.show(images)
        }catch(e){
            console.error("Ошибка синхронизации ассетов:", e);                        
        }    
        
    },
    show(images) {
        modalManager.show(battleModalTemplate, this.data, { closable: true });
        const canvas = document.querySelector("#battleCanvas");
        const ctx = canvas.getContext("2d");
        ctx.drawImage(images.you, 10,10, 100,100);
        ctx.drawImage(images.opponent, 120, 10, 100, 100);
    }

}