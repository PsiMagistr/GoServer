import { gameState } from "./game.js";
import { graphs } from "./graphic_settings.js";
export const engine = {
    canvas: null,
    ctx: null,
    loopId: null,
    images: null,
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
        if (!this.images) return;
        this.drawMap();
        this.drawCityNodes();
        this.drawPlayer();

    },
    stop() {
        cancelAnimationFrame(this.loopId);
    },
    drawMap() {
        this.ctx.drawImage(this.images.map, 0, 0, this.canvas.width, this.canvas.height);
    },
    drawPlayer() {
        const char = gameState.player;                      
        this.ctx.drawImage(this.images.hero, graphs.x, graphs.y, graphs.w, graphs.h);        
        const healParams = {           
            x:graphs.x,
            y: graphs.y,
            w:graphs.w,
            h:graphs.barH,            
            current:char.hp,
            currentRounded:Math.ceil(char.hp),            
            max:char.max_hp,
            backColor:graphs.barHealBackColor,
            color:graphs.barHealColor,
        }
        const manaParams = {           
            x:graphs.x,
            y: graphs.y,
            w:graphs.w,
            h:graphs.barH,            
            current:char.mana,
            currentRounded:Math.floor(char.mana),
            max:char.max_mana,
            backColor:graphs.barHealBackColor,
            color:graphs.barManaColor,
        }
        const expaParams = {           
            x:graphs.x,
            y: graphs.y,
            w:graphs.w,
            h:graphs.barH,            
            current:char.exp,
            currentRounded:Math.floor(char.exp),
            max:char.max_exp,
            backColor:graphs.barHealBackColor,
            color:graphs.barExpColor,
        }         
                                                 
        this.drawBar(1, healParams);
        this.drawBar(2, manaParams);
        this.drawBar(3, expaParams);        
    },
    drawBar(counter, {x, y, w, h, current, max, color, backColor, currentRounded}) {        
        const dy = y + w + graphs.barGap + (graphs.barH + graphs.barGap) * (counter - 1);
        const percent = Math.min(Math.max(current / max, 0), 1)
        this.ctx.fillStyle = backColor;
        this.ctx.fillRect(x, dy, w, h);
        this.ctx.fillStyle = color;
        this.ctx.fillRect(x, dy, w * percent, h);
        this.ctx.fillStyle = "white";
        this.ctx.font = `bold ${h * 0.8}px Arial`;
        this.ctx.textAlign = "center";
        this.ctx.fillText(`${currentRounded}/${max}`, x + w / 2, dy + h * 0.85);
    },
    drawCityNodes() {
        const points = gameState.world?.points;
        if (!points) return;
        for (const id in points) {
            const node = points[id];
            this.drawPointer(node, id);
        }
    },
    drawPointer(node, id){
        const isCurrent = gameState.player.location_id === id;
            const isHovered = gameState.hoveredNodeId === id;
            let color = "#FFDEAD";
            let backgroundTableColor = "rgba(22, 22, 22, 0.9)";
            if(isCurrent){
                color = "#DC143C";
            }
            else if(isHovered){
                color = "#90EE90";
                backgroundTableColor = "#696969";
            }
            const x = node.x;
            const y = node.y;
            // 1. Рисуем кружок (основание)
            this.ctx.beginPath();
            this.ctx.arc(x, y, 8, 0, Math.PI * 2);
            this.ctx.fillStyle = color;
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
            this.ctx.fillStyle = backgroundTableColor; // Темный фон таблички
            this.ctx.strokeStyle = "#d4af37";
            this.ctx.lineWidth = 1;
            this.ctx.fillRect(rectX, rectY, rectW, rectH);
            this.ctx.strokeRect(rectX, rectY, rectW, rectH);
            // 4. Текст (название локации)
            this.ctx.fillStyle = "#ffffff";
            this.ctx.font = "12px Arial";
            this.ctx.textAlign = "center";
            this.ctx.fillText(node.name, x, rectY + 17);
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