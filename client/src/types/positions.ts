import { ICompleteOperation, IOperation } from "./operations"

export interface IPositionForOrder {
    id: string
    position: number
    count: number
    title: string
    ring: string
    deadline: string
    connected: string
    done: boolean
    lastOperation: string | null
    curOperation: string | null
}

export interface IPosition {
    id: string
    order: string
    position: number
    count: number
    title: string
    ring: string
    deadline: string
    connected: string
    done: boolean
    operations: IOperation[]
}

export interface ICompletePosition {
    id: string
    count: number
    isFinish: boolean
    connected: string
    operation: ICompleteOperation
}
