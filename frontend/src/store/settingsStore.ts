import { create } from "zustand"

export interface SpeedExpressionsSettings {
    multiplicationSpeed: number
    divisionSpeed: number
    exponentiationSpeed: number
    subtractionSpeed: number
    additionSpeed: number
    divisionWithRemainderSpeed: number
}

interface SettingsState {
    setMultiplicationSpeed(speed: number): void
    setDivisionSpeed(speed: number): void
    setExponentiationSpeed(speed: number): void
    setSubtractionSpeed(speed: number): void
    setAdditionSpeed(speed: number): void
    setDivisionWithRemainderSpeed(speed: number): void
}

export const useSettingsStore = create<SettingsState & SpeedExpressionsSettings>((set) => ({
    multiplicationSpeed: 200,
    divisionSpeed: 200,
    exponentiationSpeed: 200,
    subtractionSpeed: 200,
    additionSpeed: 200,
    divisionWithRemainderSpeed: 200,

    setAdditionSpeed: (speed: number) => set({ additionSpeed: speed }),
    setDivisionSpeed: (speed: number) => set({ divisionSpeed: speed }),
    setExponentiationSpeed: (speed: number) => set({ exponentiationSpeed: speed }),
    setSubtractionSpeed: (speed: number) => set({ subtractionSpeed: speed }),
    setMultiplicationSpeed: (speed: number) => set({ multiplicationSpeed: speed }),
    setDivisionWithRemainderSpeed: (speed: number) => set({ divisionWithRemainderSpeed: speed })
}))