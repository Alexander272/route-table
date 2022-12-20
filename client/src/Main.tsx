import React, { useContext } from "react"
import { SWRConfig } from "swr"
import { AuthContext } from "./context/AuthProvider"
import { useRefresh } from "./hooks/refresh"
import { MyRoutes } from "./routes"
import { refresh } from "./service/auth"

export default function Main() {
    const { ready } = useRefresh()
    const { setUser } = useContext(AuthContext)

    if (!ready) return <></>
    return (
        <SWRConfig
            value={{
                onErrorRetry: async (error, key, config, revalidate, { retryCount }) => {
                    if (error.status === 401) {
                        try {
                            const res = await refresh()
                            setUser(res.data)
                        } catch (error: any) {
                            setUser(null)
                            return
                        }
                    }
                    if (error.status === 403 || error.status === 404) return
                    if (retryCount >= 5) return
                    setTimeout(() => revalidate({ retryCount }), 5000)
                },
            }}
        >
            <MyRoutes />
        </SWRConfig>
    )
}
