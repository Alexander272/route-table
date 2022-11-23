import api from "./api"

export const orderParse = async (data: FormData) => {
    try {
        await api.post("/orders/parse", data)
    } catch (error: any) {
        console.log(error)
        throw error.response.message
    }
}
