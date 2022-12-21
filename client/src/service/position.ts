import { IUpdateCount } from '../types/positions'
import api from './api'

export const updateCount = async (data: IUpdateCount) => {
	try {
		await api.put(`/positions/${data.id}/count`, data)
	} catch (error: any) {
		throw error.response.message
	}
}
