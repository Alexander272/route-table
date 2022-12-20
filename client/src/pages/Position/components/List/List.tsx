import { IconButton, List, ListItem, Stack, Typography } from "@mui/material"
import DoDisturbOnIcon from "@mui/icons-material/DoDisturbOn"
import React, { FC } from "react"
import { IOperation } from "../../../../types/operations"

type Props = {
    operations: IOperation[]
    count: number
    changeErrorHandler: (error: string) => void
}

export const OperList: FC<Props> = ({ operations, count, changeErrorHandler }) => {
    const isFinish = operations[operations.length - 1].done

    const roolbackHandler = (operationId: string) => () => {}

    return (
        <List dense sx={{ marginY: 1, width: "100%", maxWidth: "800px" }}>
            {operations?.map(o => (
                <ListItem key={o.id}>
                    {o.reasons ? (
                        <>
                            <Typography sx={{ flexBasis: "40%" }}>{o.title}</Typography>
                            <Typography sx={{ fontSize: 16, flexBasis: "25%" }} color='primary'>
                                Осталось: {o.remainder}{" "}
                                <IconButton aria-label='delete' size='small'>
                                    <DoDisturbOnIcon />
                                </IconButton>
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
                            <Stack direction={"row"} alignItems='center' sx={{ flexBasis: "30%" }}>
                                <Typography
                                    sx={{ fontSize: 16 }}
                                    color={isFinish ? "green" : o.done ? "red" : "primary"}
                                >
                                    Осталось: {o.remainder}
                                </Typography>
                                {o.remainder < count && (
                                    <IconButton
                                        aria-label='delete'
                                        size='small'
                                        sx={{ padding: 0, marginLeft: 1 }}
                                        onClick={roolbackHandler(o.id)}
                                    >
                                        <DoDisturbOnIcon color='error' />
                                    </IconButton>
                                )}
                            </Stack>
                        </>
                    )}
                </ListItem>
            ))}
        </List>
    )
}
