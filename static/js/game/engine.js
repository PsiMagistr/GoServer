import { gameState } from "./game.js";
import { graphs } from "./graphic_settings.js";
export const engine = {
    main: {
        canvas: null,
        ctx: null,
        loopId: null,
    },
    battle: {
        canvas: null,
        ctx: null,
        loopId: null,
        data: null,
        loopId: null,
    },
    images: null,
    init(canvasElement) {
        this.main.canvas = canvasElement;
        this.main.ctx = canvasElement.getContext('2d');
        this.startMainLoop();
    },
    initBattle(canvasElement, battleData) {        
        this.battle.canvas = canvasElement;
        this.battle.ctx = canvasElement.getContext('2d');
        this.battle.data = battleData;
        this.startBattleLoop();
    },

    startMainLoop() {
        if (this.main.loopId) cancelAnimationFrame(this.main.loopId);
        const render = () => {
            this.drawWorld()
            this.main.loopId = window.requestAnimationFrame(render);
            console.log("Главный цикл");
        }
        render();
    },
    startBattleLoop() {
        if (this.battle.loopId) cancelAnimationFrame(this.battle.loopId);
        const render = () => {
            this.renderBattle();
            this.battle.loopId = requestAnimationFrame(render);            
        }
        render();
    },
    drawWorld() {
        const canvas = this.main.canvas;
        const ctx = this.main.ctx;
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        if (!gameState.player) {
            ctx.fillStyle = 'lime';
            ctx.fillRect(280, 180, 40, 40);
            return;
        }
        this.drawMap();
        this.drawCityNodes();
    },
    stopMainLoop() {
        cancelAnimationFrame(this.main.loopId);
    },
    stopBattleLoop() {
        cancelAnimationFrame(this.battle.loopId);
    },
    renderBattle() {
        const ctx = this.battle.ctx;
        const canvas = this.battle.canvas;
        if (!ctx || !this.battle.data) return;
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        this.drawUnit(ctx, 10, 10, this.battle.data.you, true);
        if (this.battle.data.opponent) {
            this.drawUnit(ctx, 400, 10, this.battle.data.opponent, false);
        }
    },

    drawMap() {
        const canvas = this.main.canvas;
        const ctx = this.main.ctx;
        ctx.drawImage(this.images.map, 0, 0, canvas.width, canvas.height);
        this.drawUnit(ctx, 10, 10, gameState.player, true);
    },
    drawUnit(ctx, x, y, char, isPlayer) {
        const avatar = isPlayer ? this.images.hero : this.images.opponent;
        if (avatar) {
            ctx.drawImage(avatar, x, y, 100, 100);
        }
        const healParams = {
            x: x,
            y: y,
            w: graphs.w,
            h: graphs.barH,
            current: char.hp,
            currentRounded: Math.ceil(char.hp),
            max: char.max_hp,
            backColor: graphs.barHealBackColor,
            color: graphs.barHealColor,
        }
        const manaParams = {
            x: x,
            y: y,
            w: graphs.w,
            h: graphs.barH,
            current: char.mana,
            currentRounded: Math.floor(char.mana),
            max: char.max_mana,
            backColor: graphs.barHealBackColor,
            color: graphs.barManaColor,
        }
        const expaParams = {
            x: x,
            y: y,
            w: graphs.w,
            h: graphs.barH,
            current: char.exp,
            currentRounded: Math.floor(char.exp),
            max: char.max_exp,
            backColor: graphs.barHealBackColor,
            color: graphs.barExpColor,
        }

        this.drawBar(1, ctx, healParams);
        this.drawBar(2, ctx, manaParams);
        this.drawBar(3, ctx, expaParams);
    },
    drawBar(counter, ctx, { x, y, w, h, current, max, color, backColor, currentRounded }) {
        const dy = y + w + graphs.barGap + (graphs.barH + graphs.barGap) * (counter - 1);
        const percent = Math.min(Math.max(current / max, 0), 1)
        ctx.fillStyle = backColor;
        ctx.fillRect(x, dy, w, h);
        ctx.fillStyle = color;
        ctx.fillRect(x, dy, w * percent, h);
        ctx.fillStyle = "white";
        ctx.font = `bold ${h * 0.8}px Arial`;
        ctx.textAlign = "center";
        ctx.fillText(`${currentRounded}/${max}`, x + w / 2, dy + h * 0.85);
    },
    drawCityNodes() {
        const points = gameState.world?.points;
        if (!points) return;
        for (const id in points) {
            const node = points[id];
            this.drawPointer(node, id);
        }
    },
    drawPointer(node, id) {
        const ctx = this.main.ctx;
        const isCurrent = gameState.player.location_id === id;
        const isHovered = gameState.hoveredNodeId === id;
        let color = "#FFDEAD";
        let backgroundTableColor = "rgba(22, 22, 22, 0.9)";
        if (isCurrent) {
            color = "#DC143C";
        }
        else if (isHovered) {
            color = "#90EE90";
            backgroundTableColor = "#696969";
        }
        const x = node.x;
        const y = node.y;
        // 1. Рисуем кружок (основание)
        ctx.beginPath();
        ctx.arc(x, y, 8, 0, Math.PI * 2);
        ctx.fillStyle = color;
        ctx.fill();
        ctx.closePath();
        // 2. Рисуем палочку (шест указателя)
        ctx.beginPath();
        ctx.moveTo(x, y);
        ctx.lineTo(x, y - 30);
        ctx.strokeStyle = "#8a6d3b"; // Цвет дерева/меди
        ctx.stroke()
        // 3. Рисуем прямоугольник (табличка)
        const textWidth = ctx.measureText(node.name).width;
        const rectW = textWidth + 20;
        const rectH = 25;
        const rectX = x - rectW / 2;
        const rectY = y - 55;
        ctx.fillStyle = backgroundTableColor; // Темный фон таблички
        ctx.strokeStyle = "#d4af37";
        ctx.lineWidth = 1;
        ctx.fillRect(rectX, rectY, rectW, rectH);
        ctx.strokeRect(rectX, rectY, rectW, rectH);
        // 4. Текст (название локации)
        ctx.fillStyle = "#ffffff";
        ctx.font = "12px Arial";
        ctx.textAlign = "center";
        ctx.fillText(node.name, x, rectY + 17);
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