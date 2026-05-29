import { gameState } from "./game.js";
import { gameActions } from "./actions.js";
import { engine } from "./engine.js";
const CLICK_RADIUS = 20;
const CLICK_RADIUS_SQ = CLICK_RADIUS * CLICK_RADIUS;
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
        const id = event.target.id;
        if (id == "gameCanvas") {
            this.handleCanvasClick(event);
            return;
        }
        if (id == "chat-send-btn") {
            this.sendChat()
            return;
        }
        const worldLink = target.closest(".world-link");
        if (worldLink) {
            this.handlePortalClick(worldLink)
        }
    },
    handleMouseMove(event) {
        const node = this.getNodeAt(event.offsetX, event.offsetY);        
        const found = node ? node.id : null;
        if (gameState.hoveredNodeId !== found) {
            gameState.hoveredNodeId = found;
            event.target.style.cursor = found ? 'pointer' : 'default';
        }
    },
    handleGlobalKeyPress(event) {
        if (event.key === 'Enter' && event.target.id === 'chat-input') {
            this.sendChat()
        }
    },

    handleCanvasClick(event, canvas) {        
        const node = this.getNodeAt(event.offsetX, event.offsetY);
        if (!node || node.id === gameState.player.location_id) return;
        const goToQuestion = confirm(`Вы хотите перейти в ${node.name}`)
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
        const worldName = "totlhaim";/*element.innerText.replace('🌀', '').trim();*/
        const goToQuestion = confirm(`Вы уверены, что хотите покинуть этот мир и отправиться в "${worldName}"? (Переход займет 200 сек)`)
        if (!goToQuestion) return
        console.log("Запрос на телепортацию в:", worldId);
        gameActions.enterPortal(worldId)
    },
}