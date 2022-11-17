import React, { FC } from "react"
import { UseFormRegisterReturn } from "react-hook-form"
import classes from "./input.module.scss"

type Props = {
    id?: string
    register?: UseFormRegisterReturn
}

export const InputBase: FC<Props & React.InputHTMLAttributes<HTMLInputElement>> = ({
    register,
    ...attr
}) => {
    return <input className={classes.base} {...register} {...attr} />
}
