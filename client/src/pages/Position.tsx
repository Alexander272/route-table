import React from "react"
import {
    Button,
    Container,
    Divider,
    FormControl,
    InputLabel,
    List,
    ListItem,
    MenuItem,
    Paper,
    Select,
    Stack,
    TextField,
    Typography,
} from "@mui/material"
import { useParams } from "react-router-dom"

const operations = [
    { id: "1", title: "01 Лазерная резка", remainder: 0, done: true },
    { id: "2", title: "02 Гибка на ребро", remainder: 0, done: true },
    { id: "3", title: "03 Резка основания на готовый размер", remainder: 0, done: true },
    { id: "4", title: "04 Сварка", remainder: 13, done: false },
    { id: "5", title: "05 Зачистка шва", remainder: 25, done: false },
    { id: "6", title: "06 Полирование шва", remainder: 26, done: false },
    { id: "7", title: "07 Нарезание канавки", remainder: 26, done: false },
    { id: "8", title: "08 Маркировка", remainder: 26, done: false },
    { id: "9", title: "09 Парооксидирование (до 1000 мм)", remainder: 26, done: false },
]

export default function Position() {
    const params = useParams()

    console.log(params)

    const data = {
        order: 8730,
        type: "in",
    }

    const operationId = "4"

    return (
        <Container sx={{ margin: "auto" }}>
            <Paper elevation={3} sx={{ borderRadius: 4, padding: [2, 3] }}>
                <Typography
                    variant='h6'
                    component='h5'
                    sx={{ textAlign: "center", marginBottom: 1, wordBreak: "break-all" }}
                >
                    774,7x685,8x628,7x603,3x4,5-333 СНП-Д-3
                </Typography>
                <Stack
                    direction={{ xs: "column", sm: "row" }}
                    divider={<Divider orientation='vertical' flexItem />}
                    spacing={{ xs: 0, sm: 2, md: 4 }}
                >
                    <Stack direction='row' spacing={2}>
                        <Typography>Заказ</Typography>
                        <Typography sx={{ fontSize: 16 }} color='primary'>
                            № {data.order}
                        </Typography>
                    </Stack>

                    <Stack direction='row' spacing={2}>
                        <Typography>Количество, шт</Typography>
                        <Typography sx={{ fontSize: 16 }} color='primary'>
                            26
                        </Typography>
                    </Stack>
                </Stack>

                <Stack direction='row' spacing={2}>
                    <Typography>Ограничительное кольцо</Typography>
                    <Typography sx={{ fontSize: 16 }} color='primary'>
                        {data.type === "in" ? "внутрен." : "наружное"}
                    </Typography>
                </Stack>

                <List dense>
                    {operations.map(s => (
                        <ListItem key={s.id}>
                            <Stack
                                direction={{ xs: "column", sm: "row" }}
                                spacing={{ xs: 0, sm: 2, md: 4 }}
                            >
                                <Typography>{s.title}</Typography>
                                <Typography sx={{ fontSize: 16 }} color='primary'>
                                    Осталось: {s.remainder}
                                </Typography>
                            </Stack>
                        </ListItem>
                    ))}
                </List>
                <Stack
                    direction={{ xs: "column", sm: "row" }}
                    divider={<Divider orientation='vertical' flexItem />}
                    spacing={{ xs: 1, sm: 2, md: 4 }}
                >
                    <FormControl>
                        <InputLabel id='operation-label'>Операция</InputLabel>
                        <Select
                            labelId='operation-label'
                            id='operation'
                            value={operationId}
                            label='Операция'
                        >
                            {operations.map(o => {
                                if (!o.done) {
                                    return (
                                        <MenuItem key={o.id} value={o.id}>
                                            {o.title}
                                        </MenuItem>
                                    )
                                } else return null
                            })}
                        </Select>
                    </FormControl>
                    <TextField id='count' label='Количество' variant='outlined' />
                    <Button variant='contained'>Сделано</Button>
                </Stack>
            </Paper>
        </Container>
    )
}
