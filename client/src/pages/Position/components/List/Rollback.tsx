import { Button, Dialog, DialogActions, DialogTitle, IconButton } from '@mui/material'
import RemoveIcon from '@mui/icons-material/RemoveCircleOutline'
import React, { FC, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useSWRConfig } from 'swr'
import { IOperation } from '../../../../types/operations'
import { IPosition, IRollbackPosition } from '../../../../types/positions'
import { rollbackOperation } from '../../../../service/operation'

type Props = {
	position: IPosition
	operation: IOperation
	changeErrorHandler: (error: string) => void
}

export const Rollback: FC<Props> = ({ position, operation, changeErrorHandler }) => {
	const params = useParams()
	const { mutate } = useSWRConfig()

	const [open, setOpen] = useState(false)

	const handleClickOpen = () => {
		setOpen(true)
	}

	const handleClose = () => {
		setOpen(false)
	}

	const rollbackHandler = async () => {
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
		handleClose()
		mutate(`/positions/${params.id}`)
	}

	return (
		<>
			<Dialog
				open={open}
				onClose={handleClose}
				aria-labelledby='alert-dialog-title'
				aria-describedby='alert-dialog-description'
			>
				<DialogTitle id='alert-dialog-title'>{'Отменить операцию?'}</DialogTitle>

				<DialogActions>
					<Button onClick={handleClose}>Отмена</Button>
					<Button onClick={rollbackHandler} autoFocus>
						ОК
					</Button>
				</DialogActions>
			</Dialog>
			<IconButton aria-label='delete' size='small' sx={{ padding: 0, marginLeft: 1 }} onClick={handleClickOpen}>
				<RemoveIcon color='error' />
			</IconButton>
		</>
	)
}
