import { gameState } from "./game.js";
import { network } from "./network.js";
import { engine } from "./engine.js";
export const interaction = {
    app: null,
    init(appElement) {
        this.app = appElement;
        this.app.onclick = this.handleGlobalClick.bind(this);
        this.app.onmousemove = this.handleMouseMove.bind(this);
        this.app.onkeypress = this.handleGlobalKeyPress.bind(this);
    },
    handleGlobalClick(event) {
        const target = event.target;
        const id = event.target.id;
        if (id == "gameCanvas") {
            this.handleCanvasClick(event);
            return;
        }
        if (id == "chat-send-btn") {
            this.sendMessage();
            return;
        }
    },
    handleMouseMove(event) {
        const points = gameState.world?.points;
        if (!points) return
        let found = null;
        const radius = 20;
        for (const id in points) {
            const node = points[id];
            const dx = event.offsetX - node.x;
            const dy = event.offsetY - node.y;
            if (Math.sqrt(dx * dx + dy * dy) < radius) {
                 found = id;
                 break;
            }
        }
        if (gameState.hoveredNodeId !== found) {
            gameState.hoveredNodeId = found;
            event.target.style.cursor = found ? 'pointer' : 'default';
        }
    },
    handleGlobalKeyPress(event) {
        if (event.key === 'Enter' && event.target.id === 'chat-input') {
            this.sendMessage();
        }
    },
    handleCanvasClick(event, canvas) {
        const points = gameState.world?.points;
        if (!points) return
        const radius = 20;
        for (const id in points) {
            const node = points[id];
            const dx = event.offsetX - node.x;
            const dy = event.offsetY - node.y;
            if (Math.sqrt(dx * dx + dy * dy) < radius) {
                alert("3333");
                break;
            }
        }
    },
    sendMessage(){
        const chatInput = this.app.querySelector("#chat-input");
        const text = chatInput.value.trim();
        if (text != "" && network.socket){
            const packet = {
                type:"chat_msg",
                text:text,
            }
            network.send(packet);
            chatInput.value = "";
        }        
    }
}