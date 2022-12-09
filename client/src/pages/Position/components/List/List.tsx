import { List, ListItem, Stack, Typography } from "@mui/material"
import React, { FC } from "react"
import { IOperation } from "../../../../types/operations"

type Props = {
    operations: IOperation[]
}

export const OperList: FC<Props> = ({ operations }) => {
    const isFinish = operations[operations.length - 1].done

    return (
        <List dense sx={{ marginY: 1, width: "100%", maxWidth: "800px" }}>
            {operations?.map(o => (
                <ListItem key={o.id}>
                    {o.reasons ? (
                        <>
                            <Typography sx={{ flexBasis: "40%" }}>{o.title}</Typography>
                            <Typography sx={{ fontSize: 16, flexBasis: "25%" }} color='primary'>
                                Осталось: {o.remainder}
                            </Typography>
                            <Stack sx={{ flexBasis: "35%" }}>
                                {o.reasons.map(r => (
                                    <Typography key={r.id}>
                                        {r.value} {r.date}
                                    </Typography>
                                ))}
                            </Stack>
                        </>
                    ) : (
                        <>
                            <Typography sx={{ flexBasis: "70%" }}>{o.title}</Typography>
                            <Typography
                                sx={{ fontSize: 16, flexBasis: "30%" }}
                                color={isFinish ? "green" : o.done ? "red" : "primary"}
                            >
                                Осталось: {o.remainder}
                            </Typography>
                        </>
                    )}
                </ListItem>
            ))}
        </List>
    )
}
