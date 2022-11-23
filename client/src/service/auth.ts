import { ISignIn, IUser } from "../types/user"
import api from "./api"

export const signIn = async (data: ISignIn): Promise<{ data: IUser }> => {
    try {
        const res = await api.post("/auth/sign-in", data)
        return res.data
    } catch (error: any) {
        throw error.response.data
    }
}

export const signOut = async () => {
    try {
        const res = await api.post("/auth/sign-out")
        return res.data
    } catch (error: any) {
        throw error.response.data
    }
}

export const refresh = async (): Promise<{ data: IUser }> => {
    try {
        const res = await api.post("/auth/refresh")
        return res.data
    } catch (error: any) {
        throw error.response.data
    }
}
