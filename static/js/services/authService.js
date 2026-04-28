import { apiCall } from "../api.js";
export const authService = {
    async login(email, password){
        const response = await fetch("./api/login", {
            method:"POST",
            headers: {
                "Content-Type":"application/json",
            },
            body:JSON.stringify({email, password})    
        });
        if(response.ok){
            const data = await response.json();
            return data;
        }
         throw new Error(await response.text());
    },
    async register(body){
        const response = await apiCall("/api/register", "POST", body);
        if (!response.ok) throw new Error(await response.text());
        return await response.json();
    }
}