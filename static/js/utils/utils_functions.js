export const utils = {
    getElementsBySelectors(objSelectors, context= document) {
        const result = {};
        let currentContext = context;
        if(typeof objSelectors.elem === "string"){
            result.elem = currentContext.querySelector(objSelectors.elem);
            if(result.elem instanceof Element){
                    currentContext = result.elem
            }
            else{
                throw new Error(`CRITICAL UI ERROR: Не удалось найти контейнер по селектору ${objSelectors.elem} внутри контекста ${context}`);
            }
        }
        for(let key in objSelectors){
            if(key === "elem" && typeof objSelectors.elem === "string"){
                continue;
            }
            const value = objSelectors[key];
            if(typeof value === 'string'){
                result[key] = currentContext.querySelector(value);
            }
            else if(typeof value === 'object' && value !== null){
                result[key] = this.getElementsBySelectors(value, currentContext);
            }
            else{
                result[key] = value;
            }
        }
        return result;
    },
    getValuesBySchema(schema){
        const result = {};
        for(let key in schema){
            const item = schema[key]
            if(item === null){
                result[key] = null;
                continue;
            }
            if(item instanceof HTMLElement){
                result[key] = ("value" in item) ? item.value : null;           
            }
            else if(typeof item == "object"){
                result[key] = item;
            }            
        }
        return result;
    },
    createTimer(seconds, onTick, onEnd){
        let timerId = null;
        const endTime = Date.now() + (seconds * 1000);
        let sec = seconds;
        const time = ()=>{
            const now = Date.now();
            const diff = endTime - now;
            const sec = Math.ceil(diff / 1000);
            //sec--;
            if(onTick) onTick(sec)            
            if(sec == 0){                
                if(onEnd) onEnd()                  
                return     
            }
            timerId = setTimeout(time, 1000);            
        }
        timerId = setTimeout(time, 1000);
        return ()=>{
            clearTimeout(timerId)
        }
    },
}
