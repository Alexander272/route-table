import React, { FC, useContext } from "react"
import { Link } from "react-router-dom"
import { AuthContext } from "../../context/AuthProvider"
import { signOut } from "../../service/auth"
import { getReason } from "../../service/reason"
import classes from "./header.module.scss"

type Props = {}

export const Header: FC<Props> = () => {
    const { user, setUser } = useContext(AuthContext)

    const logoutHandler = async () => {
        try {
            await signOut()
            setUser(null)
        } catch (error) {}
    }

    const saveHandler = async () => {
        try {
            const res = await getReason()
            const blob = new Blob([res.data])

            const href = URL.createObjectURL(blob)
            const link = document.createElement("a")
            link.href = href
            link.download = res.headers["content-disposition"]?.split("=")[1] || ""
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
        } catch (error) {
            console.log(error)
        }
    }

    return (
        <div className={classes.header}>
            <div className={classes.content}>
                <Link className={classes.logoLink} to='/'>
                    <img
                        className={classes.logo}
                        width={192}
                        height={192}
                        loading='lazy'
                        src='/logo192.webp'
                        alt='logo'
                    />
                    <span>Sealur Route</span>
                </Link>

                {user && (
                    <div className={classes.nav}>
                        {user?.role === "master" ? (
                            <p className={classes.profile} onClick={saveHandler}>
                                <img
                                    src='/image/download.svg'
                                    alt='download'
                                    width='28'
                                    height='28'
                                />
                            </p>
                        ) : null}

                        <Link to='/' className={classes.profile}>
                            <img src='/image/home.svg' alt='home' width='32' height='32' />
                        </Link>

                        {user?.role === "master" || user?.role === "display" ? (
                            <Link to='/orders/group' className={classes.profile}>
                                <img src='/image/list.svg' alt='orders' width='30' height='30' />
                            </Link>
                        ) : null}

                        <div className={classes.profile} onClick={logoutHandler}>
                            <img src='/image/logout.svg' alt='log-out' width='30' height='30' />
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}
