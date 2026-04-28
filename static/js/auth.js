import { regTemplate } from './templates/registration.js';
import { loginTemplate } from './templates/login.js';
import {getConfigData} from './utils/utils_functions.js';
import { authService } from './services/authService.js';

export function initAuth(onSuccess) {
    const app = document.getElementById('app');
    let state = {
        isLoginMode:true,
        error:"",
    }
    app.onclick = (event)=>{
        const id = event.target.id;
        if(id === "toggleMode"){
            state.isLoginMode = !state.isLoginMode;
            state.error = "";
            render();
        }
        else if(id === "submitBtn"){
            handleSubmit(onSuccess);
        }
    }
    const render = () => {
        const template = state.isLoginMode ? loginTemplate(state.error) : regTemplate(state.error);
        app.innerHTML = template;
    };

    const handleSubmit = async (callback) => {
        
        /////
        try{
            if(state.isLoginMode){//Вход.
                const configData = getConfigData(["email", "password"]);
                const user = await authService.login(configData.email, configData.password);
                onSuccess(user.username);
            }
            else{// Регистрация.
                const userData = getConfigData(['username', 'email', 'password', 'confirm_password']);
                    await authService.register(userData);
                    state.isLoginMode = true;
                    state.error = "";
                    render();               
                
            }
        }
        catch(e){
            state.error = "Ошибка: " + e.message;
            render();
        }

    }
    render(); // Первый запуск
}