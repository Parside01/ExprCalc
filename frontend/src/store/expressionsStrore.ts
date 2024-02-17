import { api } from "src/utils/axios";
import { create } from "zustand";
import { SpeedExpressionsSettings, useSettingsStore } from "./settingsStore";

export enum WorkerState {
    COMPLETED = "COMPLETED",
    PENDING = "PENDING"
}

export interface Expression {
    id: string // guid
    state: WorkerState // Завершено или нет
    expression: string // операция
    result: string // результат опреации 
    executionTime: string // время выполнения в мс 
}

export interface ExpressionsState {
    expressions: Expression[],
    workers: Worker[]
    getExpressions({ take, skip }: { take: number, skip: number }): Promise<void>
    sendExpression(expression: string, settings: SpeedExpressionsSettings): Promise<void>
    getWorkersInfo({ take, skip }: { take: number, skip: number }): Promise<void>
}

export const useExpressionsStore = create<ExpressionsState>((set) => ({
    expressions: [],
    workers: [],

    async getExpressions({ take, skip }: { take: number, skip: number }) {
        try {
            const res = await api.get("/expr/getAllExpressions", {params: {take, skip}})
            const data: Expression[] = res.data.map((v: any) => ({
                state: v["is-done"] ? WorkerState.COMPLETED : WorkerState.PENDING,
                id: v.guid,
                expression: v.expression,
                result: v.result,
                executionTime: v["execute-time"]
            }))
            
            set((state) => ({
                expressions: data.reduce((acc: Expression[], cur: Expression) => {
                    if (acc.find(v => v.id === cur.id) === undefined) acc.push(cur)
                    
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

    async getWorkersInfo({ take, skip }: { take: number, skip: number }) {
        try {
            const res = await api.get("/expr/getWorkersInfo", {params: {take, skip}})

            this.workers.concat(res.data)
        } catch (e) {
            throw e
        }
    } 
}))