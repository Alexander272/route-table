import { createContext } from "react"

export const OrderContext = createContext({
    orderId: "",
    changeOrderId: (orderId: string) => {},
})
