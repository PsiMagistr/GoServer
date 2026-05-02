import { apiCall } from "../api.js";
export const characterService = {
    async create(dataChar){
        const response = await apiCall("/api/character/create", "POST", dataChar)
            if(response.ok){
                const data = await response.json()
                console.log("Персонаж сохранен!")
                console.log(data)
                return;
            }
            throw new Error(await response.text())                
    }
}