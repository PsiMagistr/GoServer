import { CreateCharacterTemplate } from "./templates/create_character.js";
import { apiCall } from "./api.js";
import { characterService } from "./services/characterService.js";
export function showCreateCharacter(){
    const app = document.getElementById('app');
    const state = {
        name: "",
        gender: "male",
        selectedAvatar: 1,
        error: ""
    };
    const syncInputToState = () => {
        const input = app.querySelector("#charName");
        if (input) state.name = input.value;
    };    
    const render = ()=>{
       app.innerHTML = CreateCharacterTemplate(state);   
    }    
    const commands = {
        setMale(){
            syncInputToState();                        
            state.gender = 'male';            
            render();
        },
        setFemale(){
            syncInputToState();               
            state.gender = 'female';            
            render();
        },
        async createConfirmBtn(){
            try{
                syncInputToState();
                await characterService.create(/*Тут будут данные из state и еще кое что */);                
            }
            catch(e){
                const errMessage = `Ошибка сети: ${e}`
                console.log(errMessage);
                state.error = errMessage;
                render();
            }                      
            
        }        
    }   
    app.onclick = async (event) => {
         const avatarItem = event.target.closest('.avatar-item');
         const id = event.target.id;                               
         if(commands[id]){
            await commands[id]();
         }
         else if(avatarItem){
            syncInputToState();
            state.selectedAvatar = parseInt(avatarItem.dataset.id);
            render();              
         }         
    }
    render();    
}