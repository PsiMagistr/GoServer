export function getConfigData(params){
            const body = {}
            for(let param of params){
                const elem = document.querySelector(`#${param}`)
                if(elem){
                    body[param] = elem.value
                }
                else {
                    throw new Error("Неверное значение идентификатора поля формы.")
                }
            }
            return body
        }