export const statsModalTemplate = (draft) => {
    const hasPoints = draft.free_points > 0;
    const saveBtnFlag = draft.hasChanges ? "" : "disabled";
    const addBtn = (statName)=>{
        let statBtnFlag = ""       
        if(!hasPoints){            
            statBtnFlag = "disabled"
        }       
        const button = `<button class="add-stat-btn" ${statBtnFlag} data-state_name="${statName}" id="add-${statName}">+</button>`
        return button        
    }
    return `
     <div class="stats-view">
            <h2>Характеристики героя</h2>
             <div class="free-points-block">
                Свободных очков: <b id="free-points-count">${draft.free_points}</b>
            </div>
            <div class="stat-line">
                <span>Сила: <b>${draft.strength}</b></span>
                ${addBtn("strength")}
            </div>
            <div class="stat-line">
                <span>Ловкость:<b>${draft.agility}</b></span>
                ${addBtn("agility")}
            </div>
            <div class="stat-line">
                <span>Интуиция:<b>${draft.intuition}</b></span>
                ${addBtn("intuition")}    
            </div>
            <div class="stat-line">
                <span>Мудрость:<b>${draft.wisdom}</b></span>
                ${addBtn("wisdom")}    
            </div>
            <div class="stat-line">
                <span>Обaяние:<b>${draft.charm}</b></span>
                ${addBtn("charm")}
            </div> 
            <div class="stat-line">
                <span>Выносливость:<b>${draft.vitality}</b></span>
                ${addBtn("vitality")}
            </div>                        
            <br>
            <p style="font-size: 12px; color: #666;">Нажмите на крестик, чтобы закрыть</p>
            <div class="modal-actions">
                <button id="save-stats-btn" class="main-btn" ${saveBtnFlag}>Сохранить</button>
            </div>
        </div>
    `;
}