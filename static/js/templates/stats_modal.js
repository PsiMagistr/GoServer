export const statsModalTemplate = (player) => {
    const hasPoints = player.free_points > 0;
    const addBtn = (statName)=>{
        return hasPoints ? `<button class="add-stat-btn" id="add-${statName}">+</button>`:""
    }
    return `
     <div class="stats-view">
            <h2>Характеристики героя</h2>
             <div class="free-points-block">
                Свободных очков: <b id="free-points-count">${player.free_points}</b>
            </div>
            <div class="stat-line">
                <span>Сила: <b>${player.strength}</b></span>
                ${addBtn("strength")}
            </div>
            <div class="stat-line">
                <span>Ловкость:<b>${player.agility}</b></span>
                ${addBtn("agility")}
            </div>
            <div class="stat-line">
                <span>Интуиция:<b>${player.intuition}</b></span>
                ${addBtn("intuition")}    
            </div>
            <div class="stat-line">
                <span>Мудрость:<b>${player.wisdom}</b></span>
                ${addBtn("wisdom")}    
            </div>
            <div class="stat-line">
                <span>Выносливость:<b>${player.vitality}</b></span>
                ${addBtn("vitality")}
            </div>            
            <br>
            <p style="font-size: 12px; color: #666;">Нажмите на крестик, чтобы закрыть</p>
            <div class="modal-actions">
                <button id="save-stats-btn" class="main-btn">Сохранить</button>
            </div>
        </div>
    `;
}