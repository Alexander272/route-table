import {
    Box,
    FormControl,
    IconButton,
    InputAdornment,
    InputBase,
    OutlinedInput,
    Paper,
} from "@mui/material"
import React, { FC } from "react"

type Props = {}

export const Find: FC<Props> = () => {
    return (
        <Box sx={{ display: "flex", justifyContent: "center" }}>
            {/* <OutlinedInput
                type='search'
                size='small'
                endAdornment={
                    <InputAdornment position='end'>
                        <IconButton type='button' sx={{ p: "6px" }} aria-label='search'>
                            &#128269;
                        </IconButton>
                    </InputAdornment>
                }
                placeholder='Введите номер заказа'
                fullWidth
            /> */}
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
