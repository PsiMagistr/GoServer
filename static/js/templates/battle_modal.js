export const battleModalTemplate = (data) => {
    const you = data.you;
    const opponent = data.opponent;
    const round = data.round;

    // Генерируем 5 слотов (2 защиты, 3 атаки)
    let slotsHtml = '';
    for (let i = 1; i <= 5; i++) {
        //const isShield = i <= 2;
        slotsHtml += `
            <div class="battle-slot" id="slot-${i}" data-slot-index="${i}">
                <div id="label-slot-${i}" class="slot-type-label">Пусто</div>
                <div class="slot-icon" id="slot-icon-${i}"></div>
            </div>
        `;
    }

    

    return `
        <div class="battle-interface">
            <!-- ВЕРХНЯЯ ПАНЕЛЬ: Таймер и Имена -->
            <div class="battle-header">
                <div class="battle-title">Поединок: <b>${you.name}</b> vs <b>${opponent.name}</b> Раунд:${round}</div>
                <div id="battle-timer" class="battle-timer">${data.time_left || 30}s</div>
            </div>

            <!-- ЭКРАН БОЯ (Канвас) -->
            <div class="battle-screen">
                <canvas id="battleCanvas" width="650" height="200"></canvas>
            </div>

            <!-- ВЫБРАННЫЕ СЛОТЫ -->
            <div class="battle-slots-row">
                ${slotsHtml}
            </div>

            <!-- КНИГА ЗАКЛИНАНИЙ (Две колонки) -->
            <div class="spells-library">
             <!-- Колонка Атаки -->
                <div class="spells-column">
                        <h4>Боевая магия</h4>
                        <div id="attack-spells" class="spells-list">
                            <div class="spell-item attack fire" data-type="atk" data-element="fire">🔥 Огненный шар</div>
                            <div class="spell-item attack water" data-type="atk" data-element="water">💧 Ледяная стрела</div>
                            <div class="spell-item attack air" data-type="atk" data-element="air">🌬️ Молния</div>
                            <div class="spell-item attack earth" data-type="atk" data-element="earth">🌱 Каменный шип</div>
                        </div>
                    </div>
                <!-- Колонка Защиты -->
                <div class="spells-column">
                    <h4>Щиты</h4>
                    <div id="defense-spells" class="spells-list">
                        <div class="spell-item shield fire" data-type="def" data-element="fire">🔥 Щит Огня</div>
                        <div class="spell-item shield water" data-type="def" data-element="water">💧 Щит Воды</div>
                        <div class="spell-item shield air" data-type="def" data-element="air">🌬️ Щит Воздуха</div>
                        <div class="spell-item shield earth" data-type="def" data-element="earth">🌱 Щит Земли</div>
                    </div>
                </div>             
               
            </div>

            <!-- ЛОГ БОЯ -->
            <div id="battle-log" class="battle-log">
                <div class="log-entry sys">Бой начался! Соберите комбинацию стихий.</div>
            </div>

            <!-- КНОПКИ ДЕЙСТВИЙ -->
            <div class="battle-footer-actions">
                <button id="btn-submit-turn" class="battle-btn btn-confirm" disabled>ПОДТВЕРДИТЬ ХОД</button>
                <button id="btn-surrender" class="battle-btn btn-surrender">СДАТЬСЯ</button>
            </div>
        </div>
    `;
};