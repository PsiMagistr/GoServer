export function initGame(username) {
    const app = document.getElementById('app');

    app.innerHTML = `
        <div class="game-ui">
            <p>Игрок: <b>${username}</b></p>
            <canvas id="gameCanvas" width="600" height="400"></canvas>
            <div id="gameLog">Вы вошли в игру.</div>
        </div>
    `;

    const canvas = document.getElementById('gameCanvas');
    const ctx = canvas.getContext('2d');

    // Простейший цикл отрисовки
    function draw() {
        ctx.fillStyle = 'black';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = 'lime';
        ctx.fillRect(280, 180, 40, 40); // Наш "игрок"
        requestAnimationFrame(draw);
    }
    draw();
}