import React, { FC, useContext } from "react"
import { Link } from "react-router-dom"
import { AuthContext } from "../../context/AuthProvider"
import { signOut } from "../../service/auth"
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

    return (
        <div className={classes.header}>
            <div className={classes.content}>
                <a
                    className={classes.logoLink}
                    href='https://route.sealur.ru/'
                    target='_blank'
                    rel='noreferrer'
                >
                    <img
                        className={classes.logo}
                        width={340}
                        height={100}
                        loading='lazy'
                        src='/image/logo.webp'
                        alt='logo'
                    />
                </a>

                <div className={classes.nav}>
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
            </div>
        </div>
    )
}
