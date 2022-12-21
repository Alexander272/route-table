import { ICompletePosition, IRollbackPosition } from '../types/positions'
import api from './api'

export const operationComplete = async (data: ICompletePosition) => {
	try {
		await api.put(`/operations/${data.id}`, data)
	} catch (error) {
		console.log(error)
	}
}

export const rollbackOperation = async (operationId: string, pos: IRollbackPosition) => {
	try {
		await api.post(`/operations/rollback/${operationId}`, pos)
	} catch (error: any) {
		throw error.response
	}
}
