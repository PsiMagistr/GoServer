import { gameState } from "./game.js";
import { gameActions } from "./actions.js";
import { engine } from "./engine.js";
import { modalManager } from "./modalManager.js";
import { statsController } from "./modal_controllers/statsController.js";
import { battleController } from "./modal_controllers/battleController.js";
import { contextNenu } from "../templates/context_menu.js";
import { ui } from "./ui.js"

const CLICK_RADIUS = 20;
const CLICK_RADIUS_SQ = CLICK_RADIUS * CLICK_RADIUS;
const clickers = {
    "gameCanvas": (obj, event) => {
        obj.handleCanvasClick(event);
    },
    "chat-send-btn": (obj, event) => {
        obj.sendChat();
    },
    "stats": (obj, event) => {
        statsController.open();
    },
    "save-stats-btn": (obj, event) => {
        statsController.commit();
    },
    "modal-close-btn": (obj, event) => {
        console.log("какой шаблон активен")
        modalManager.hide();
    },
    "btn-submit-turn": (obj, event) => {
        console.log("=======")        
        const battleId = battleController.battleData.battle_id
        const round = battleController.battleData.round;
        const spells = battleController.slots.map(slot=>slot.id);    
        gameActions.sendBattleTurn(battleId, round, spells)
    }
}
const movers = {
    "gameCanvas": (obj, event) => {
        const node = obj.getNodeAt(event.offsetX, event.offsetY);
        const found = node ? node.id : null;
        if (gameState.hoveredNodeId !== found) {
            gameState.hoveredNodeId = found;
            event.target.style.cursor = (found && found !== gameState.player.location_id) ? 'pointer' : 'default';
        }
    },
}
export const interaction = {
    app: null,
    init(appElement) {
        this.app = appElement;
        this.app.onclick = this.handleGlobalClick.bind(this);
        this.app.onmousemove = this.handleMouseMove.bind(this);
        this.app.onkeypress = this.handleGlobalKeyPress.bind(this);
    },
    getNodeAt(offsetX, offsetY) {
        const points = gameState.world?.points;
        if (!points) return;
        for (const id in points) {
            const node = points[id];
            const dx = offsetX - node.x;
            const dy = offsetY - node.y;
            if ((dx * dx + dy * dy) < CLICK_RADIUS_SQ) {
                return node;
            }
        }
        return null;
    },
    handleGlobalClick(event) {
        const target = event.target;
        this.hideContextMenu();
        // 1. Проверка через диспетчер (ищем ID у цели или ближайшего родителя)
        const clickable = target.closest('[id]');
        const id = clickable.id;
        if (clickable && clickers[id]) {
            clickers[id](this, event)
            return
        }
        const worldLink = event.target.closest(".world-link");
        if (worldLink) {
            this.handlePortalClick(worldLink)
        }
        const playerLink = event.target.closest(".player-link")
        if (playerLink) {
            this.showContextMenu(event, playerLink)
            return
        }
        if (id.startsWith("add-")) {
            console.log(event.target)
            const stateName = event.target.dataset.state_name;
            statsController.increment(stateName);
            return
        }
        const contextBtn = target.closest('.context-menu-btn');
        if (contextBtn) {
            this.handleContextAction(contextBtn.dataset.action, contextBtn.dataset.id);
            return;
        }
        //Обработка кликов кнопок на приглашениях.
        if (id.startsWith("accept-") || id.startsWith("decline")) {
            const senderId = id.split("-")[1];
            const stopTimer = gameState.challengeTimers[senderId]
            if (stopTimer) {
                stopTimer();
                delete gameState.challengeTimers[senderId];
            }
            if (id.startsWith("accept-")) {
                console.log("Принимаем заявку на бой")
                gameActions.acceptBattle(senderId);
            } else {
                ui.removeItemFromUI("invite", senderId);
            }

        }
        const spellItem = target.closest(".spell-item");
        if (spellItem) {
            battleController.pickSpell(spellItem.dataset.id)
        }
        const slotItem = target.closest(".battle-slot");
        if (slotItem) {            
            const index = parseInt(slotItem.dataset.slotIndex) - 1;            
            battleController.unpickSlots(index);

        }
    },
    handleMouseMove(event) {
        if (movers[event.target.id]) {
            movers[event.target.id](this, event);
        }
    },
    handleGlobalKeyPress(event) {
        if (event.key === 'Enter' && event.target.id === 'chat-input') {
            this.sendChat()
        }
    },
    handleCanvasClick(event) {//Перемещение по локациям
        const node = this.getNodeAt(event.offsetX, event.offsetY);
        if (!node || node.id === gameState.player.location_id) return;
        if (gameState.player.state == 1) {
            alert("Вы уже в пути")
            return
        }
        const goToQuestion = confirm(`Вы хотите перейти в ${node.name}`);
        if (!goToQuestion) return
        gameActions.moveToNode(node.id);
    },
    sendChat() {
        const toInput = document.getElementById('chat-to');
        const msgInput = document.getElementById('chat-input');
        if (!msgInput) return;
        const recipient = toInput.value.trim();
        const text = msgInput.value.trim();
        if (text === "") return;
        if (recipient !== "") {
            gameActions.sendWhisper(recipient, text)
        }
        else {
            gameActions.sendPublicChat(text);
        }
        msgInput.value = "";
    },
    handlePortalClick(element) {
        if (gameState.player.state == 1) {
            alert("Вы уже в пути")
            return
        }
        const worldId = element.dataset.id; // Берем ID из data-world-id        
        const worldName = element.innerText.replace('🌀', '').trim();
        const goToQuestion = confirm(`Вы уверены, что хотите покинуть этот мир и отправиться в "${worldName}"? (Переход займет 200 сек)`)
        if (!goToQuestion) return
        console.log("Запрос на телепортацию в:", worldId);
        gameActions.enterPortal(worldId)
    },

    showContextMenu(event, element) {
        const menu = document.getElementById('context-menu');
        const charId = element.dataset.id;
        const charName = element.dataset.name
        if (!charId || !charName) {
            console.warn("Попытка открыть меню для пустого или удаленного элемента");
            return;
        }
        if (parseInt(charId) === gameState.player.id) return;
        // Наполняем меню кнопками
        menu.innerHTML = contextNenu(charId, charName);
        // Позиционируем меню справа от клика
        menu.style.display = 'flex';
        menu.style.left = `${event.pageX + 5}px`;
        menu.style.top = `${event.pageY + 5}px`;
    },

    hideContextMenu() {
        const menu = document.getElementById('context-menu');
        if (menu) menu.style.display = 'none';
    },
    handleContextAction(action, id) {
        const name = event.target.dataset.name; // Достаем имя из дата-атрибута кнопки

        switch (action) {
            case 'private':
                document.getElementById('chat-to').value = name;
                document.getElementById('chat-input').focus();
                break;
            case 'challenge':
                if (confirm(`Вызвать на бой ${name}?`)) {
                    gameActions.sendBattleChallenge(id);
                }
                break;
            case 'trade':
                alert("Торговля будет доступна позже");
                break;
        }
    },

}