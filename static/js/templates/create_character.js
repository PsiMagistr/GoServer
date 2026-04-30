export const CreateCharacterTemplate = (state) => {
    const avatars = [];    
    // Цикл для создания сетки 3x3 (9 штук)
    for (let i = 1; i <= 9; i++) {
        // Путь теперь учитывает папку пола: /assets/avatars/male/1.png
        const imagePath = `/assets/avatars/${state.gender}/${i}.png`;
        const isSelected = state.selectedAvatar === i ? 'selected' : '';        
        avatars.push(`
            <div class="avatar-item ${isSelected}" data-id="${i}">
                <img src="${imagePath}" alt="Avatar ${i}">
            </div>
        `);
    }

    return `
        <div class="creation-container">
            <h2>Создание героя</h2>
            
            <div class="form-group">
                <input type="text" id="charName" placeholder="Имя персонажа" value="${state.name || ''}">
            </div>

            <div class="gender-toggle">
                <button id="setMale" class="${state.gender === 'male' ? 'active' : ''}">Мужчина</button>
                <button id="setFemale" class="${state.gender === 'female' ? 'active' : ''}">Женщина</button>
            </div>

            <div class="avatar-grid">
                ${avatars.join('')}
            </div>

            <div id="charError" class="error-msg">${state.error || ''}</div>
            
            <button id="createConfirmBtn" class="btn-hero">Войти в мир</button>
        </div>
    `;
};