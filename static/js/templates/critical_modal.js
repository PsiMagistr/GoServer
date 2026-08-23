export const modalCriticalTemplate = (data)=>{
    console.log("Шаблон кирпича.")    
    return `    
    <div class="overlay-content">
        <h2>Сессия прервана</h2>
        <p>Вы вошли в игру с другого устройства или в новой вкладке.</p>
        <p>Это окно больше не получает обновлений.</p>
        <button onclick="window.location.reload()" class="main-btn">ВЕРНУТЬ УПРАВЛЕНИЕ СЮДА</button>
    </div>`    
}