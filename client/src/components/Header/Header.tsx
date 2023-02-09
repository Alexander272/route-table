import { Tooltip } from '@mui/material'
import React, { FC, useContext } from 'react'
import { Link } from 'react-router-dom'
import { AuthContext } from '../../context/AuthProvider'
import { signOut } from '../../service/auth'
import DownloadFiles from '../DownloadFiles/DownloadFiles'
import Setting from '../Setting/Setting'
import classes from './header.module.scss'

type Props = {}

export const Header: FC<Props> = () => {
	const { user, setUser } = useContext(AuthContext)

	const logoutHandler = async () => {
		try {
			await signOut()
			setUser(null)
		} catch (error) {}
	}

	return (
		<div className={classes.header}>
			<div className={classes.content}>
				<Link className={classes.logoLink} to='/'>
					<img
						className={classes.logo}
						width={192}
						height={192}
						loading='lazy'
						src='/logo192.webp'
						alt='logo'
					/>
					<span>Sealur Route</span>
				</Link>

				{user && (
					<div className={classes.nav}>
						{user?.role === 'master' ? <DownloadFiles className={classes.profile} /> : null}

						{user?.role === 'master' ? <Setting className={classes.profile} /> : null}

						<Tooltip title='Главная'>
							<Link to='/' className={classes.profile}>
								<img src='/image/home.svg' alt='home' width='32' height='32' />
							</Link>
						</Tooltip>

						{user?.role === 'master' || user?.role === 'display' || user?.role === 'manager' ? (
							<Tooltip title='Список заказов'>
								<Link to='/orders' className={classes.profile}>
									<img src='/image/list.svg' alt='orders' width='30' height='30' />
								</Link>
							</Tooltip>
						) : null}

						<Tooltip title='Выход'>
							<div className={classes.profile} onClick={logoutHandler}>
								<img src='/image/logout.svg' alt='log-out' width='30' height='30' />
							</div>
						</Tooltip>
					</div>
				)}
			</div>
		</div>
	)
}
