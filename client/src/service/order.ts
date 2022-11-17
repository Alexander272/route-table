import axios from "axios"

export const orderParse = async (data: FormData) => {
    try {
        await axios.post("/api/v1/orders/parse", data)
    } catch (error: any) {
        console.log(error)
        throw error.response.message
    }
}
