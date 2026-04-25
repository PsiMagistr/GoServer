import { initAuth } from './auth.js';
import { initGame } from './game.js';
import {apiCall} from "./api.js";

/**
 * Главная точка входа в приложение
 */


const onSuccess = (username)=>{
    console.log(`Пользователь ${username} вошел. Запускаем игру...`);
    // 2. После успешного входа инициализируем игровой движок
    // Передаем имя пользователя в игру, чтобы отобразить его на экране
    initGame(username);

    // Здесь в будущем мы добавим создание WebSocket соединения
    // setupNetwork(username);
}

async function main() {
    console.log("Приложение запущено. Проверка авторизации...");
    const token = localStorage.getItem("game_token");

    if (!token) {
        console.log("Токена нет в localStorage");
        initAuth(onSuccess);
        return;
    }

    try {
        // ИСПОЛЬЗУЕМ ЧИСТЫЙ ПУТЬ БЕЗ СЛЭША В КОНЦЕ
        const response = await apiCall("api/me")

        if (response.ok) {
            const userData = await response.json();
            console.log("Авторизация подтверждена для:", userData.username);
            onSuccess(userData.username);
        } else {
            console.log("Сервер отклонил токен. Статус:", response.status);
            localStorage.removeItem('game_token');
            initAuth(onSuccess);
        }
    } catch (e) {
        console.error("Ошибка связи с сервером:", e);
        initAuth(onSuccess);
    }
}
// Запускаем наше приложение
main();