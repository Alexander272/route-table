import { useContext, useEffect } from "react"
import { useLocation, useNavigate } from "react-router-dom"
import { SignInForm } from "./components/AuthForms/SignInForm"
import { AuthContext } from "../../context/AuthProvider"
import classes from "./auth.module.scss"
import { Header } from "../../components/Header/Header"

export default function Auth() {
    const navigate = useNavigate()
    const location = useLocation()

    const { isAuth } = useContext(AuthContext)

    const from: string = (location.state as any)?.from?.pathname || "/"

    useEffect(() => {
        if (isAuth) {
            if (from !== "/") navigate(-1)
            else navigate(from, { replace: true })
        }
    }, [isAuth, navigate, from])

    return (
        <div className={classes.page}>
            <Header />
            <div className={`${classes.container}`}>
                <SignInForm />
            </div>
        </div>
    )
}
