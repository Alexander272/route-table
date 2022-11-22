import { useContext } from "react"
import { Navigate, useLocation } from "react-router-dom"
import { AuthContext } from "../../context/AuthProvider"

export default function RequireAuth({ children }: { children: JSX.Element }) {
    const { isAuth, user } = useContext(AuthContext)
    const location = useLocation()

    if (!isAuth) return <Navigate to='/auth' state={{ from: location }} />

    return children
}
