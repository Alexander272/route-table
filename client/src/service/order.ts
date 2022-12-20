import api from "./api"
import { IUpdateOrder } from "../types/order"

export const orderParse = async (data: FormData) => {
    try {
        await api.post("/orders/parse", data)
    } catch (error: any) {
        console.log(error)
        throw error.response.message
    }
}

export const orderUpdate = async (data: IUpdateOrder) => {
    try {
        await api.put(`/orders/${data.id}`, data)
    } catch (error: any) {
        console.log(error)
        throw error.response.message
    }
}

export const orderDelete = async (id: string) => {
    try {
        await api.delete(`/orders/${id}`)
    } catch (error: any) {
        throw error.response.message
    }
}
