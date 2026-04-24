export async function apiCall(url, method = "GET", body = null){
    let token = localStorage.getItem("game_token");
    const headers = {
        'Content-Type':'application/json'
    };
    if(token){
        headers['Authorization'] = `Bearer ${token}`;
    }
    const config = {
        method:method,
        headers:headers,
    }
    if(body){
        config.body = JSON.stringify(body);
    }
    let response = await fetch(url, config);
    if(response.status === 401){
        const refreshed = await refresh()
        if(refreshed){
            const newToken = localStorage.getItem("game_token");
            config.headers['Authorization'] = `Bearer${newToken}`;
            return await fetch(url, config);
        }
        else{
            localStorage.removeItem('game_token');
            localStorage.removeItem('refresh_token');
            return response; // Вернет исходный 401, и main.js покажет логин
        }
    }
    return response;
}

async function  refresh(){
    const refreshToken = localStorage.getItem("refresh_token");
    if(!refreshToken) return false;
    const response = await fetch("api/refresh", {
        method:'POST',
        headers:{'Content-Type':'application/json'},
        body:JSON.stringify({'refresh_token':refreshToken})
    });
    if(response.ok){
        const data = await response.json();
        localStorage.setItem("game_token", data.token)
        localStorage.setItem("refresh_token", data.refresh_token)
        return true;
    }
    return false;
}