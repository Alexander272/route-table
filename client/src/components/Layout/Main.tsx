import React from "react"
import { Outlet } from "react-router-dom"
import { Header } from "../Header/Header"
import classes from "./main.module.scss"

export default function Main() {
    return (
        <div className={classes.page}>
            <Header />
            <Outlet />
        </div>
    )
}
