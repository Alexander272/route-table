import axios from "axios"
import { ICompletePosition } from "../types/positions"

export const operationComplite = async (data: ICompletePosition) => {
    try {
        await axios.patch(`/api/v1/operations/${data.id}`, data)
    } catch (error) {
        console.log(error)
    }
}
