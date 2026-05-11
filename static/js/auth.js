import { regTemplate } from './templates/registration.js';
import { loginTemplate } from './templates/login.js';
import {utils} from './utils/utils_functions.js';
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
                console.log("ENTER2")
                const params = {
                    email:"#email",
                    password:"#password",
                }
                const configData = utils.getElementsBySelectors(params);
                const user = await authService.login(
                    configData.email.value,
                    configData.password.value
                );
                onSuccess(user.username); // Вход в аккаунт.
            }
            else{// Регистрация.
                 const params = {
                    username:"#username",
                    email:"#email",                   
                    password:"#password",
                    confirm_password:"#confirm_password",
                }
                const userData = utils.getElementsBySelectors(params); 
                console.log("===========")
                const schema = utils.getValuesBySchema(userData)
                console.log(schema)            
                await authService.register(schema)                 
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