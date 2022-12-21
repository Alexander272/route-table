import React, { useContext, useState } from 'react'
import {
	Box,
	Container,
	Divider,
	Paper,
	Stack,
	Typography,
	CircularProgress,
	IconButton,
	Snackbar,
	Alert,
} from '@mui/material'
import CloseIcon from '@mui/icons-material/Close'
import { useNavigate, useParams } from 'react-router-dom'
import useSWR from 'swr'
import { IPosition } from '../../types/positions'
import { fetcher } from '../../service/read'
import { OperList } from './components/List/List'
import { Operations } from './components/Operations/Operations'
import { AuthContext } from '../../context/AuthProvider'
import { Edit } from './components/Edit/Edit'

export default function Position() {
	const params = useParams()
	const navigate = useNavigate()

	const { user } = useContext(AuthContext)

	const [open, setOpen] = useState(false)
	const [error, setError] = useState('')

	const { data: position } = useSWR<{ data: IPosition }>(`/positions/${params.id}`, fetcher)

	const backHandler = () => {
		navigate(-1)
	}

	const changeErrorHandler = (error: string) => {
		setError(error)
		if (error !== '') setOpen(true)
	}

	const handleClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
		if (reason === 'clickaway') {
			return
		}
		setOpen(false)
	}

	return (
		<Container sx={{ flexGrow: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
			{!position && (
				<Box sx={{ display: 'flex', justifyContent: 'center' }}>
					<CircularProgress />
				</Box>
			)}
			<Snackbar
				open={open}
				anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
				autoHideDuration={6000}
				onClose={handleClose}
			>
				<Alert onClose={handleClose} severity='error' sx={{ width: '100%' }}>
					{error}
				</Alert>
			</Snackbar>
			{position && (
				<Paper
					elevation={3}
					sx={{
						marginTop: 2,
						borderRadius: 4,
						paddingX: [2, 4],
						paddingY: [2, 3],
						flexGrow: 1,
						display: 'flex',
						flexDirection: 'column',
						alignItems: 'center',
						position: 'relative',
					}}
				>
					<IconButton onClick={backHandler} sx={{ position: 'absolute', right: '3px', top: '3px' }}>
						<CloseIcon />
					</IconButton>
					<Typography
						variant='h5'
						component='h5'
						color='primary'
						sx={{ textAlign: 'center', marginBottom: 2, wordBreak: 'break-all' }}
					>
						{position.data.title} {user?.role === 'master' && <Edit position={position?.data} />}
					</Typography>
					<Stack
						direction={{ xs: 'column', sm: 'row' }}
						divider={<Divider orientation='vertical' flexItem />}
						spacing={{ xs: 0, sm: 2, md: 4 }}
						alignItems='center'
					>
						<Stack direction='row' spacing={2} alignItems='center'>
							<Typography>Заказ/Позиция</Typography>
							<Typography sx={{ fontSize: 20 }} color='primary'>
								№ {position.data.order}/{position.data.position}
							</Typography>
						</Stack>

						<Stack direction='row' spacing={2}>
							<Typography>Количество, шт</Typography>
							<Typography sx={{ fontSize: 16 }} color='primary'>
								{position.data.count}
							</Typography>
						</Stack>

						{position.data.ring && (
							<Stack direction='row' spacing={2}>
								<Typography>Ограничительное кольцо</Typography>
								<Typography sx={{ fontSize: 16 }} color='primary'>
									{position.data.ring}
								</Typography>
							</Stack>
						)}
					</Stack>

					<OperList
						position={position?.data}
						operations={position?.data?.operations || []}
						count={position.data.count}
						changeErrorHandler={changeErrorHandler}
					/>

					{!position.data.done && user?.role !== 'manager' ? (
						<Operations
							position={position.data}
							operations={position?.data?.operations.filter(o => !o.done) || []}
							changeErrorHandler={changeErrorHandler}
						/>
					) : null}
				</Paper>
			)}
		</Container>
	)
}
