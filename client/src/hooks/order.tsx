import { useCallback, useState } from "react"

export const useOrder = () => {
    const [orderId, setOrderId] = useState("")

    const changeOrderId = useCallback((orderId: string) => {
        setOrderId(orderId)
    }, [])

    return { orderId, changeOrderId }
}
