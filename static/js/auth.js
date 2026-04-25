import { regTemplate } from './templates/registration.js';
import { loginTemplate } from './templates/login.js';
import {apiCall} from "./api.js";

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
        const getConfigData = (params)=>{
            const body = {}
            for(let param of params){
                const elem = document.querySelector(`#${param}`)
                if(elem){
                    body[param] = elem.value
                }
                else {
                    throw new Error("Неверное значение идентификатора поля формы.")
                }
            }
            return body
        }
        /////
        try{
            if(state.isLoginMode){//Вход.
                const body = getConfigData(["email","password"]);
                const response = await fetch("api/login", {
                    method:"POST",
                    headers:{"Content-Type":"application/json"},
                    body:JSON.stringify(body),
                });
                if(response.ok){
                    const data = await response.json();
                    localStorage.setItem("game_token", data.token)
                    localStorage.setItem("refresh_token", data.refresh_token)
                    console.log("Успешный вход, оба токена сохранены.");
                    callback(data.username);
                }
                else{
                    state.error = await response.text();
                    render();
                }

            }
            else{// Регистрация.
                const body = getConfigData(["username","email","password","confirm_password"]);
                const response = await apiCall("/api/register", "POST", body);
                if(response.ok){
                    state.isLoginMode = true;
                    state.error = "";
                    render();
                }
                else{
                    state.error = await response.text();
                    render();
                }
            }
        }
        catch(e){
            state.error = "Ошибка: " + e.message;
            render();
        }

    }
    render(); // Первый запуск
}