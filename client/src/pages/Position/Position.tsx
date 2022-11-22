import React, { useContext } from "react"
import { Box, Container, Divider, Paper, Stack, Typography, CircularProgress } from "@mui/material"
import { useParams } from "react-router-dom"
import useSWR from "swr"
import { IPosition } from "../../types/positions"
import { fetcher } from "../../service/read"
import { OperList } from "./components/List/List"
import { Operations } from "./components/Operations/Operations"
import { AuthContext } from "../../context/AuthProvider"

export default function Position() {
    const params = useParams()
    const { user } = useContext(AuthContext)

    const { data: position } = useSWR<{ data: IPosition }>(
        `/api/v1/positions/${params.id}`,
        fetcher
    )

    return (
        <Container sx={{ margin: "auto" }}>
            {!position && (
                <Box sx={{ display: "flex", justifyContent: "center" }}>
                    <CircularProgress />
                </Box>
            )}
            {position && (
                <Paper
                    elevation={3}
                    sx={{
                        borderRadius: 4,
                        paddingX: [2, 4],
                        paddingY: [2, 3],
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
                    >
                        <Stack direction='row' spacing={2}>
                            <Typography>Заказ/Позиция</Typography>
                            <Typography sx={{ fontSize: 16 }} color='primary'>
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

                    <OperList operations={position.data.operations || []} />

                    {!position.data.done && user?.role !== "master" ? (
                        <Operations
                            position={position.data}
                            operations={position.data.operations}
                        />
                    ) : null}
                </Paper>
            )}
        </Container>
    )
}
