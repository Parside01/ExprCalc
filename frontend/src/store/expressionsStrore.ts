import { api } from "src/utils/axios";
import { create } from "zustand";
import { SpeedExpressionsSettings } from "./settingsStore";

export enum ExpressionState {
    COMPLETED = "COMPLETED",
    PENDING = "PENDING"
}

export enum WorkerState {
    WAITING, WORKING
}

export interface Expression {
    id: string // guid
    state: ExpressionState // Завершено или нет
    expression: string // операция
    result: string // результат опреации 
    executionTime: string // время выполнения в мс 
}

export interface Worker {
    id: string // guid
    state: WorkerState // состояние отыхает или работает
    lastExpression: string // последнее выражение
    currentExpression: string // текущее выражение
    lastTouch: Date // когда закончил свою последнюю работу
}

export interface ExpressionsState {
    expressions: Expression[],
    workers: Worker[]
    getExpressions({ take, skip }: { take: number, skip: number }): Promise<void>
    sendExpression(expression: string, settings: SpeedExpressionsSettings): Promise<void>
    getWorkersInfo(): Promise<void>
}

export const useExpressionsStore = create<ExpressionsState>((set) => ({
    expressions: [],
    workers: [],

    async getExpressions({ take, skip }: { take: number, skip: number }) {
        try {
            const res = await api.get("/expr/getAllExpressions", {params: {take, skip}})
            const data: Expression[] = res.data.map((v: any) => ({
                state: v["is-done"] ? ExpressionState.COMPLETED : ExpressionState.PENDING,
                id: v.guid,
                expression: v.expression,
                result: v.result,
                executionTime: v["execute-time"]
            }))
            
            set((state) => ({
                expressions: data.reduce((acc: Expression[], cur: Expression) => {
                    const existExpression = acc.find(v => v.id === cur.id)

                    if (existExpression === undefined) acc.push(cur)
                    if (existExpression !== undefined && cur.state !== existExpression.state) existExpression.state = cur.state
                    if (existExpression !== undefined && cur.executionTime !== existExpression.executionTime) existExpression.executionTime = cur.executionTime
                    if (existExpression !== undefined && cur.result !== existExpression.result) existExpression.result = cur.result
                    
                    return acc
                }, state.expressions)
            }))
        } catch(e) {
            throw e
        }
    },

    async sendExpression(
        expression: string, 
        settings: SpeedExpressionsSettings = {
            multiplicationSpeed: 200,
            divisionSpeed: 200,
            exponentiationSpeed: 200,
            subtractionSpeed: 200,
            additionSpeed: 200,
            divisionWithRemainderSpeed: 200
        }
    ) {
        try {
            await api.post("/expr/calc", { 
                expr: expression, 
                "*": settings.multiplicationSpeed,
                "/": settings.divisionSpeed,
                "**": settings.exponentiationSpeed,
                "-": settings.subtractionSpeed,
                "+": settings.additionSpeed,
                "%": settings.divisionWithRemainderSpeed
            })
        } catch (e) {
            throw e
        }
    },

    async getWorkersInfo() {
        try {
            const res = await api.get("/expr/getWorkersInfo")
            const data: Worker[] = res.data.map((v: any): Worker => ({
                state: v["is-employ"] ? WorkerState.WORKING : WorkerState.WAITING,
                id: v["worker-id"],
                lastExpression: v["prev-job"],
                currentExpression: v["current-job"],
                lastTouch: new Date(v["last-touch"])
            }))
            
            set((state) => ({
                workers: data.reduce((acc: Worker[], cur: Worker) => {
                    const existWorker = acc.find(v => v.id === cur.id)

                    if (existWorker === undefined) acc.push(cur)
                    if (existWorker !== undefined && cur.state !== existWorker.state) existWorker.state = cur.state
                    if (existWorker !== undefined && cur.currentExpression !== existWorker.currentExpression) existWorker.currentExpression = cur.currentExpression
                    if (existWorker !== undefined && cur.lastExpression !== existWorker.lastExpression) existWorker.lastExpression = cur.lastExpression
                    if (existWorker !== undefined && cur.lastTouch !== existWorker.lastTouch) existWorker.lastTouch = cur.lastTouch
                    
                    return acc
                }, state.workers)
            }))

            console.log(...this.workers)
        } catch(e) {
            throw e
        }
    } 
}))