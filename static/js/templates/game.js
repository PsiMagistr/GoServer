export const gameTemplate = (char) => {
    return `
        <div class="game-wrapper">
            <div class="game-layout">
                <!-- Левая панель: Списки (фиксированная 180px) -->
                <aside class="side-panel">
                    <div class="panel-section">
                        <h3>Люди</h3>
                        <ul id="players-list" class="game-list">
                            <li class="empty">Пусто</li>                            
                        </ul>
                    </div>
                    <div class="panel-section">
                        <h3>Монстры</h3>
                        <ul id="monsters-list" class="game-list">
                            <li class="empty">Пусто</li>
                        </ul>
                    </div>
                    <div class="panel-section">
                        <h3>Приглашения</h3>
                        <ul id="challenges-list" class="game-list">
                            <li class="empty">Пусто</li>                           
                        </ul>
                    </div>
                    <div class="panel-section">
                        <h3>Персонажи</h3>
                        <ul id="npc-list" class="game-list">
                            <li class="empty">Пусто</li>                                                    
                        </ul>
                    </div>
                    <div class="panel-section">
                        <h3>Миры</h3>
                        <ul id="worlds-list" class="game-list">
                            <li class="empty">Пусто</li>                                                       
                        </ul>
                    </div>
                </aside>

                <!-- Основная зона -->
                <main class="main-screen">
                <div class="hero-nav-group">
                    <div class="top-bar-simple">
                        <span id = "char_name_label" class="char-name-label">Герой: <b>${char.name}</b></span>
                        <span id = "location_label" class="location-label">Локация: <b>${char.location_id || 'Начало'}</b></span>
                    </div>
                    <div class="hero-bar-simple">
                        <span id="stats" class="hero-panel-link">Характеристики</span>
                        <span class="hero-panel-link">Магия</span>
                        <span class="hero-panel-link">Рюкзак</span>
                    </div>
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
                        <!-- Маленькое поле для имени (автоматически заполняется при клике по списку) -->
                        <input type="text" id="chat-to" placeholder="Всем" title="Имя получателя (оставьте пустым для общего чата)" maxlength="20">
                        
                        <!-- Основное поле сообщения -->
                        <input type="text" id="chat-input" placeholder="Введите сообщение..." maxlength="120" autocomplete="off">
                        
                        <button id="chat-send-btn">ОК</button>
                    </div>
                </div>
                </main>
            </div>
            <div id="modal-overlay" class="modal-overlay" style="display:none;">
                <div class="modal-window">
                    <button id="modal-close-btn" class="modal-close">&times;</button>
                    <div id="modal-content">Пространство модального окна</div>
                </div>
            </div>
            <div id="context-menu" class="context-menu" style="display: none;"></div>
        </div>
    `;
};