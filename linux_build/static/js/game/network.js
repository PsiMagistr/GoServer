import { socket_events } from "./socket_events.js";
export const network = {
    socket:null,
    connect(){
        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
        this.socket = new WebSocket(`${protocol}//${window.location.host}/ws`);
        this.socket.onopen = ()=>{
            console.log("Сетевой узел активен.");
        }
        this.socket.onmessage = (event)=>{
            try{
                const msg = JSON.parse(event.data);
                if(socket_events[msg.type]){
                    socket_events[msg.type](msg);
                }
            }catch(e){
                console.error("Ошибка сети:", e)
            }
        }
        this.socket.onclose = ()=>{
            alert("Связь с сервером потеряна.")
        }
    },
    send(data = {}){
        if(this.socket && this.socket.readyState === WebSocket.OPEN){
            this.socket.send(JSON.stringify(data));
        }
    }
}