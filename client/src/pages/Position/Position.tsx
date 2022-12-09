import React from "react"
import { Box, Container, Divider, Paper, Stack, Typography, CircularProgress } from "@mui/material"
import { useParams } from "react-router-dom"
import useSWR from "swr"
import { IPosition } from "../../types/positions"
import { fetcher } from "../../service/read"
import { OperList } from "./components/List/List"
import { Operations } from "./components/Operations/Operations"

export default function Position() {
    const params = useParams()

    const { data: position } = useSWR<{ data: IPosition }>(`/positions/${params.id}`, fetcher)

    return (
        <Container
            sx={{ flexGrow: 1, display: "flex", alignItems: "center", justifyContent: "center" }}
        >
            {!position && (
                <Box sx={{ display: "flex", justifyContent: "center" }}>
                    <CircularProgress />
                </Box>
            )}
            {position && (
                <Paper
                    elevation={3}
                    sx={{
                        marginTop: 2,
                        borderRadius: 4,
                        paddingX: [2, 4],
                        paddingY: [2, 3],
                        flexGrow: 1,
                        display: "flex",
                        flexDirection: "column",
                        alignItems: "center",
                    }}
                >
                    <Typography
                        variant='h5'
                        component='h5'
                        color='primary'
                        sx={{ textAlign: "center", marginBottom: 2, wordBreak: "break-all" }}
                    >
                        {position.data.title}
                    </Typography>
                    <Stack
                        direction={{ xs: "column", sm: "row" }}
                        divider={<Divider orientation='vertical' flexItem />}
                        spacing={{ xs: 0, sm: 2, md: 4 }}
                        alignItems='center'
                    >
                        <Stack direction='row' spacing={2} alignItems='center'>
                            <Typography>Заказ/Позиция</Typography>
                            <Typography sx={{ fontSize: 20 }} color='primary'>
                                № {position.data.order}/{position.data.position}
                            </Typography>
                        </Stack>

                        <Stack direction='row' spacing={2}>
                            <Typography>Количество, шт</Typography>
                            <Typography sx={{ fontSize: 16 }} color='primary'>
                                {position.data.count}
                            </Typography>
                        </Stack>

                        <Stack direction='row' spacing={2}>
                            <Typography>Ограничительное кольцо</Typography>
                            <Typography sx={{ fontSize: 16 }} color='primary'>
                                {position.data.ring}
                            </Typography>
                        </Stack>
                    </Stack>

                    <OperList operations={position?.data?.operations || []} />

                    {!position.data.done ? (
                        <Operations
                            position={position.data}
                            operations={position?.data?.operations.filter(o => !o.done) || []}
                        />
                    ) : null}
                </Paper>
            )}
        </Container>
    )
}
