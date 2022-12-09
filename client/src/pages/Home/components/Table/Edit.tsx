import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    IconButton,
    TextField,
} from "@mui/material"
import EditIcon from "@mui/icons-material/Edit"
import React, { FC, useEffect, useState } from "react"
import { IOrder } from "../../../../types/order"
import { orderUpdate } from "../../../../service/order"
import { useSWRConfig } from "swr"

type Props = {
    order: IOrder
}

export const Edit: FC<Props> = ({ order }) => {
    const [open, setOpen] = useState(false)
    const [deadline, setDeadline] = useState(Date.now())
    const { mutate } = useSWRConfig()

    useEffect(() => {
        if (order) setDeadline(+order.deadline * 1000)
    }, [order])

    const handleDeadline = (event: React.ChangeEvent<HTMLInputElement>) => {
        setDeadline(event.target.valueAsNumber)
    }

    const handleClickOpen = () => {
        setOpen(true)
    }

    const handleClose = () => {
        setOpen(false)
    }

    const handleSave = async () => {
        try {
            const data = { id: order.id, deadline: (deadline / 1000).toString() }
            await orderUpdate(data)
            mutate(`/orders/${order.id}`)
            handleClose()
        } catch (error) {
            console.log(error)
        }
    }

    return (
        <>
            <Dialog open={open} onClose={handleClose}>
                <DialogTitle>Редактировать заказ</DialogTitle>
                <DialogContent>
                    {/* <DialogContentText>
                        To subscribe to this website, please enter your email address here. We will
                        send updates occasionally.
                    </DialogContentText> */}
                    <TextField
                        id='date'
                        label='Дата отгрузки'
                        type='date'
                        value={new Date(deadline).toISOString().split("T")[0]}
                        onChange={handleDeadline}
                        sx={{ width: 320, marginTop: "20px" }}
                        InputLabelProps={{
                            shrink: true,
                        }}
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>Отмена</Button>
                    <Button onClick={handleSave}>Сохранить</Button>
                </DialogActions>
            </Dialog>
            <IconButton color='secondary' aria-label='edit' size='small' onClick={handleClickOpen}>
                <EditIcon />
            </IconButton>
        </>
    )
}
