import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Chip,
    CircularProgress,
    Container,
    Stack,
    Typography,
    Box,
    Divider,
} from "@mui/material"
import React from "react"
import useSWR from "swr"
import { fetcher } from "../../service/read"
import { IGroupedOrder } from "../../types/order"

export default function Orders() {
    const { data: res } = useSWR<{ data: IGroupedOrder[] }>("/api/v1/orders/", fetcher, {
        refreshInterval: 60 * 1000,
    })

    return (
        <Container sx={{ marginTop: 5 }}>
            {res?.data.map(o => (
                <Accordion
                    key={o.id}
                    TransitionProps={{ unmountOnExit: true }}
                    defaultExpanded={true}
                    // disableGutters
                >
                    <AccordionSummary
                        expandIcon={<>&#9660;</>}
                        aria-controls='panel1a-content'
                        id='panel1a-header'
                    >
                        <Stack
                            direction={{ xs: "column", sm: "row" }}
                            alignItems='center'
                            divider={<Divider orientation='vertical' flexItem />}
                            spacing={{ xs: 0, sm: 2 }}
                        >
                            <Typography color='primary' variant='h6'>
                                Дата отгрузки {o.deadline}
                            </Typography>

                            <Stack
                                direction={{ xs: "column", sm: "row" }}
                                spacing={{ xs: 0, sm: 2 }}
                                alignItems='center'
                            >
                                <Typography>Срочность</Typography>
                                <Chip
                                    label={o.urgency}
                                    sx={{ width: "90px" }}
                                    color={
                                        o.urgency === "Высокая"
                                            ? "error"
                                            : o.urgency === "Средняя"
                                            ? "primary"
                                            : "success"
                                    }
                                />
                            </Stack>
                        </Stack>
                    </AccordionSummary>
                    <AccordionDetails>
                        {o.orders.map(o => (
                            <Stack
                                key={o.id}
                                direction={{ xs: "column", sm: "row" }}
                                justifyContent='space-between'
                                alignItems='center'
                                sx={{
                                    paddingY: "4px",
                                    borderBottom: "1px solid var(--primary-color)",
                                }}
                            >
                                <Stack
                                    direction={{ xs: "column", sm: "row" }}
                                    spacing={{ xs: 0, sm: 2 }}
                                    alignItems='center'
                                >
                                    <Typography color='primary'>Заказ №{o.number}</Typography>
                                    <Typography>От {o.date}</Typography>
                                </Stack>

                                <Stack
                                    direction={{ xs: "column", sm: "row" }}
                                    spacing={{ xs: 0, sm: 2 }}
                                    alignItems='center'
                                >
                                    {o.done ? (
                                        <Chip label='Выполнен' color='success' />
                                    ) : (
                                        <Chip label='В работе' color='primary' />
                                    )}
                                    <Box sx={{ position: "relative", display: "inline-flex" }}>
                                        <CircularProgress
                                            size={50}
                                            variant='determinate'
                                            value={o.progress || 0}
                                        />
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
                            </Stack>
                        ))}
                    </AccordionDetails>
                </Accordion>
            ))}
        </Container>
    )
}
