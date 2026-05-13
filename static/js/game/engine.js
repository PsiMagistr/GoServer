export const engine = {
    canvas:null,
    ctx:null,
    loopId:null,
    init(canvasElement){
        this.canvas = canvasElement;
        this.ctx = canvasElement.getContext('2d');
        this.start();
    },
    start(){
        if(this.loopId) this.stop();
        const render = ()=>{
            this.draw()
            this.loopId = window.requestAnimationFrame(render)
        }
        render();
    },
    draw(){
        this.ctx.fillStyle = 'black';
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);       
        this.ctx.fillStyle = 'lime';
        this.ctx.fillRect(280, 180, 40, 40);
    },
    stop(){
         cancelAnimationFrame(this.loopId);
    },
}