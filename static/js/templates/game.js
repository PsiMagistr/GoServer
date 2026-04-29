export function gameTemplate(username) {
    return `
        <div class="game-ui">
            <p>Игрок: <b>${username}</b></p>
            <canvas id="gameCanvas" width="600" height="400"></canvas>
            <div id="gameLog">Вы вошли в игру.</div>
        </div>
    `;
}