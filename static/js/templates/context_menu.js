export const contextNenu = (charID, charName)=>{
    return `
            <div style="padding: 5px; font-size: 10px; color:#FFD700; border-bottom: 1px solid #222;">${charName}</div>
            <button class="context-menu-btn" data-action="private" data-id="${charID}" data-name="${charName}">📨 Приват</button>
            <button class="context-menu-btn" data-action="challenge" data-id="${charID}" data-name="${charName}">⚔️ Вызвать</button>
            <button class="context-menu-btn" data-action="trade" data-id="${charID}">💰 Торговля</button>
            <button class="context-menu-btn" data-action="info" data-id="${charID}">📜 Инфо</button>
        `;
}