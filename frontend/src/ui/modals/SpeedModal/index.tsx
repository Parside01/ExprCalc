import { useSettingsStore } from "src/store/settingsStore"
import { BaseModal } from "../BaseModal"
import "./index.scss"

interface SpeedModalProps {
    changer: React.Dispatch<React.SetStateAction<boolean>>,
}

export default function SpeedModal({ changer }: SpeedModalProps) {
    const settingsStore = useSettingsStore()

    return (
        <BaseModal changer={changer}>
            <div className="modal__content"> 
                <span className="title">Настройки скорости операций</span>
                <div className="settings">
                    <label>
                        <span>Скорость выполнения умножения (в мс). Текущая: {settingsStore.multiplicationSpeed}</span>
                        <input 
                            type="range" 
                            value={settingsStore.multiplicationSpeed} 
                            max={10000}
                            onChange={(e) => settingsStore.setMultiplicationSpeed(parseInt(e.target.value))} 
                        />
                    </label>
                    <label>
                        <span>Скорость выполнения возведения в степень (в мс). Текущая: {settingsStore.exponentiationSpeed}</span>
                        <input 
                            type="range"
                            value={settingsStore.exponentiationSpeed} 
                            max={10000}
                            onChange={(e) => settingsStore.setExponentiationSpeed(parseInt(e.target.value))} 
                        />
                    </label>
                    <label>
                        <span>Скорость выполнения деления (в мс). Текущая: {settingsStore.divisionSpeed}</span>
                        <input 
                            type="range"
                            value={settingsStore.divisionSpeed} 
                            max={10000}
                            onChange={(e) => settingsStore.setDivisionSpeed(parseInt(e.target.value))} 
                        />
                    </label>
                    <label>
                        <span>Скорость выполнения вычитания (в мс). Текущая: {settingsStore.subtractionSpeed}</span>
                        <input 
                            type="range"
                            value={settingsStore.subtractionSpeed} 
                            max={10000}
                            onChange={(e) => settingsStore.setSubtractionSpeed(parseInt(e.target.value))} 
                        />
                    </label>
                    <label>
                        <span>Скорость выполнения сложения (в мс). Текущая: {settingsStore.additionSpeed}</span>
                        <input 
                            type="range"
                            value={settingsStore.additionSpeed} 
                            max={10000}
                            onChange={(e) => settingsStore.setAdditionSpeed(parseInt(e.target.value))} 
                        />
                    </label>
                    <label>
                        <span>Скорость выполнения деления с остатком (в мс). Текущая: {settingsStore.divisionWithRemainderSpeed}</span>
                        <input 
                            type="range"
                            value={settingsStore.divisionWithRemainderSpeed} 
                            max={10000}
                            onChange={(e) => settingsStore.setDivisionWithRemainderSpeed(parseInt(e.target.value))} 
                        />
                    </label>
                </div>
            </div>
        </BaseModal>
    )
}