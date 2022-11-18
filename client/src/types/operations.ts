import { IReason } from "./reason"

export interface IOperation {
    id: string
    title: string
    done: boolean
    remainder: number
    reasons?: IReason[]
}

export interface ICompleteOperation {
    id: string
    done: boolean
    remainder: number
    reason: string
}
