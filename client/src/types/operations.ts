export interface IOperation {
    id: string
    title: string
    done: boolean
    remainder: number
}

export interface ICompleteOperation {
    id: string
    done: boolean
    remainder: number
    reason: string
}
