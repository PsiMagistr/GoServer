import { initAuth } from './auth.js';
import {initGame} from './game/game.js';
import {showCreateCharacter} from './char_creation.js';
import {apiCall} from "./api.js";

/**
 * Главная точка входа в приложение
 */


const onSuccess = async (username)=>{
    try{
        console.log(`Пользователь ${username} вошел. Запускаем игру...`);
        const response = await apiCall("/api/character/me","GET");
        if(response.ok){
            const char = await response.json();            
            initGame(char);
        }
        else if(response.status == 404){       
            showCreateCharacter();
            console.log("Аккаунт не имеет персонажа. Создайте его.");              
        }
        else{
           console.error("Неизвестная ошибка")
        }   
    }catch(e){
          console.error("Критическая ошибка (сеть):", e);    
    }
    
    // 2. После успешного входа инициализируем игровой движок
    // Передаем имя пользователя в игру, чтобы отобразить его на экране
    

    // Здесь в будущем мы добавим создание WebSocket соединения
    // setupNetwork(username);
}

async function main() {
    console.log("Приложение запущено. Проверка авторизации...");
      
    try {
        // ИСПОЛЬЗУЕМ ЧИСТЫЙ ПУТЬ БЕЗ СЛЭША В КОНЦЕ
        const response = await apiCall("/api/me")

        if (response.ok) {
            const userData = await response.json();
            console.log("Авторизация подтверждена для:", userData.username);
            onSuccess(userData.username);           
        } else {
            console.log("Сервер отклонил токен. Статус:", response.status);           
            initAuth(onSuccess);
        }
    } catch (e) {
        if (e.message === "RATE_LIMIT_REACHED") {            
            console.warn("Ждем восстановления лимита...");
            return;
        }
        console.error("Ошибка связи с сервером:", e);
        initAuth(onSuccess);
    }
}
// Запускаем наше приложение
main();