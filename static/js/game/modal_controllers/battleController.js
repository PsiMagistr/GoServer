import { modalManager } from "../modalManager.js";
import { battleModalTemplate } from "../../templates/battle_modal.js";
import { engine } from "../engine.js";
import {utils} from "../../utils/utils_functions.js";
export const battleController = {
    stopTimerFunc: null,
    async open(data) {
        engine.stopMainLoop();
        this.battleData = data;
        console.log("++++++")
        console.log(this.battleData)
        modalManager.show(battleModalTemplate, this.battleData, { closable: true }, this);
        const you = data.you;
        const opponent = data.opponent;
        this.startTurnTimer(this.battleData.time_left)
        const assetsToLoad = {}
        if (!engine.images || !engine.images.hero) {
            assetsToLoad.hero = `../../assets/avatars/${data.you.gender}/${data.you.avatar_id}.png`;
        }        
        assetsToLoad.opponent =  `../../../assets/avatars/${opponent.gender}/${opponent.avatar_id}.png`;
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
    show() {
        modalManager.show(battleModalTemplate, this.data, { closable: true });       
    },
    onHide() {
        engine.stopBattleLoop();
         if (this.stopTimerFunc) {
            this.stopTimerFunc();
            this.stopTimerFunc = null;
        }
        engine.startMainLoop();
    },
    hide() {       
        modalManager.hide();
    },

    startTurnTimer(seconds){
        if(this.stopTimerFunc){
            this.stopTimerFunc();
        }
        const timerEl = document.getElementById('battle-timer');
        if (!timerEl){
            alert()
            return
        }
        this.stopTimerFunc = utils.createTimer(
            seconds,
            (sec)=>{                
                 timerEl.innerText = `${sec}s`;

        },
        ()=>{

        })

    },

}