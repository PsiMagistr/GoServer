import { gameState } from "./game.js";
import { gameActions } from "./actions.js";
import { engine } from "./engine.js";
/*import {modalManager} from "./modalManager.js";
import {statsModalTemplate} from "../templates/stats_modal.js";*/
import { statsController } from "./modal_controllers/statsController.js";
const CLICK_RADIUS = 20;
const CLICK_RADIUS_SQ = CLICK_RADIUS * CLICK_RADIUS;
const clickers = {
    "gameCanvas":(obj, event)=>{
        obj.handleCanvasClick(event);
    },
    "chat-send-btn":(obj, event)=>{        
        obj.sendChat();    
    },
    "stats":(obj, event)=>{
        statsController.open();
    },
    "save-stats-btn":(obj, event)=>{
        statsController.commit();
    },
    "modal-close-btn":(obj, event)=>{
        statsController.hide();
    },   
}
const movers = {
    "gameCanvas":(obj, event)=>{       
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
        if(clickable && clickers[id]){
            clickers[id](this, event)
            return    
        }
        const worldLink = event.target.closest(".world-link");
        if (worldLink) {
            this.handlePortalClick(worldLink)
        }
        const playerLink = event.target.closest(".player-link")
        if(playerLink){
            this.showContextMenu(event, playerLink)
            return
        }
        if(id.startsWith("add-")){
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
    handleCanvasClick(event) {
        const node = this.getNodeAt(event.offsetX, event.offsetY);
        if (!node || node.id === gameState.player.location_id) return;
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
        if (gameState.isMoving) {
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

    showContextMenu(event, element){        
        const menu = document.getElementById('context-menu');
        const charId = element.dataset.id;
        const charName = element.querySelector('.p-name').innerText;

        // Наполняем меню кнопками
        menu.innerHTML = `
            <div style="padding: 5px; font-size: 10px; color:#FFD700; border-bottom: 1px solid #222;">${charName}</div>
            <button class="context-menu-btn" data-action="private" data-id="${charId}" data-name="${charName}">📨 Приват</button>
            <button class="context-menu-btn" data-action="challenge" data-id="${charId}" data-name="${charName}">⚔️ Вызвать</button>
            <button class="context-menu-btn" data-action="trade" data-id="${charId}">💰 Торговля</button>
            <button class="context-menu-btn" data-action="info" data-id="${charId}">📜 Инфо</button>
        `;

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