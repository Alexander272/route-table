import {
    Button,
    Divider,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    SelectChangeEvent,
    Stack,
    TextField,
} from "@mui/material"
import React, { FC, useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { useSWRConfig } from "swr"
import { operationComplite } from "../../../../service/operation"
import { ICompleteOperation, IOperation } from "../../../../types/operations"
import { ICompletePosition, IPosition } from "../../../../types/positions"

type Props = {
    position: IPosition
    operations: IOperation[]
}

export const Operations: FC<Props> = ({ position, operations }) => {
    const [operationIdx, setOperationIdx] = useState("0")
    const [remainder, setRemainder] = useState(0)
    const [count, setCount] = useState(0)
    const [reason, setReason] = useState("")

    const params = useParams()

    const { mutate } = useSWRConfig()

    useEffect(() => {
        for (let i = 0; i < (operations.length || 0); i++) {
            const o = operations[i]
            if (!o?.done) {
                setOperationIdx(i.toString())
                setRemainder(o?.remainder || 0)
                setCount(o?.remainder || 0)
                break
            }
        }
    }, [operations])

    const countHandler = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        if (+event.target.value > remainder || +event.target.value < 1) return
        setCount(+event.target.value)
    }

    const operationHandler = (event: SelectChangeEvent) => {
        setOperationIdx(event.target.value)
        const op = operations[+event.target.value]
        setCount(op?.remainder || 0)
        setRemainder(op?.remainder || 0)
    }

    const reasonHandler = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setReason(event.target.value)
    }

    const submitHandler = async () => {
        const operation: ICompleteOperation = {
            id: operations[+operationIdx].id || "",
            done: remainder === count,
            remainder: remainder - count,
            reason: reason,
        }

        const data: ICompletePosition = {
            id: position.id || "",
            count: position.count || 0,
            isFinish: +operationIdx === (operations.length || 0) - 1 || false,
            connected: position.connected || "",
            operation: operation,
        }

        await operationComplite(data)
        mutate(`/api/v1/positions/${params.id}`)
    }

    return (
        <Stack
            direction={{ xs: "column", sm: "row" }}
            divider={<Divider orientation='vertical' flexItem />}
            spacing={{ xs: 1, sm: 2, md: 4 }}
        >
            <FormControl>
                <InputLabel id='operation-label'>Операция</InputLabel>
                <Select
                    sx={{ minWidth: 220 }}
                    labelId='operation-label'
                    id='operation'
                    value={operationIdx}
                    label='Операция'
                    size='small'
                    onChange={operationHandler}
                >
                    {operations?.map((o, i) => {
                        if (!o.done) {
                            return (
                                <MenuItem key={o.id} value={i}>
                                    {o.title}
                                </MenuItem>
                            )
                        } else return null
                    })}
                </Select>
            </FormControl>
            <TextField
                sx={{ minWidth: 150 }}
                id='count'
                label='Количество'
                variant='outlined'
                type='number'
                value={count}
                onChange={countHandler}
                size='small'
                inputProps={{
                    inputMode: "numeric",
                    min: 1,
                    max: remainder,
                    // pattern: "[0-9]*",
                }}
            />
            {count < remainder && (
                <TextField
                    sx={{ minWidth: 150 }}
                    id='reason'
                    label='Причина'
                    variant='outlined'
                    size='small'
                    value={reason}
                    onChange={reasonHandler}
                />
            )}
            <Button variant='contained' onClick={submitHandler}>
                Сделано
            </Button>
        </Stack>
    )
}
