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
        } else {
            // Если и рефреш сдох — на выход
            window.location.reload();
        }
        return response
    }
    return response
}

