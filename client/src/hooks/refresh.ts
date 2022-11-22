import { useCallback, useContext, useEffect, useState } from "react"
import { AuthContext } from "../context/AuthProvider"
import { refresh } from "../service/auth"

export function useRefresh() {
    const [ready, setReady] = useState(false)
    const { user, setUser } = useContext(AuthContext)

    const refreshUser = useCallback(async () => {
        if (ready || user) return
        try {
            const res = await refresh()
            setUser(res.data)
        } catch (error: any) {
            console.log(error.response.message)
        } finally {
            setReady(true)
        }
    }, [setUser, ready, user])

    useEffect(() => {
        refreshUser()
    }, [refreshUser])

    return { ready }
}
