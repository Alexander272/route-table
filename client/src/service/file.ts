import { AxiosError } from "axios"
import api from "./api"

export async function getFile(url: string) {
    try {
        const res = await api.get(url, {
            responseType: "blob",
        })
        return res
    } catch (error: any) {
        throw (error as AxiosError).response
    }
}
