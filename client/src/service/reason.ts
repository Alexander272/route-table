import { AxiosError } from "axios"
import api from "./api"

export async function getReason() {
    try {
        const res = await api.get("/reasons/file", {
            responseType: "blob",
        })
        return res
    } catch (error: any) {
        throw (error as AxiosError).response
    }
}
