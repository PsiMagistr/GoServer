import { modalManager } from "../modalManager.js";
import { battleModalTemplate } from "../../templates/battle_modal.js";
import { engine } from "../engine.js";
import { utils } from "../../utils/utils_functions.js";
import { ui } from '../ui.js';
import { gameState } from "../game.js";

export const battleController = {
    battleData: null,
    stopTimerFunc: null,
    slots: [null, null, null, null, null],
    // 1. Точка входа (вызывается из сокета)
    async open(data) {
        this.slots = [null, null, null, null, null];
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
     setLock(locked) {
        this.isLoading = locked;

        // 1. Управляем списками заклинаний
        const lists = document.querySelectorAll('.spells-list');
        lists.forEach(list => {
            if (locked) {
                list.classList.add('locked');
            } else {
                list.classList.remove('locked');
            }
        });

        // 2. Управляем кнопками внизу
        const submitBtn = document.getElementById('btn-submit-turn');
        const surrenderBtn = document.getElementById('btn-surrender');

        if (submitBtn) {
            submitBtn.disabled = locked || !this.isSlotsFull();
            submitBtn.innerText = locked ? "ОЖИДАНИЕ..." : "ПОДТВЕРДИТЬ ХОД";
        }
        
        // Кнопку "Сдаться" обычно оставляют активной, но можно тоже приглушить
        if (surrenderBtn) {
            surrenderBtn.style.opacity = locked ? "0.5" : "1";
        }
    },

    // Вспомогательная проверка
    isSlotsFull() {
        return this.slots.every(s => s !== null);
    },

    startTurnTimer(seconds) {
        if (this.stopTimerFunc) this.stopTimerFunc();

        // 1. Находим элемент таймера ОДИН раз
        const timerEl = document.getElementById('battle-timer');
        if (!timerEl) return;

        // Сбрасываем стили (на случай если прошлый таймер был "красным")
        timerEl.style.color = "";
        timerEl.classList.remove("blink");

        // 2. Создаем "умный" таймер на базе Date.now()
        this.stopTimerFunc = utils.createTimer(
            seconds,
            (sec) => {
                // ПРЯМОЕ ОБНОВЛЕНИЕ: меняем только текст, не трогая остальной HTML
                if (timerEl) {
                    timerEl.innerText = `${sec}s`;

                    if (sec <= 5) {
                        timerEl.style.color = "#ff4d4d";
                    }
                }
            },
            () => {
                if (timerEl) timerEl.innerText = "0s";
            }
        );
    },

    cleanup() {
        if (this.stopTimerFunc) this.stopTimerFunc();
        engine.stopBattleLoop();
        engine.startMainLoop();
    },
    getStats() {
        let shields = 0;
        let attacks = 0;
        let totalMana = 0;
        for (const slot of this.slots) {
            if (slot) {
                if (slot.type == "shield") {
                    shields++;
                }
                if (slot.type == "attack") {
                    attacks++
                }
                totalMana += slot.mana_cost
            }

        }
        return { shields, attacks, totalMana }
    },
    pickSpell(spellId) {
        const spell = gameState.player.spells.find(s => s.id == spellId)
        const stats = this.getStats();
        if (spell.type == "shield" && stats.shields >= 2) {
            alert("Нельзя повесть больше двух щитов");
            return;
        }
        if (spell.type == "attack" && stats.attacks >= 3) {
            alert("Нельзя повесить больше 3 атакующих заклов.");
            return;
        }
        if (stats.totalMana + spell.mana_cost > gameState.player.mana) {
            alert("Недостаточно маны для такой комбинации!");
            return;
        }
        const freeIndex = this.slots.indexOf(null)
        if (freeIndex !== -1) {
            this.slots[freeIndex] = spell;
            this.renderSlots()
        }
        //if(this.slots.length == 5)
    },
    unpickSlots(spellId) {
        this.slots[spellId] = null;
        this.renderSlots();
    },
    renderSlots() {
        const isReady = this.slots.every(s => s !== null)
        for (let i = 0; i < this.slots.length; i++) {
            const slot = this.slots[i]
            let value = "Пусто"
            if (slot) {
                value = this.slots[i].name;
            }
            const slotElement = document.querySelector(`#label-slot-${i + 1}`);
            slotElement.textContent = value;
            const btnSubmitTurn = document.querySelector("#btn-submit-turn");           
            if (btnSubmitTurn) {
                btnSubmitTurn.disabled = !isReady
            }
        }
    },
    end(msg) {
        // 1. Останавливаем боевой таймер
        if (this.stopTimerFunc) {
            this.cleanup()
            this.stopTimerFunc = null;
        }
        const timerEl = document.getElementById('battle-timer');
        if (timerEl) {
            timerEl.innerText = "0s";
            timerEl.style.color = "#555"; // Делаем его серым, "неактивным"
            timerEl.classList.remove("blink"); // Убираем мигание, если оно было
        }

        // 2. Блокируем кнопку "Сделать ход", если она была активна
        const submitBtn = document.getElementById('btn-submit-turn');
        if (submitBtn) submitBtn.disabled = true;

        // 3. Выводим результат в лог боя (красиво)
        const log = document.getElementById('battle-log');
        if (log) {
            log.innerHTML += `<div class="log-entry result"><b>ИТОГ: ${msg.reason}</b></div>`;
            log.scrollTop = log.scrollHeight;
        }
        // 4. САМОЕ ГЛАВНОЕ: Показываем крестик закрытия через менеджер
        modalManager.setClosable(true);
        this.setLock(true)
        console.log("Бой окончен. Окно разблокировано.");
    },

};    