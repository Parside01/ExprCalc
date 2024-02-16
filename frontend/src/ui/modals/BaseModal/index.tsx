import { createPortal } from "react-dom"
import "./index.scss"

interface BaseModalProps {
    changer: any
}

export const BaseModal = ({ changer }: BaseModalProps) => {
    return (
        <>
            {createPortal(
                (
                    <div onClick={() => changer(false)} className="modal">
                        <div onClick={(e) => e.stopPropagation()} className="modal__content">
                            <span className="title">Настройки скорости операций</span>
                            <div className="settings">
                                <label>
                                    <span>Скорость выполнения умножения (в мс)</span>
                                    <input type="range" />
                                </label>
                                <label>
                                    <span>Скорость выполнения возведения в степень (в мс)</span>
                                    <input type="range" />
                                </label>
                                <label>
                                    <span>Скорость выполнения деления (в мс)</span>
                                    <input type="range" />
                                </label>
                                <label>
                                    <span>Скорость выполнения вычитания (в мс)</span>
                                    <input type="range" />
                                </label>
                                <label>
                                    <span>Скорость выполнения сложения (в мс)</span>
                                    <input type="range" />
                                </label>
                                <label>
                                    <span>Скорость выполнения деления с остатком (в мс)</span>
                                    <input type="range" />
                                </label>
                                <label>
                                    <span>Скорость выполнения целочисленного деления (в мс)</span>
                                    <input type="range" />
                                </label>
                            </div>
                        </div>
                    </div>
                ),
                document.body
            )}
        </>
    )
}