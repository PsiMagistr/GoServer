export const battleModalTemplate = (data) => {
    // data содержит: { player, opponent, timeLeft }
    /*const p = data.player;
    const o = data.opponent;*/

    // Генерируем 5 слотов (2 щита, 3 атаки)
    let slotsHtml = '';
    for (let i = 1; i <= 5; i++) {
        const isShield = i <= 2;
        slotsHtml += `
            <div class="magic-slot ${isShield ? 'slot-def' : 'slot-atk'}" id="slot-${i}" data-slot="${i}">
                <div class="slot-label">${isShield ? 'ЩИТ' : 'УДАР'}</div>
                <div class="slot-icon">?</div>
            </div>
        `;
    }

    return `
        <div class="modal-content">
            <!-- ВЕРХНЯЯ ПАНЕЛЬ: Таймер и Статус -->
            <div class="battle-top-bar">
                <div class="battle-title">Поединок: Вася vs Петя</div>
                <div id="battle-timer" class="battle-timer">30s</div>
            </div>

            <!-- ЭКРАН БОЯ: Канвас (Аватары и Шкалы рисует Engine) -->
            <div class="battle-screen">
                <canvas id="battleCanvas" width="650" height="250"></canvas>
            </div>

            <!-- ЦЕНТРАЛЬНАЯ ПАНЕЛЬ: Выбор стихий и Слоты -->
            <div class="battle-management">
                <div class="element-picker">
                    <button class="el-btn fire" data-element="fire" title="Огонь">🔥</button>
                    <button class="el-btn water" data-element="water" title="Вода">💧</button>
                    <button class="el-btn air" data-element="air" title="Воздух">🌬️</button>
                    <button class="el-btn earth" data-element="earth" title="Земля">🌱</button>
                </div>

                <div class="magic-slots-container">
                    ${slotsHtml}
                </div>
            </div>

            <!-- НИЖНЯЯ ПАНЕЛЬ: Лог и Кнопки -->
            <div class="battle-footer">
                <div id="battle-log" class="battle-log">
                    <div class="log-entry sys">Бой начался! Выберите заклинания для хода.</div>
                </div>
                
                <div class="battle-actions">
                    <button id="btn-submit-turn" class="main-btn btn-confirm" disabled>СДЕЛАТЬ ХОД</button>
                    <button id="btn-surrender" class="main-btn btn-quit">СДАТЬСЯ</button>
                </div>
            </div>
        </div>
    `;
};