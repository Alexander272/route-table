import axios from "axios"
import { ISignIn, IUser } from "../types/user"

export const signIn = async (data: ISignIn): Promise<{ data: IUser }> => {
    try {
        const res = await axios.post("/api/v1/auth/sign-in", data)
        return res.data
    } catch (error: any) {
        throw error.response.data
    }
}

export const signOut = async () => {
    try {
        const res = await axios.post("/api/v1/auth/sign-out")
        return res.data
    } catch (error: any) {
        throw error.response.data
    }
}

export const refresh = async (): Promise<{ data: IUser }> => {
    try {
        const res = await axios.post("/api/v1/auth/refresh")
        return res.data
    } catch (error: any) {
        throw error.response.data
    }
}
