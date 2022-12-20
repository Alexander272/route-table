import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Stack,
    TextField,
} from "@mui/material"
import React, { useEffect, useState } from "react"
import { Controller, useForm } from "react-hook-form"
import useSWR, { useSWRConfig } from "swr"
import { fetcher } from "../../service/read"
import { IUrgency } from "../../types/urgency"
import { changeUrgency } from "../../service/urgency"

export default function Setting({ className }: { className: string }) {
    const [open, setOpen] = useState(false)
    const { data: res } = useSWR<{ data: IUrgency }>("/urgency", fetcher)
    const { mutate } = useSWRConfig()
    const { control, setValue, handleSubmit } = useForm<IUrgency>()

    useEffect(() => {
        if (res) {
            setValue("high", res.data.high)
            setValue("middle", res.data.middle)
        }
    }, [setValue, res])

    const handleClickOpen = () => {
        setOpen(true)
    }

    const handleClose = () => {
        setOpen(false)
    }

    const handleSave = async (data: IUrgency) => {
        try {
            const newData: IUrgency = {
                high: +data.high,
                middle: +data.middle,
            }
            await changeUrgency(newData)
            mutate(`/urgency`)
            handleClose()
        } catch (error) {
            console.log(error)
        }
    }

    return (
        <>
            <Dialog open={open} onClose={handleClose}>
                <DialogTitle>Настройки срочности</DialogTitle>
                <DialogContent>
                    {/* <DialogContentText>
                        To subscribe to this website, please enter your email address here. We will
                        send updates occasionally.
                    </DialogContentText> */}
                    <Stack>
                        <Controller
                            control={control}
                            name='high'
                            render={({ field }) => (
                                <TextField
                                    label='Высокая срочность (в часах)'
                                    type='number'
                                    value={field.value}
                                    onChange={field.onChange}
                                    sx={{ width: 320, marginTop: "20px" }}
                                    InputLabelProps={{
                                        shrink: true,
                                    }}
                                    inputProps={{
                                        inputMode: "numeric",
                                        min: 1,
                                    }}
                                />
                            )}
                        />
                        <Controller
                            control={control}
                            name='middle'
                            render={({ field }) => (
                                <TextField
                                    label='Средняя срочность (в часах)'
                                    type='number'
                                    value={field.value}
                                    onChange={field.onChange}
                                    sx={{ width: 320, marginTop: "20px" }}
                                    InputLabelProps={{
                                        shrink: true,
                                    }}
                                    inputProps={{
                                        inputMode: "numeric",
                                        min: 1,
                                    }}
                                />
                            )}
                        />
                    </Stack>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>Отмена</Button>
                    <Button onClick={handleSubmit(handleSave)}>Сохранить</Button>
                </DialogActions>
            </Dialog>
            <div className={className} onClick={handleClickOpen}>
                <img src='/image/setting.svg' alt='setting' width='30' height='30' />
            </div>
        </>
    )
}
