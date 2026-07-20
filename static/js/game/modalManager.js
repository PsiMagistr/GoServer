export const modalManager = {
    currentTemplateFunc: null,
    activeController: null,
    show(templateFunc, data, options = {closable: true}, controller=null) {
        this.activeController = controller;
        const overlay = document.querySelector("#modal-overlay");
        const content = document.querySelector("#modal-content");
        if (!overlay || !content) return
        const closeBtn = document.querySelector("#modal-close-btn");
        // Запоминаем, какой шаблон мы используем сейчас
        this.currentTemplateFunc = templateFunc;
        content.innerHTML = this.currentTemplateFunc(data);
        // ВАЖНО: используем flex для центрирования из CSS
        overlay.style.display = 'flex';
        if (options.closable) {
            closeBtn.style.display = 'block';
        }
        else {
            closeBtn.style.display = "none"
        }

    },    
    hide() {       
        if (this.activeController && this.activeController.onHide) {
            this.activeController.onHide();
        }
        let overlay = document.querySelector("#modal-overlay");
        if (overlay) {
            overlay.style.display = 'none';
            this.currentTemplateFunc = null; // Очищаем память
        }
    },
    refresh(data) {
        const overlay = document.querySelector("#modal-overlay");
        const content = document.querySelector("#modal-content");
        if (overlay && overlay.style.display === 'flex' && this.currentTemplateFunc) {
            content.innerHTML = this.currentTemplateFunc(data);
        }
    },    
    setClosable(status) {
        const closeBtn = document.querySelector("#modal-close-btn");
        const overlay = document.querySelector("#modal-overlay");
        if (!closeBtn || !overlay) return;

        if (status) {
            closeBtn.style.display = 'block';
            // Если окно можно закрыть, разрешаем закрытие кликом по фону (оверлею)
            overlay.style.cursor = 'pointer';
        } else {
            closeBtn.style.display = "none";
            // Если окно закрывать нельзя, убираем курсор и возможность клика по фону
            overlay.style.cursor = 'default';
        }
    },
}