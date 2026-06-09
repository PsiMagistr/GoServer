import { utils } from "../utils/utils_functions.js";

export const ui = {
    renderList(containerSelector, items, idPrefix, className, labelGenerator){
        const container = document.querySelector(containerSelector);
        if (!container) return;
        container.innerHTML = ""; // Очищаем список

        if (!items || items.length === 0) {
            container.innerHTML = '<li class="empty">Пусто</li>';
            return;
        }
        for(const item of items){
            this.addItemToList(containerSelector, item, idPrefix, className, labelGenerator);
        }
    },
    addItemToList(containerSelector, item, idPrefix, className, labelGenerator){
        const container = document.querySelector(containerSelector);
        if(!container) return
        const itemId = `${idPrefix}-${item.id}`;
        const empty = container.querySelector('.empty');
        if (empty) empty.remove();
        const li = document.createElement("li");
        li.id = itemId;
        li.dataset.id = item.id;
        li.className = className;
        li.innerHTML = labelGenerator(item);
        container.appendChild(li);
    },
    removeItemFromUI(idPrefix, id){
        const el = document.querySelector(`#${idPrefix}-${id}`);
        if(el) el.remove();    
    }
}