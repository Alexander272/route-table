import { List, ListItem, Stack, Typography } from '@mui/material'
import React, { FC, useContext } from 'react'
import { AuthContext } from '../../../../context/AuthProvider'
import { IOperation } from '../../../../types/operations'
import { IPosition } from '../../../../types/positions'
import { Rollback } from './Rollback'

type Props = {
	position: IPosition
	operations: IOperation[]
	count: number
	changeErrorHandler: (error: string) => void
}

export const OperList: FC<Props> = ({ position, operations, count, changeErrorHandler }) => {
	const { user } = useContext(AuthContext)

	const isFinish = operations[operations.length - 1].done

	return (
		<List dense sx={{ marginY: 1, width: '100%', maxWidth: '800px' }}>
			{operations?.map(o => (
				<ListItem key={o.id}>
					{o.reasons ? (
						<>
							<Typography sx={{ flexBasis: '40%' }}>{o.title}</Typography>
							<Stack direction={'row'} alignItems='center' sx={{ flexBasis: '25%' }}>
								<Typography
									sx={{ fontSize: 16 }}
									color={isFinish ? 'green' : o.done ? 'red' : 'primary'}
								>
									Осталось: {o.remainder}
									<br />
									{o.date}
								</Typography>
								{o.remainder < count && user?.role !== 'manager' ? (
									<Rollback
										position={position}
										operation={o}
										changeErrorHandler={changeErrorHandler}
									/>
								) : null}
							</Stack>
							<Stack sx={{ flexBasis: '35%' }}>
								{o.reasons.map(r => (
									<Typography key={r.id}>
										{r.value} {r.date}
									</Typography>
								))}
							</Stack>
						</>
					) : (
						<>
							<Typography sx={{ flexBasis: '70%' }}>{o.title}</Typography>
							<Stack direction={'row'} alignItems='center' sx={{ flexBasis: '30%' }}>
								<Typography
									sx={{ fontSize: 16 }}
									color={isFinish ? 'green' : o.done ? 'red' : 'primary'}
								>
									Осталось: {o.remainder}
									<br />
									{o.date}
								</Typography>
								{o.remainder < count && user?.role !== 'manager' ? (
									<Rollback
										position={position}
										operation={o}
										changeErrorHandler={changeErrorHandler}
									/>
								) : null}
							</Stack>
						</>
					)}
				</ListItem>
			))}
		</List>
	)
}
