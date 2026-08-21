export const moveModalTemplate = (data)=>{
    return `
    <div id="move-content">
        <p>Идем в <span class="target-name">${data.world_name}, ${data.location_name}</span>...</p>
        <span class="timer-count" id="move-timer">${data.duration}</span>
    </div>`
}