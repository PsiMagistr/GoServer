import { gameState } from "./game.js";
import { graphs } from "./graphic_settings.js";
export const engine = {
    canvas: null,
    ctx: null,
    loopId: null,
    images:null,
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
        if(!this.images) return;
        this.drawMap();
        this.drawCityNodes();
        this.drawPlayer();      

    },
    stop() {
        cancelAnimationFrame(this.loopId);
    },
    drawMap(){
        this.ctx.drawImage(this.images.map, 0, 0, this.canvas.width, this.canvas.height);
    },
    drawPlayer(){
         const char = gameState.player;
         const healthPercent = Math.max(char.hp / char.max_hp);        
         this.ctx.drawImage(this.images.hero, graphs.x, graphs.y, graphs.w, graphs.h);
         this.ctx.fillStyle = graphs.barHealBackColor;
         this.ctx.fillRect(graphs.x, graphs.y + graphs.w + graphs.barGap, graphs.w, graphs.barH);        
         this.ctx.fillStyle = graphs.barHealColor;
         this.ctx.fillRect(graphs.x, graphs.y + graphs.w + graphs.barGap, healthPercent * graphs.w, graphs.barH);
    },
    drawBar(){
        
    },
    drawCityNodes(){
        if(!gameState.world) return;
        for (const id in gameState.world) {
            const node = gameState.world[id];
            const isCurrent = gameState.player.location_id === id;
            
            const x = node.x;
            const y = node.y;

            // 1. Рисуем кружок (основание)
            this.ctx.beginPath();
            this.ctx.arc(x, y, 6, 0, Math.PI * 2);
            this.ctx.fillStyle = isCurrent ? "#DC143C" : "#708090"; // Зеленый если мы тут, иначе золото
            this.ctx.fill();
            this.ctx.closePath();

            // 2. Рисуем палочку (шест указателя)
            this.ctx.beginPath();
            this.ctx.moveTo(x, y);
            this.ctx.lineTo(x, y - 30);
            this.ctx.strokeStyle = "#8a6d3b"; // Цвет дерева/меди
            this.ctx.lineWidth = 2;
            this.ctx.stroke();

            // 3. Рисуем прямоугольник (табличка)
            const textWidth = this.ctx.measureText(node.name).width;
            const rectW = textWidth + 20;
            const rectH = 25;
            const rectX = x - rectW / 2;
            const rectY = y - 55;

            this.ctx.fillStyle = "rgba(22, 22, 22, 0.9)"; // Темный фон таблички
            this.ctx.strokeStyle = "#d4af37";
            this.ctx.lineWidth = 1;
            this.ctx.fillRect(rectX, rectY, rectW, rectH);
            this.ctx.strokeRect(rectX, rectY, rectW, rectH);

            // 4. Текст (название локации)
            this.ctx.fillStyle = "#ffffff";
            this.ctx.font = "12px Arial";
            this.ctx.textAlign = "center";
            this.ctx.fillText(node.name, x, rectY + 17);            

        }
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