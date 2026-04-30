import { CreateCharacterTemplate } from "./templates/create_character.js";
import { apiCall } from "./api.js";
export function showCreateCharacter(){
    const render = ()=>{
       app.innerHTML = CreateCharacterTemplate(state);   
    }
    const state = {
        name: "",
        gender: "female",
        selectedAvatar: 1,
        error: ""
    };
    const commands = {
        setMale:(id)=>{
            state.gender = 'male';
            render();
        },
        setFemale:(id)=>{
             state.gender = 'female';
             render();
        }
    }    
    const app = document.getElementById('app');
    app.onclick = async (event) => {
         const id = event.target.id;        
         if(commands[id]){
            commands[id](id);
         }
         
    }
    render();    
}