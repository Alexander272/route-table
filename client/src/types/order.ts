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
    deadline: string
    positions: IPositionForOrder[]
}

export interface IUpdateOrder {
    id: string
    deadline: string
}

export interface ISearchForm {
    search: string
    resultIndex: number
}

export interface IGroupedOrder {
    id: string
    deadline: string
    urgency: string
    orders: IOrderItem[]
}

export interface IUrgencyGroup {
    high: IGroupedOrder[]
    middle: IGroupedOrder[]
    low: IGroupedOrder[]
}

export interface IOrderItem {
    id: string
    number: string
    done: boolean
    date: string
    progress: number
}
