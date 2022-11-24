import { Box, Chip, CircularProgress, Stack, Typography } from "@mui/material"
import React, { FC } from "react"
import { IGroupedOrder } from "../../types/order"
import classes from "./orders.module.scss"

type Props = {
    order: IGroupedOrder
}

export const OrderItem: FC<Props> = ({ order }) => {
    return (
        <div className={classes.item}>
            <Stack alignItems='center' spacing={1} sx={{ marginBottom: 1 }}>
                <Typography color='primary' variant='h6'>
                    Дата отгрузки {order.deadline}
                </Typography>

                <Stack
                    direction={{ xs: "column", sm: "row" }}
                    spacing={{ xs: 0, sm: 2 }}
                    alignItems='center'
                >
                    <Typography>Срочность</Typography>
                    <Chip
                        label={order.urgency}
                        sx={{ width: "90px" }}
                        color={
                            order.urgency === "Высокая"
                                ? "error"
                                : order.urgency === "Средняя"
                                ? "primary"
                                : "success"
                        }
                    />
                </Stack>
            </Stack>

            {order.orders.map(o => (
                <Stack
                    key={o.id}
                    direction={{ xs: "column", sm: "row" }}
                    spacing={{ xs: 0, sm: 2 }}
                    alignItems='center'
                    sx={{
                        paddingY: "4px",
                        borderBottom: "1px solid var(--primary-color)",
                    }}
                >
                    <Stack spacing={0} alignItems='center'>
                        <Typography color='primary'>Заказ №{o.number}</Typography>
                        <Typography>От {o.date}</Typography>
                    </Stack>

                    <Box sx={{ position: "relative", display: "inline-flex" }}>
                        <CircularProgress size={50} variant='determinate' value={o.progress || 0} />
                        <Box
                            sx={{
                                top: 0,
                                left: 0,
                                bottom: 0,
                                right: 0,
                                position: "absolute",
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                            }}
                        >
                            <Typography
                                variant='caption'
                                component='div'
                                color='text.secondary'
                                fontSize={12}
                            >{`${o.progress || 0}%`}</Typography>
                        </Box>
                    </Box>
                </Stack>
            ))}
        </div>
    )
}
