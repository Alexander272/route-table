import React, { useCallback, useContext, useEffect, useRef, useState } from "react"
import useSWR from "swr"
import Masonry from "@mui/lab/Masonry"
import { Navigate } from "react-router-dom"
import { AuthContext } from "../../../context/AuthProvider"
import { fetcher } from "../../../service/read"
import { IUrgencyGroup } from "../../../types/order"
import { OrderItem } from "./OrderItem"
import classes from "./orders.module.scss"
import { Typography } from "@mui/material"

const itemSize = 820

export default function OrdersGroup() {
    const { user } = useContext(AuthContext)
    const container = useRef<HTMLDivElement>(null)
    const scroll = useRef<number>(0)

    const { data: res, error } = useSWR<{ data: IUrgencyGroup }>("/orders/group", fetcher, {
        refreshInterval: 60 * 1000,
    })
    const [columns, setColumns] = useState([0, 0])

    useEffect(() => {
        if (!res) return

        let count = 0
        let newColumns = [0, 0]

        if (res.data.high) {
            count++
            newColumns[0] = -1
        }
        if (res.data.middle) {
            count++
            newColumns[1] = -1
        }

        newColumns = newColumns.map(c => {
            if (c === -1) return Math.trunc((window.innerWidth - 40) / itemSize / count)
            else return c
        })
        console.log(newColumns)

        setColumns(newColumns)
    }, [res, setColumns])

    const scrollToMyRef = useCallback(() => {
        if (scroll.current === 0) {
            scroll.current =
                (container.current?.scrollHeight || 0) - (container.current?.clientHeight || 0)
        } else {
            scroll.current = 0
        }

        container.current?.scrollTo({ left: 0, top: scroll.current, behavior: "smooth" })
    }, [])

    useEffect(() => {
        const timer = setInterval(() => {
            scrollToMyRef()
        }, 20 * 1000)
        return () => clearInterval(timer)
    }, [scrollToMyRef])

    useEffect(() => {
        if (!user) return
        if (user?.role !== "display" && user?.role !== "master") {
            Navigate({ to: "/", replace: true })
        }
    }, [user])

    if (error) return null

    return (
        <>
            <Typography
                variant='h5'
                sx={{ textAlign: "center", marginTop: 1, fontWeight: "bold", fontSize: "4rem" }}
                className={classes.time}
            >
                {new Date().toLocaleString("ru-RU", { dateStyle: "short", timeStyle: "short" })}
            </Typography>

            <div className={classes.container} ref={container}>
                {res?.data.high && (
                    <div className={classes.column}>
                        <Masonry columns={columns[0]} spacing={3}>
                            {res?.data.high.map(o => (
                                <OrderItem key={o.id} order={o} />
                            ))}
                        </Masonry>
                    </div>
                )}
                {res?.data.middle && (
                    <div className={classes.column}>
                        <Masonry columns={columns[1]} spacing={3}>
                            {res?.data.middle.map(o => (
                                <OrderItem key={o.id} order={o} />
                            ))}
                        </Masonry>
                    </div>
                )}
            </div>
        </>
    )
}
