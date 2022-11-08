import {
    Box,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
} from "@mui/material"
import React, { FC } from "react"
import { useNavigate } from "react-router-dom"

type Props = {
    order: number
}

const rows = [
    {
        id: "1",
        position: 1,
        count: 15,
        title: 'Прокладка СПН-Д-24"-300#-316-FG-316 ASME B 16.20',
        ring: "внутрен.",
        deadline: "03.09.2022",
        lastOperation: "03 Резка основания на готовый размер",
        curOperation: "04 Сварка",
    },
    {
        id: "2",
        position: 1,
        count: 15,
        title: 'Прокладка СПН-Д-24"-300#-316-FG-316 ASME B 16.20',
        ring: "наружное",
        deadline: "03.09.2022",
        lastOperation: "04 Сварка",
        curOperation: "05 Зачистка шва",
    },
    {
        id: "3",
        position: 4,
        count: 10,
        title: "Прокладка СНП-В-3-433-1,0-3,2 ОСТ 26.260.454",
        ring: "внутрен.",
        deadline: "03.09.2022",
        lastOperation: "02 Гибка на ребро",
        curOperation: "03 Резка основания на готовый размер",
    },
    {
        id: "4",
        position: 9,
        count: 24,
        title: "Прокладка СНП-В-3-29-1,6-3,2 ОСТ 26.260.454",
        ring: "внутрен.",
        deadline: "03.09.2022",
        lastOperation: "07 Нарезание канавки",
        curOperation: "08 Маркировка",
    },
    {
        id: "5",
        position: 16,
        count: 13,
        title: 'СНП-Д-1/2"-150#-F.G ASME B 16.20',
        ring: "внутрен.",
        deadline: "03.09.2022",
        lastOperation: "",
        curOperation: "01 Лазерная резка",
    },
    {
        id: "6",
        position: 16,
        count: 13,
        title: 'СНП-Д-1/2"-150#-F.G ASME B 16.20',
        ring: "наружное",
        deadline: "03.09.2022",
        lastOperation: "",
        curOperation: "01 Лазерная резка",
    },
]

export const OrderTable: FC<Props> = ({ order }) => {
    const navigate = useNavigate()

    const navigateToPositionHandler = (id: string) => (event: any) => {
        navigate(`positions/${id}`)
    }

    return (
        <Box sx={{ marginTop: 3 }}>
            <Typography variant='h5' component='h5' sx={{ textAlign: "center", marginBottom: 1 }}>
                Заказ №{order}
            </Typography>
            <TableContainer sx={{ maxHeight: 540 }}>
                <Table stickyHeader sx={{ backgroundColor: "#fff" }}>
                    <TableHead>
                        <TableRow>
                            <TableCell>№</TableCell>
                            <TableCell>Количество</TableCell>
                            <TableCell>Наименование</TableCell>
                            <TableCell>Ограничительное кольцо</TableCell>
                            <TableCell>Срок выполнения</TableCell>
                            <TableCell>Последная выполенная операция</TableCell>
                            <TableCell>Текущая операция</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {rows.map(row => (
                            <TableRow key={row.id} onClick={navigateToPositionHandler(row.id)}>
                                <TableCell>{row.position}</TableCell>
                                <TableCell>{row.count}</TableCell>
                                <TableCell>{row.title}</TableCell>
                                <TableCell>{row.ring}</TableCell>
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
