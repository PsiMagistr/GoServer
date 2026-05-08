export const gameTemplate = (char) => {
    return `
        <div class="game-layout">
            <!-- Левая колонка со списками -->
            <aside class="side-panel">
                <div class="panel-section">
                    <h3>Игроки</h3>
                    <ul id="players-list" class="game-list">
                        <li class="empty">Загрузка...</li>
                    </ul>
                </div>
                
                <div class="panel-section">
                    <h3>Монстры</h3>
                    <ul id="monsters-list" class="game-list">
                        <li class="empty">В локации пусто</li>
                    </ul>
                </div>

                <div class="panel-section">
                    <h3>Здания / Переходы</h3>
                    <ul id="buildings-list" class="game-list">
                        <li>Таверна</li>
                        <li>Арена</li>
                    </ul>
                </div>
            </aside>

            <!-- Центральная часть: Канвас и Чат -->
            <main class="main-screen">
                <div class="canvas-container">
                    <div class="char-header">Персонаж: <b>${char.name}</b></div>
                    <canvas id="gameCanvas" width="600" height="400"></canvas>
                </div>

                <div class="chat-area">
                    <div id="chat-messages" class="chat-messages">
                        <div class="sys-msg">Добро пожаловать в мир, ${char.name}!</div>
                    </div>
                    <div class="chat-input-group">
                        <input type="text" id="chat-input" placeholder="Введите сообщение..." autocomplete="off">
                        <button id="chat-send-btn">ОК</button>
                    </div>
                </div>
            </main>
        </div>
    `;
};