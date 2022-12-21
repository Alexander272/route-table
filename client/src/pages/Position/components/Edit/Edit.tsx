import { Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, TextField } from '@mui/material'
import EditIcon from '@mui/icons-material/Edit'
import React, { FC, useEffect, useState } from 'react'
import { useSWRConfig } from 'swr'
import { IPosition } from '../../../../types/positions'
import { updateCount } from '../../../../service/position'

type Props = {
	position: IPosition
}

export const Edit: FC<Props> = ({ position }) => {
	const [open, setOpen] = useState(false)
	const [count, setCount] = useState(position.count)
	const { mutate } = useSWRConfig()

	useEffect(() => {
		// if (order) setDeadline(+order.deadline * 1000)
	}, [position])

	const handleCount = (event: React.ChangeEvent<HTMLInputElement>) => {
		setCount(+event.target.value)
	}

	const handleClickOpen = () => {
		setOpen(true)
	}
	const handleClose = () => {
		setOpen(false)
	}

	const handleSave = async () => {
		try {
			const data = { id: position.id, count, done: position.done }
			await updateCount(data)
			mutate(`/positions/${position.id}`)
			handleClose()
		} catch (error) {
			console.log(error)
		}
	}

	return (
		<>
			<Dialog open={open} onClose={handleClose}>
				<DialogTitle>Редактировать позицию</DialogTitle>
				<DialogContent>
					<TextField
						id='count'
						label='Количество, шт'
						type='number'
						value={count}
						onChange={handleCount}
						sx={{ width: 320, marginTop: '20px' }}
						InputLabelProps={{
							shrink: true,
						}}
					/>
				</DialogContent>
				<DialogActions>
					<Button color='info' onClick={handleClose}>
						Отмена
					</Button>
					<Button onClick={handleSave}>Сохранить</Button>
				</DialogActions>
			</Dialog>

			<IconButton color='secondary' aria-label='edit' size='small' onClick={handleClickOpen}>
				<EditIcon />
			</IconButton>
		</>
	)
}
