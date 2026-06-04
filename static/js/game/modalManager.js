export const modalManager = {
    currentTemplate: null,
    show(templateFunc, data){
        const overlay = document.querySelector("#modal-overlay");
        const content = document.querySelector("#modal-content");
        if(!overlay || !content) return
         // Запоминаем, какой шаблон мы используем сейчас
        this.currentTemplate = templateFunc;
        content.innerHTML = templateFunc(data);        
        // ВАЖНО: используем flex для центрирования из CSS
        overlay.style.display = 'flex';

    },
    hide(){
        const overlay = document.querySelector("#modal-overlay");
        if (overlay) {
            overlay.style.display = 'none';
            this.currentTemplate = null; // Очищаем память
        }
    }
}