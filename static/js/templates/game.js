export const gameTemplate = (char) => {
    return `
        <div class="game-wrapper">
            <div class="game-layout">
                <!-- Левая панель: Списки (фиксированная 180px) -->
                <aside class="side-panel">
                    <div class="panel-section">
                        <h3>Люди</h3>
                        <ul id="players-list" class="game-list">
                            <li class="empty">Загрузка...</li>
                            <li class="player-link">Вика</li>
                        </ul>
                    </div>
                    <div class="panel-section">
                        <h3>Монстры</h3>
                        <ul id="monsters-list" class="game-list">
                            <li class="empty">Пусто</li>
                        </ul>
                    </div>
                    <div class="panel-section">
                        <h3>Навигация</h3>
                        <ul id="buildings-list" class="game-list">
                            <li class="location-link" data-loc="tavern">Таверна</li>
                            <li class="location-link" data-loc="arena">Арена</li>
                        </ul>
                    </div>
                </aside>

                <!-- Основная зона -->
                <main class="main-screen">
                    <div class="top-bar-simple">
                        <span id = "char_name_label" class="char-name-label">Герой: <b>${char.name}</b></span>
                        <span id = "location_label" class="location-label">Локация: <b>${char.location_id || 'Начало'}</b></span>
                    </div>
                    <div class="canvas-container">
                        <div id="move-overlay">
                            <p>Идем в <span class="target-name"></span>...</p>
                            <span class="timer-count">0</span>
                        </div>
                        <canvas id="gameCanvas" width="600" height="400"></canvas>
                    </div>
                    
                    <div class="chat-section">
                        <div id="chat-messages" class="chat-window"></div>
                        <div class="chat-controls">
                            <input type="text" id="chat-input" placeholder="Введите сообщение..." maxlength="120" autocomplete="off">
                            <button id="chat-send-btn">Чат</button>
                        </div>
                    </div>
                </main>
            </div>
        </div>
    `;
};