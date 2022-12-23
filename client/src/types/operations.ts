import { IReason } from './reason'

export interface IOperation {
	id: string
	title: string
	done: boolean
	remainder: number
	isFinish: boolean
	date?: string
	reasons?: IReason[]
}

export interface ICompleteOperation {
	id: string
	done: boolean
	remainder: number
	count: number
	reason: string
}
