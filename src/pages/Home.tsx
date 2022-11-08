import { Container } from "@mui/material"
import React from "react"
import { Find } from "../components/Find"
import { OrderTable } from "../components/Table"

export default function Home() {
    return (
        <Container sx={{ marginTop: 5 }}>
            <Find />
            <OrderTable order={8730} />
        </Container>
    )
}
