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
    sendPlayerStats(draft) {
    // Формируем "худой" пакет - только то, что сервер должен записать
    const statsOnly = {
        strength:  draft.strength,
        agility:   draft.agility,
        intuition: draft.intuition,
        vitality:  draft.vitality,
        charm:     draft.charm,
        wisdom:    draft.wisdom
    };

    const packet = {
        type: "commit_stats",
        stats: statsOnly, // Теперь тут нет тяжелых массивов с магией
    }      
    network.send(packet);

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