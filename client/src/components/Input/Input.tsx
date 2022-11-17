import React, { FC } from "react"
import classes from "./input.module.scss"

type Props = {}

export const Input: FC<Props & React.InputHTMLAttributes<HTMLInputElement>> = ({ ...attr }) => {
    return <input />
}
