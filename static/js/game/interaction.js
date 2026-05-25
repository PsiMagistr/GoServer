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
        const worldLink = target.closest(".world-link");
        if(worldLink){
            this.handlePortalClick(worldLink)
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
                console.log(gameState.player)
                if(node.id === gameState.player.location_id) return
                const goToQuestion = confirm(`Вы хотите перейти в ${node.name}`)                
                if(!goToQuestion) return
                const packet = {
                    type:"move",
                    target_id: node.id,
                }               
                network.send(packet);   
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
    },
    handlePortalClick(element){
        if(gameState.isMoving){
            alert("Вы уже в пути")
            return
        }
        const worldId = element.dataset.id; // Берем ID из data-world-id
        //console.log("eweweeeeewe")
        //console.log(element)
        const worldName = "totlhaim";/*element.innerText.replace('🌀', '').trim();*/
        if (confirm(`Вы уверены, что хотите покинуть этот мир и отправиться в "${worldName}"? (Переход займет 200 сек)`)) {
            console.log("Запрос на телепортацию в:", worldId);
            const packet = {
                    type:"portal_request",
                    world_id: worldId, 
            }
            network.send(packet);
        }
    },
}