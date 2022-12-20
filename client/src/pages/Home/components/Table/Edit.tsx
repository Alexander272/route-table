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
import React, { FC, useContext, useEffect, useState } from "react"
import { IOrder } from "../../../../types/order"
import { orderDelete, orderUpdate } from "../../../../service/order"
import { useSWRConfig } from "swr"
import { OrderContext } from "../../../../context/order"

type Props = {
    order: IOrder
}

export const Edit: FC<Props> = ({ order }) => {
    const [open, setOpen] = useState(false)
    const [openDelete, setOpenDelete] = useState(false)
    const [deadline, setDeadline] = useState(Date.now())
    const { mutate } = useSWRConfig()
    const { changeOrderId } = useContext(OrderContext)

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

    const handleOpenDelete = () => {
        setOpenDelete(true)
    }
    const handleCloseDelete = () => {
        setOpenDelete(false)
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

    const handleDelete = async () => {
        try {
            await orderDelete(order.id)
            changeOrderId("")
            handleCloseDelete()
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
                    <Button color='info' onClick={handleClose}>
                        Отмена
                    </Button>
                    <Button color='error' onClick={handleOpenDelete}>
                        Удалить
                    </Button>
                    <Button onClick={handleSave}>Сохранить</Button>
                </DialogActions>
            </Dialog>
            <Dialog open={openDelete} onClose={handleCloseDelete}>
                <DialogTitle>Удалить заказ?</DialogTitle>
                <DialogActions>
                    <Button onClick={handleCloseDelete}>Отмена</Button>
                    <Button onClick={handleDelete}>Удалить</Button>
                </DialogActions>
            </Dialog>
            <IconButton color='secondary' aria-label='edit' size='small' onClick={handleClickOpen}>
                <EditIcon />
            </IconButton>
        </>
    )
}
