import { network } from "./network.js";

export const gameActions = {
    moveToNode(targetId) {        
        const packet = {
            "type": "move",
            "target_id": targetId,
        }
        network.send(packet);
    },
    enterPortal(worldId) {        
        const packet = {
            "type": "portal_request",
            "world_id": worldId,
        };
        network.send(packet);
    },
    sendPublicChat(text){
        const packet = {
            type: "chat_msg",
            text: text,
        };
        network.send(packet);
    },
    sendWhisper(recipient, text){
        const packet = {
            type: "private_chat",
            target_name: recipient,
            text: text,
        };
        network.send(packet);
    },
}