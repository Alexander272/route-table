import { List, ListItem, Typography } from "@mui/material"
import React, { FC } from "react"
import { IOperation } from "../../../../types/operations"

type Props = {
    operations: IOperation[]
}

export const OperList: FC<Props> = ({ operations }) => {
    return (
        <List dense sx={{ marginY: 1, width: "100%", maxWidth: "700px" }}>
            {operations?.map(s => (
                <ListItem key={s.id}>
                    <Typography sx={{ flexBasis: "70%" }}>{s.title}</Typography>
                    <Typography sx={{ fontSize: 16, flexBasis: "30%" }} color='primary'>
                        Осталось: {s.remainder}
                    </Typography>
                </ListItem>
            ))}
        </List>
    )
}
