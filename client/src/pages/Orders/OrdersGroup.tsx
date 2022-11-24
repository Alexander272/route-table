import React, { useEffect, useState } from "react"
import useSWR from "swr"
import Masonry from "@mui/lab/Masonry"
import { fetcher } from "../../service/read"
import { IUrgencyGroup } from "../../types/order"
import { OrderItem } from "./OrderItem"
import classes from "./orders.module.scss"

const itemSize = 307

export default function OrdersGroup() {
    const { data: res, error } = useSWR<{ data: IUrgencyGroup }>("/orders/group", fetcher, {
        refreshInterval: 60 * 1000,
    })
    const [columns, setColumns] = useState([0, 0, 0])

    useEffect(() => {
        if (!res) return

        let count = 0
        let newColumns = [0, 0, 0]

        if (res.data.high) {
            count++
            newColumns[0] = -1
        }
        if (res.data.middle) {
            count++
            newColumns[1] = -1
        }
        if (res.data.low) {
            count++
            newColumns[2] = -1
        }

        newColumns = newColumns.map(c => {
            if (c === -1) return Math.trunc((window.innerWidth - 40) / itemSize / count)
            else return c
        })
        console.log(newColumns)

        setColumns(newColumns)
    }, [res, setColumns])

    if (error) return null

    return (
        <div className={classes.container}>
            {res?.data.high && (
                <div className={classes.column}>
                    <Masonry columns={columns[0]} spacing={2}>
                        {res?.data.high.map(o => (
                            <OrderItem key={o.id} order={o} />
                        ))}
                    </Masonry>
                </div>
            )}
            {res?.data.middle && (
                <div className={`${classes.column} ${classes["column-middle"]}`}>
                    <Masonry columns={2} spacing={2}>
                        {res?.data.middle.map(o => (
                            <OrderItem key={o.id} order={o} />
                        ))}
                    </Masonry>
                </div>
            )}
            {res?.data.low && (
                <div className={classes.column}>
                    <Masonry columns={columns[2]} spacing={2}>
                        {res?.data.low.map(o => (
                            <OrderItem key={o.id} order={o} />
                        ))}
                    </Masonry>
                </div>
            )}
        </div>
    )
}
