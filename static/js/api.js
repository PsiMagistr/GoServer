export async function apiCall(url, method = "GET", body = null) {
    const headers = {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest' // ОБЯЗАТЕЛЬНО ДОБАВИТЬ ЭТО
    };
    const config = {
        method: method,
        headers: headers,
    }
    if (body) {
        config.body = JSON.stringify(body);
    }
    let response = await fetch(url, config);
    if (response.status === 401) {
        const refreshed = await fetch("/api/refresh", {method: "POST"});
        if (refreshed.ok) {
            // Если рефреш прошел, повторяем исходный запрос
            return await fetch(url, config);
        }        
    }
    else if(response.status === 429){ 
        alert("Сервер перегружен. Пожалуйста, не обновляйте страницу так часто.");
        // Возвращаем специальную ошибку, чтобы main.js не рисовал форму логина
        throw new Error("RATE_LIMIT_REACHED");

    }
    return response
}

