import React, { useContext } from "react"
import { Container } from "@mui/material"
import { Find } from "./components/Find/Find"
import { OrderTable } from "./components/Table/Table"
import { AuthContext } from "../../context/AuthProvider"
import { Navigate } from "react-router-dom"

export default function Home() {
    const { user } = useContext(AuthContext)

    if (user?.role === "display") return <Navigate to='/orders/group' />

    return (
        <Container sx={{ marginTop: 5 }}>
            <Find />
            <OrderTable />
        </Container>
    )
}
