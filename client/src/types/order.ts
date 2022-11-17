import { IPositionForOrder } from "./positions"

export interface IFindedOrder {
    id: string
    number: string
    done: boolean
}

export interface IOrder {
    id: string
    number: string
    done: boolean
    positions: IPositionForOrder[]
}

export interface ISearchForm {
    search: string
    resultIndex: number
}
