import React, { createContext, FC, useState } from "react"
import { IUser } from "../types/user"

type Context = {
    isAuth: boolean
    // setIsAuth: (isAuth: boolean) => void
    user: IUser | null
    setUser: (user: IUser | null) => void
}

export const AuthContext = createContext<Context>({
    isAuth: false,
    // setIsAuth: (isAuth: boolean) => {},
    user: null,
    setUser: (user: IUser | null) => {},
})

export const AuthProvider: FC = ({ children }) => {
    const [isAuth, setIsAuth] = useState(false)
    const [user, setUser] = useState<IUser | null>(null)

    // const authHandler = (isAuth: boolean) => setIsAuth(isAuth)
    const userHandler = (user: IUser | null) => {
        setUser(user)
        setIsAuth(!!user)
    }

    return (
        <AuthContext.Provider value={{ isAuth, user, setUser: userHandler }}>
            {children}
        </AuthContext.Provider>
    )
}
