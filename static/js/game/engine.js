import { gameState } from "./game.js";
export const engine = {
    canvas: null,
    ctx: null,
    loopId: null,
    init(canvasElement) {
        this.canvas = canvasElement;
        this.ctx = canvasElement.getContext('2d');
        this.start();
    },
    start() {
        if (this.loopId) this.stop();
        const render = () => {
            this.draw()
            this.loopId = window.requestAnimationFrame(render)
        }
        render();
    },
    draw() {
        this.ctx.fillStyle = 'black';
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        if (!gameState.player) {
            this.ctx.fillStyle = 'lime';
            this.ctx.fillRect(280, 180, 40, 40);
            return;
        }
        this.ctx.fillStyle = 'white';
        this.ctx.fillText(`${gameState.player.name} HP: ${gameState.player.hp}`, 280, 170);

    },
    stop() {
        cancelAnimationFrame(this.loopId);
    },

    async loaderAssets(assetsMap) {
        const keys = Object.keys(assetsMap);
        const promises = [];
        function load(path) {
            const img = new Image();
            return new Promise((resolve, reject) => {
                img.addEventListener("load", () => {
                    resolve(img);
                })
                img.addEventListener("error", () => {
                    reject("Ошибка загрузки");
                })
                img.src = path;
            })
        }

        for (const key in assetsMap) {
            promises.push(load(assetsMap[key]));
        }
        const resolvedPromises = await Promise.all(promises);        
        const library = {}
        keys.forEach((key, index) => {
            library[key] = resolvedPromises[index];
        });
        return library;
    }
}