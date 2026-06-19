export const battleModalTemplate = (data) => {
    const p = data.player;
    const o = data.opponent;

    // Генерируем 5 пустых слотов
    let slotsHtml = '';
    for (let i = 1; i <= 5; i++) {
        const type = i <= 2 ? 'Щит' : 'Атака';
        slotsHtml += `
            <div class="battle-slot" data-slot="${i}" id="slot-${i}">
                <span class="slot-type">${type}</span>
                <div class="slot-icon">?</div>
            </div>
        `;
    }

    return `
        <div class="battle-container">
            <div class="battle-header">
                <div class="fighter-info">
                    <b>${p.name}</b> vs <b>${o.name}</b>
                </div>
                <div class="battle-timer" id="battle-timer">${data.timeLeft}s</div>
            </div>

            <div class="battle-arena">
                <canvas id="battleCanvas" width="500" height="200"></canvas>
            </div>

            <!-- ПАНЕЛЬ ВЫБОРА МАГИИ -->
            <div class="magic-selector">
                <div class="elements-grid">
                    <button class="el-btn fire" data-element="fire">🔥 Огонь</button>
                    <button class="el-btn water" data-element="water">💧 Вода</button>
                    <button class="el-btn air" data-element="air">🌬️ Воздух</button>
                    <button class="el-btn earth" data-element="earth">🌱 Земля</button>
                </div>
                
                <div class="selected-slots">
                    ${slotsHtml}
                </div>
            </div>

            <div id="battle-log" class="battle-log">
                <div class="log-entry">Бой начался! Выберите 2 щита и 3 атаки.</div>
            </div>

            <div class="battle-actions">
                <button id="btn-submit-turn" class="main-btn" disabled>ПОДТВЕРДИТЬ ХОД</button>
                <button id="btn-surrender" class="btn-surrender">Сдаться</button>
            </div>
        </div>
    `;
};