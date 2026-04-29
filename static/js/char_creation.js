import { CreateCharacterTemplate } from "./templates/create_character.js";
import { apiCall } from "./api.js";
export function showCreateCharacter(){
    const app = document.getElementById('app');
    app.innerHTML = CreateCharacterTemplate();    
}