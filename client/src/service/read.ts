import api from "./api"

export const fetcher = async (url: string) => {
    try {
        const res = await api.get(url)
        return res.data
    } catch (error: any) {
        throw error.response
    }
}
