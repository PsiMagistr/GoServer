export const statsModalTemplate = (draft) => {    
    const canClickSave = draft.hasChanges && !draft.isLoading ? "" : "disabled";
    const canClickAdd = draft.free_points > 0 && !draft.isLoading;
    const labelSave = !draft.isLoading ? "Сохранить":"Сохранение";   
    const errorBlock = draft.error 
        ? `<div class="modal-error-msg">Ошибка сохранения в БД: ${draft.error}</div>` 
        : '';
    const addBtn = (statName)=>{
        let statBtnFlag = ""       
        if(!canClickAdd){            
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
                <span>Выносливость:<b>${draft.vitality}</b></span>
                ${addBtn("vitality")}
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
            <br>
            <p style="font-size: 12px; color: #666;">Нажмите на крестик, чтобы закрыть</p>
            ${errorBlock}            
            <div class="modal-actions">
                <button id="save-stats-btn" class="main-btn" ${canClickSave}>${labelSave}</button>
            </div>
        </div>
    `;
}