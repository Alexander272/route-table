import React, { FC } from "react"
import { IFindedOrder } from "../../../../types/order"

import classes from "./find.module.scss"

type Props = {
    order: IFindedOrder
    index: number
    selected: boolean
    selectHandler: (order: IFindedOrder) => void
}

export const ResultItem: FC<Props> = ({ order, index, selectHandler, selected }) => {
    return (
        <li
            role='option'
            aria-selected='false'
            tabIndex={index}
            className={[classes.item, selected ? classes.active : null].join(" ")}
            onClick={() => selectHandler(order)}
        >
            <span>{order.number}</span>
            <span></span>
        </li>
    )
}
