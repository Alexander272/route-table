import React from "react"
import { Container } from "@mui/material"
import { Find } from "./components/Find/Find"
import { OrderTable } from "./components/Table/Table"

export default function Home() {
    return (
        <Container sx={{ marginTop: 5 }}>
            <Find />
            <OrderTable />
        </Container>
    )
}
