import { Box, IconButton, InputBase, Paper } from "@mui/material"
import React, { FC } from "react"

type Props = {}

export const Find: FC<Props> = () => {
    return (
        <Box sx={{ display: "flex", justifyContent: "center" }}>
            <Paper
                component='form'
                sx={{ p: "2px 4px", display: "flex", alignItems: "center", width: 600 }}
            >
                <InputBase
                    sx={{ ml: 1, flex: 1 }}
                    placeholder='Введите номер заказа'
                    inputProps={{ "aria-label": "Введите номер заказа" }}
                />
                <IconButton type='button' sx={{ p: "6px" }} aria-label='search'>
                    {/* <SearchIcon />  */}
                    &#128269;
                </IconButton>
            </Paper>
        </Box>
    )
}
