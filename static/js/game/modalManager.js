export const modalManager = {
    currentTemplateFunc: null,
    show(templateFunc, data){
        const overlay = document.querySelector("#modal-overlay");
        const content = document.querySelector("#modal-content");
        if(!overlay || !content) return
         // Запоминаем, какой шаблон мы используем сейчас
        this.currentTemplateFunc = templateFunc;
        content.innerHTML = this.currentTemplateFunc(data);        
        // ВАЖНО: используем flex для центрирования из CSS
        overlay.style.display = 'flex';

    },
    hide(){
        const overlay = document.querySelector("#modal-overlay");
        if (overlay) {
            overlay.style.display = 'none';
            this.currentTemplateFunc = null; // Очищаем память
        }
    },
    refresh(data){
        const overlay = document.querySelector("#modal-overlay");
        const content = document.querySelector("#modal-content");
        if (overlay && overlay.style.display === 'flex' && this.currentTemplateFunc) {
            content.innerHTML = this.currentTemplateFunc(data);
        }    
    },        
}