import React from "react"
import { useRefresh } from "./hooks/refresh"
import { MyRoutes } from "./routes"

export default function Main() {
    const { ready } = useRefresh()

    if (!ready) return <></>
    return <MyRoutes />
}
