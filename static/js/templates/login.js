export function loginTemplate(error) {
    return `
    <div class="auth-container">
                <h1>'Вход'</h1>
                <div class="auth-form">                    
                    <input type="email" id="email" placeholder="E-mail">
                    <input type="password" id="password" placeholder="Пароль">                  
                    <button id="submitBtn">'Войти'</button>
                </div>
                <p id="toggleMode">'Нет аккаунта? Регистрация'</p>
                <div id="authError">${error}</div>
            </div>    
    `;
}