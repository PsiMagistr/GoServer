import { ui } from "./ui.js";
export const socket_events = {
    room_presence(msg){
         console.log("Игроки в нашей комнате:");         
         ui.renderPlayerList(msg.players)          
    },
    player_joined(msg){
        console.log(`К нам присоединился(лась) ${msg.player.name}`);
        ui.addPlayerToUI(msg.player)
    },
    player_left(msg){
        console.log("+++++++++")
        console.log(msg.player)
        console.log(`Нас покинул(ла) ${msg.player.name}`);        
        ui.removePlayerFromUI(msg.player.id);
    },
}