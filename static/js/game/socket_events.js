export const socket_events = {
    room_presence(msg){
        console.log("Игроки в нашей комнате:");            
            for(const player of msg.players){
                console.log(player.name);
            }      
    },
    player_joined(msg){
        console.log(`К нам присоединился(лась) ${msg.player.name}`);
    },
    player_left(msg){
        console.log(`Нас покинул(ла) ${msg.player.name}`);
    },
}