import { IUrgency } from "../types/urgency"
import api from "./api"

export async function changeUrgency(data: IUrgency) {
    try {
        await api.put(`/urgency`, data)
    } catch (error: any) {
        throw error.response.message
    }
}
