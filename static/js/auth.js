import { regTemplate } from './templates/registration.js';
import { loginTemplate } from './templates/login.js';
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
        const getConfig = (params)=>{
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
            return {
                method:"POST",
                headers:{"Content-Type":"application/json"},
                body:JSON.stringify(body),
            }
        }
        if(state.isLoginMode){
            const config = getConfig(["email","password"])
            const response = await fetch("/api/login", config)
            if(response.ok){
                const message = await response.json();
                console.log(message)
                localStorage.setItem('game_token', message.token)
                callback(message.username);
            }
            else{
                state.error = await response.text();
                render()
            }
        }
        else{// регистрация.
           const config = getConfig(["username","email","password","confirm_password"])
           const response = await fetch("/api/register", config)
           if(response.ok){
              const message = await response.json();
              state.isLoginMode = true
              render();
           }
           else{
               state.error = await response.text();
               render();
           }
        }
    };
    render(); // Первый запуск
}