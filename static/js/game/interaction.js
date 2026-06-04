import { gameState } from "./game.js";
import { gameActions } from "./actions.js";
import { engine } from "./engine.js";
import {modalManager} from "./modalManager.js";
import {statsModalTemplate} from "../templates/stats_modal.js";
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
        modalManager.show(statsModalTemplate, gameState.player)
    },
    "modal-close-btn":(obj, event)=>{
        modalManager.hide();
    }
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
        if(id.startsWith("add-")){
            alert(id);
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
}