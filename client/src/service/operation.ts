import { ICompletePosition } from "../types/positions"
import api from "./api"

export const operationComplite = async (data: ICompletePosition) => {
    try {
        await api.patch(`/operations/${data.id}`, data)
    } catch (error) {
        console.log(error)
    }
}
