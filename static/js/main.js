import { initAuth } from './auth.js';
import { initGame, showCreateCharachter } from './game.js';
import {apiCall} from "./api.js";

/**
 * Главная точка входа в приложение
 */


const onSuccess = async (username)=>{
    console.log(`Пользователь ${username} вошел. Запускаем игру...`);
    const response = await apiCall("/api/character/me","GET");
    if(response.ok){
        const data = await response.json();
        console.log(data.id)
        initGame(username);
    }
    else{
        const data = await response.json();
        showCreateCharachter()    
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
        const response = await apiCall("api/me")

        if (response.ok) {
            const userData = await response.json();
            console.log("Авторизация подтверждена для:", userData.username);
            onSuccess(userData.username);           
        } else {
            console.log("Сервер отклонил токен. Статус:", response.status);           
            initAuth(onSuccess);
        }
    } catch (e) {
        console.error("Ошибка связи с сервером:", e);
        initAuth(onSuccess);
    }
}
// Запускаем наше приложение
main();