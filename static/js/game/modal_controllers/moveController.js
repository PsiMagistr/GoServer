import { modalManager } from "../modalManager.js";
import { moveModalTemplate } from "../../templates/move_modal.js";
import { utils } from "../../utils/utils_functions.js";
import { gameState } from "../game.js";

export const moveController = {
    stopTimerFunc:null,
    open(data){
        modalManager.show(moveModalTemplate, data, { closable: false }, this);
        this.startCountDown(data.duration);
    },
    startCountDown(seconds){
        if(this.stopTimerFunc) this.stopTimerFunc();
        const timerEl = document.getElementById('move-timer');
        this.stopTimerFunc = utils.createTimer(seconds, 
            (sec) => {
                if (timerEl) timerEl.innerText = sec;
            },
            () => {
                if (timerEl) timerEl.innerText = "Прибытие...";
            }
        )

    },
    onHide() {
        if (this.stopTimerFunc) {
            this.stopTimerFunc();
            this.stopTimerFunc = null;
        }
    },
    hide() {
        modalManager.hide();
    }
}