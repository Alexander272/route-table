import { Box, CircularProgress, Stack, Typography } from '@mui/material'
import React, { FC, useContext } from 'react'
import { useNavigate } from 'react-router-dom'
import { OrderContext } from '../../context/order'
import { IGroupedOrder, IOrderItem } from '../../types/order'
import classes from './orders.module.scss'

type Props = {
	order: IGroupedOrder
}

export const OrderItem: FC<Props> = ({ order }) => {
	let urgency = 'high'
	if (order.urgency === 'Высокая') urgency = 'high'
	else if (order.urgency === 'Средняя') urgency = 'middle'
	else urgency = ''

	const { changeOrderId } = useContext(OrderContext)
	const navigate = useNavigate()

	const selectOrder = (order: IOrderItem) => () => {
		changeOrderId(order.id)
		navigate('/')
	}

	const chooseColor = (progress: number) => {
		if (progress < 1) return 'primary'
		if (progress > 1 && progress < 70) {
			return 'warning'
		} else {
			return 'error'
		}
	}

	return (
		<div className={[classes.item, urgency ? classes[urgency] : null].join(' ')}>
			<Stack alignItems='center' spacing={1}>
				<Typography color='primary' variant='h5' sx={{ fontWeight: 700 }}>
					Дата отгрузки {order.deadline}
				</Typography>
			</Stack>

			{order.orders.map(o => (
				<Stack
					key={o.id}
					direction={{ xs: 'column', sm: 'row' }}
					spacing={{ xs: 0, sm: 2 }}
					alignItems='center'
					justifyContent='space-between'
					sx={{
						paddingY: '4px',
						borderBottom: '1px solid var(--primary-color)',
					}}
				>
					<Stack spacing={0} onClick={selectOrder(o)} sx={{ cursor: 'pointer' }}>
						<Typography variant='h6' color='primary' sx={{ fontWeight: 700 }}>
							Заказ №{o.number}
						</Typography>
						{/* <Typography>От {o.date}</Typography> */}
						<Typography>{o.customer}</Typography>
					</Stack>

					<Box sx={{ position: 'relative', display: 'inline-flex' }}>
						<CircularProgress
							size={50}
							variant='determinate'
							value={o.progress || 0}
							color={chooseColor(o.progress)}
						/>
						<Box
							sx={{
								top: 0,
								left: 0,
								bottom: 0,
								right: 0,
								position: 'absolute',
								display: 'flex',
								alignItems: 'center',
								justifyContent: 'center',
							}}
						>
							<Typography variant='caption' component='div' color='text.secondary' fontSize={12}>{`${
								o.progress || 0
							}%`}</Typography>
						</Box>
					</Box>
				</Stack>
			))}
		</div>
	)
}
