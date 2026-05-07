export function loginTemplate(error) {
    return `
    <div class="auth-container">
                <h1>'Вход'</h1>
                <div class="auth-form">                    
                    <input type="email" id="email" placeholder="E-mail" value="chago@test.ru">
                    <input type="password" id="password" placeholder="Пароль" value="12345">                  
                    <button id="submitBtn">'Войти'</button>
                </div>
                <p id="toggleMode">'Нет аккаунта? Регистрация'</p>
                <div id="authError">${error}</div>
            </div>    
    `;
}