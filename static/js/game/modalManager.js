export const modalManager = {
    currentTemplateFunc: null,
    show(templateFunc, data, type){         
        const overlay = document.querySelector(type);
        const modalType = type == "#modal-overlay" ? "#modal-content":"#modal-battle-content";       
        const content = document.querySelector(modalType);        
        if(!overlay || !content) return
         // Запоминаем, какой шаблон мы используем сейчас
        this.currentTemplateFunc = templateFunc;
        content.innerHTML = this.currentTemplateFunc(data);        
        // ВАЖНО: используем flex для центрирования из CSS
        overlay.style.display = 'flex';              

    },
    hide(){
        let overlay = document.querySelector("#modal-overlay");
        if (overlay) {            
            overlay.style.display = 'none';
            this.currentTemplateFunc = null; // Очищаем память
        }
       overlay = document.querySelector("#batlle-overlay");
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