export function regTemplate(error){
    return `
        <div class="auth-container">
                <h1>'Регистрация'</h1>
                <div class="auth-form">
                    <input type="text" id="username" placeholder="Имя владельца" value="Павел">
                    <input type="email" id="email" placeholder="E-mail" value="chago@test.ru">
                    <input type="password" id="password" placeholder="Пароль" value="12345">
                    <input type="password" id="confirm_password" placeholder="Повтор пароля" value="12345">
                    <button id="submitBtn">'Создать аккаунт'</button>
                </div>
                <p id="toggleMode">'Есть аккаунт? Войти'</p>
                <div id="authError">${error}</div>
            </div>    
    `;
}
