//Отправление пакетов с данными на сервер.
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
    sendPublicChat(text) {
        const packet = {
            type: "chat_msg",
            text: text,
        };
        network.send(packet);
    },
    sendWhisper(recipient, text) {
        const packet = {
            type: "private_chat",
            target_name: recipient,
            text: text,
        };
        network.send(packet);
    },
    sendPlayerStats(stats) {
        const packet = {
            type: "commit_stats",
            stats: stats,
        }
        network.send(packet)
    },
    sendBattleChallenge(id) {
        const packet = {
            type: "battle_challenge",
            target_id: parseInt(id),
        }
        network.send(packet)
    },
    acceptBattle(senderId) {
        const packet = {
            type: "battle_accept",
            sender_id: parseInt(senderId),
        }
        network.send(packet);
    },
    sendBattleTurn(battleId, round, spells) {
        const packet = {
            type: "battle_turn",
            battle_id: battleId,
            spells: spells,
            round: round,
        }
        network.send(packet)
    },
    declineChallenge(senderId) {
        const packet = {
            type: "battle_decline",
            sender_id: Number(senderId)
        }
        network.send(packet)
    },
    surrender(){        
        const packet = {
            type: "battle_surrender"
        }
        network.send(packet);
    }
}