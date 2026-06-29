import { modalManager } from "../modalManager.js";
import { battleModalTemplate } from "../../templates/battle_modal.js";
import { engine } from "../engine.js";
export const battleController = {
    async open(data) {
        engine.stopMainLoop();
        this.battleData = data;
        console.log("++++++")
        console.log(this.battleData)
        modalManager.show(battleModalTemplate, this.battleData, { closable: true }, this);
        const opponent = data.opponent; 
        console.log("TTTTTT")
        console.log(opponent)
        const assetsToLoad = {
            /*you: `../../../assets/avatars/${you.gender}/${you.avatar_id}.png`,*/            
            opponent: `../../../assets/avatars/${opponent.gender}/${opponent.avatar_id}.png`,
        };
        try {
            const newImages = await engine.loaderAssets(assetsToLoad);
            engine.images = { ...engine.images, ...newImages };
            // 5. Теперь, когда картинка врага в памяти, запускаем боевой канвас
            const canvas = document.getElementById('battleCanvas');           
            engine.initBattle(canvas, this.battleData);
            console.log("Ресурсы боя загружены, начинаем отрисовку.");

        } catch (e) {
            console.error("Ошибка синхронизации ассетов:", e);
        }

    },
    show(images) {
        modalManager.show(battleModalTemplate, this.data, { closable: true });
        const canvas = document.querySelector("#battleCanvas");
        const ctx = canvas.getContext("2d");
        ctx.drawImage(images.you, 10, 10, 100, 100);
        ctx.drawImage(images.opponent, 120, 10, 100, 100);
    },
    onHide() {
        engine.stopBattleLoop();
        engine.startMainLoop();
    },
    hide() {       
        modalManager.hide();
    },

}