import {
    Alert,
    AlertColor,
    Box,
    Button,
    Divider,
    Snackbar,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TextField,
    Typography,
} from "@mui/material"
import React, { FC, useCallback, useContext, useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import useSWR from "swr"
import { AuthContext } from "../../../../context/AuthProvider"
import { OrderContext } from "../../../../context/order"
import { useDebounce } from "../../../../hooks/debounse"
import { orderParse } from "../../../../service/order"
import { fetcher } from "../../../../service/read"
import { IOrder } from "../../../../types/order"
import { IPositionForOrder } from "../../../../types/positions"
import { Edit } from "./Edit"

const messages = {
    success: {
        type: "success" as AlertColor,
        message: "Заказ успешно добавлен",
    },
    error: {
        type: "error" as AlertColor,
        message: "Возникла ошибка",
    },
}

type Props = {
    // orderId: string
}

export const OrderTable: FC<Props> = () => {
    const { user } = useContext(AuthContext)
    const { orderId } = useContext(OrderContext)
    const [positions, setPositions] = useState<IPositionForOrder[]>([])
    const [search, setSearch] = useState("")
    const [open, setOpen] = useState(false)
    const [type, setType] = useState<"success" | "error">("success")
    const searchValue = useDebounce(search, 500)

    const { data: order } = useSWR<{ data: IOrder }>(orderId ? `/orders/${orderId}` : null, fetcher)

    useEffect(() => {
        if (order) setPositions(order.data.positions)
    }, [order])

    const navigate = useNavigate()

    const navigateToPositionHandler = (id: string) => (event: any) => {
        navigate(`positions/${id}`)
    }

    const filterHandler = useCallback(
        (search: string) => {
            if (search)
                setPositions(
                    order?.data.positions.filter(p => p.position.toString() === search) || []
                )
            else setPositions(order?.data.positions || [])
        },
        [order]
    )

    useEffect(() => {
        filterHandler(searchValue)
    }, [searchValue, filterHandler])

    const searchHandler = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setSearch(event.target.value)
    }

    const uploadHandler = async (event: any) => {
        const files = event.target.files
        if (!files) return

        const data = new FormData()
        data.append("order", files[0])

        try {
            await orderParse(data)
            handleClick("success")
        } catch (error) {
            handleClick("error")
        }
    }

    const handleClick = (type: AlertColor) => {
        setOpen(true)
        setType(type as "success")
    }
    const handleClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === "clickaway") {
            return
        }

        setOpen(false)
    }

    const chooseBackroundRow = (done: boolean, ring: string) => {
        if (done) {
            return "var(--green)"
        } else {
            return "var(--white)"
        }
        // if (ring === "наружное") {
        //     return "var(--pale-red)"
        // } else {
        //     return "var(--white)"
        // }
    }

    if (!order)
        return (
            <Box
                sx={{
                    marginTop: 3,
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                }}
            >
                <Snackbar
                    open={open}
                    anchorOrigin={{ vertical: "top", horizontal: "center" }}
                    autoHideDuration={6000}
                    onClose={handleClose}
                >
                    <Alert
                        onClose={handleClose}
                        severity={messages[type].type}
                        sx={{ width: "100%" }}
                    >
                        {messages[type].message}
                    </Alert>
                </Snackbar>
                {user?.role === "master" && (
                    <Stack
                        direction={{ xs: "column", sm: "row" }}
                        divider={<Divider orientation='vertical' flexItem />}
                        spacing={{ xs: 0, sm: 2, md: 4 }}
                        sx={{ marginY: 2, marginX: [1, "130px"] }}
                    >
                        <Button variant='contained' component='label' onChange={uploadHandler}>
                            Загрузить заказ
                            <input hidden type='file' />
                        </Button>
                    </Stack>
                )}
            </Box>
        )

    return (
        <Box sx={{ marginTop: 3, display: "flex", flexDirection: "column" }}>
            <Snackbar
                open={open}
                anchorOrigin={{ vertical: "top", horizontal: "center" }}
                autoHideDuration={6000}
                onClose={handleClose}
            >
                <Alert onClose={handleClose} severity={messages[type].type} sx={{ width: "100%" }}>
                    {messages[type].message}
                </Alert>
            </Snackbar>
            <Stack
                direction={{ xs: "column", sm: "row" }}
                divider={<Divider orientation='vertical' flexItem />}
                spacing={{ xs: 0, sm: 2, md: 4 }}
                sx={{ marginY: 2, marginX: [1, "130px"] }}
            >
                <TextField
                    id='filter'
                    label='№ Позиции'
                    size='small'
                    variant='outlined'
                    autoComplete='off'
                    sx={{ background: "var(--white)" }}
                    value={search}
                    onChange={searchHandler}
                />

                <Typography variant='h5' component='h5' sx={{ textAlign: "center" }}>
                    Заказ №{order?.data.number}
                </Typography>

                {user?.role === "master" && (
                    <Button variant='contained' component='label' onChange={uploadHandler}>
                        Загрузить заказ
                        <input hidden type='file' />
                    </Button>
                )}

                {user?.role === "master" && <Edit order={order.data} />}
            </Stack>

            <TableContainer sx={{ maxHeight: 680 }}>
                <Table stickyHeader sx={{ backgroundColor: "#fff" }}>
                    <TableHead>
                        <TableRow>
                            <TableCell>№</TableCell>
                            <TableCell>Наименование</TableCell>
                            <TableCell>Количество</TableCell>
                            <TableCell>Ограничительное кольцо</TableCell>
                            <TableCell>Срок выполнения</TableCell>
                            <TableCell>Последная выполенная операция</TableCell>
                            <TableCell>Текущая операция</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {positions.map(row => (
                            <TableRow
                                key={row.id}
                                onClick={navigateToPositionHandler(row.id)}
                                sx={{ background: chooseBackroundRow(row.done, row.ring) }}
                            >
                                <TableCell>{row.position}</TableCell>
                                <TableCell>{row.title}</TableCell>
                                <TableCell>{row.count}</TableCell>
                                <TableCell
                                    sx={{
                                        fontWeight: row.ring === "наружное" ? "700" : "400",
                                        fontSize: row.ring === "наружное" ? "16px" : "14px",
                                    }}
                                >
                                    {row.ring}
                                </TableCell>
                                <TableCell>{row.deadline}</TableCell>
                                <TableCell>{row.lastOperation}</TableCell>
                                <TableCell>{row.curOperation}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </Box>
    )
}
