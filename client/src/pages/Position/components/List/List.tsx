import { IconButton, List, ListItem, Stack, Typography } from '@mui/material'
import RemoveIcon from '@mui/icons-material/RemoveCircleOutline'
import React, { FC } from 'react'
import { IOperation } from '../../../../types/operations'
import { rollbackOperation } from '../../../../service/operation'
import { useParams } from 'react-router-dom'
import { useSWRConfig } from 'swr'
import { IPosition, IRollbackPosition } from '../../../../types/positions'

type Props = {
	position: IPosition
	operations: IOperation[]
	count: number
	changeErrorHandler: (error: string) => void
}

export const OperList: FC<Props> = ({ position, operations, count, changeErrorHandler }) => {
	const params = useParams()
	const { mutate } = useSWRConfig()

	const isFinish = operations[operations.length - 1].done

	const rollbackHandler = (operation: IOperation) => async () => {
		try {
			const pos: IRollbackPosition = {
				id: position.id,
				connected: position.connected,
				reasons: operation.reasons?.map(r => r.id) || [],
				isFinishOperation: operation.isFinish,
			}
			await rollbackOperation(operation.id, pos)
		} catch (error) {
			changeErrorHandler('Отмена операции не доступна')
		}
		mutate(`/positions/${params.id}`)
	}

	return (
		<List dense sx={{ marginY: 1, width: '100%', maxWidth: '800px' }}>
			{operations?.map(o => (
				<ListItem key={o.id}>
					{o.reasons ? (
						<>
							<Typography sx={{ flexBasis: '40%' }}>{o.title}</Typography>
							<Typography sx={{ fontSize: 16, flexBasis: '25%' }} color='primary'>
								Осталось: {o.remainder}{' '}
								<IconButton
									aria-label='delete'
									size='small'
									sx={{ padding: 0, marginLeft: 1 }}
									onClick={rollbackHandler(o)}
								>
									<RemoveIcon color='error' />
								</IconButton>
							</Typography>
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
								</Typography>
								{o.remainder < count && (
									<IconButton
										aria-label='delete'
										size='small'
										sx={{ padding: 0, marginLeft: 1 }}
										onClick={rollbackHandler(o)}
									>
										<RemoveIcon color='error' />
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
