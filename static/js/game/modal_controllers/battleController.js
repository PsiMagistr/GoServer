import { modalManager } from "../modalManager.js";
import { battleModalTemplate } from "../../templates/battle_modal.js";
import { engine } from "../engine.js";
import { utils } from "../../utils/utils_functions.js";
import { ui } from '../ui.js';
import { gameState } from "../game.js";

export const battleController = {
    battleData: null,
    stopTimerFunc: null,

    // 1. Точка входа (вызывается из сокета)
    async open(data) {
        this.battleData = data;
        
        // Останавливаем мир
        engine.stopMainLoop();

        // 2. Рисуем "тело" модалки (сразу, чтобы игрок видел интерфейс)
        this.renderFrame();

        // 3. Запускаем таймер
        this.startTurnTimer(this.battleData.time_left);

        // 4. Наполняем списки заклинаний (теперь элементы уже есть в DOM)
        this.renderSpells();

        // 5. Грузим картинки и запускаем канвас
        await this.setupGraphics();
    },

    // Рисует саму оболочку окна через менеджер
    renderFrame() {
        modalManager.show(battleModalTemplate, this.battleData, { closable: false }, this);
    },

    // Наполняет пустые <ul> списками магии
    renderSpells() {
        const allSpells = gameState.player.spells || [];
        
        ui.renderList(
            "#defense-spells",
            allSpells.filter(s => s.type === "shield"),
            "spell", "spell-item",
            (s) => this.generateSpellContent(s)
        );

        ui.renderList(
            "#attack-spells",
            allSpells.filter(s => s.type === "attack"),
            "spell", "spell-item",
            (s) => this.generateSpellContent(s)
        );
    },

    // Асинхронная подготовка графики
    async setupGraphics() {
        const { you, opponent } = this.battleData;
        const assets = {
            opponent: `/assets/avatars/${opponent.gender}/${opponent.avatar_id}.png`
        };

        // Если вдруг при F5 потеряли свою картинку
        if (!engine.images?.hero) {
            assets.hero = `/assets/avatars/${you.gender}/${you.avatar_id}.png`;
        }

        try {
            const newImages = await engine.loaderAssets(assets);
            engine.images = { ...engine.images, ...newImages };
            
            const canvas = document.getElementById('battleCanvas');
            if (canvas) engine.initBattle(canvas, this.battleData);
        } catch (e) {
            console.error("Battle Graphics Error:", e);
        }
    },

    // Вспомогательный метод для содержимого <li>
    generateSpellContent(s) { 
        return s.name;
    },

    startTurnTimer(seconds) {
        if (this.stopTimerFunc) this.stopTimerFunc();
        const timerEl = document.getElementById('battle-timer');
        if (!timerEl) return;

        this.stopTimerFunc = utils.createTimer(seconds, 
            (sec) => { timerEl.innerText = `${sec}s`; },
            () => { timerEl.innerText = "0s"; }
        );
    },

    cleanup() {
        if (this.stopTimerFunc) this.stopTimerFunc();
        engine.stopBattleLoop();
        engine.startMainLoop();
    }
};    