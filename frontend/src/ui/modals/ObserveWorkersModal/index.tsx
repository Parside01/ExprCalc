import { WorkerState, useExpressionsStore } from "src/store/expressionsStrore"
import { BaseModal } from "../BaseModal"
import { useEffect } from "react"
import "./index.scss"

interface ObserveWorkersModalProps {
    changer: React.Dispatch<React.SetStateAction<boolean>>
}

export const ObserveWorkersModal = ({ changer }: ObserveWorkersModalProps) => {

    const expressionsStore = useExpressionsStore()

    useEffect(() => {
        expressionsStore.getWorkersInfo()
    }, [])

    return (
        <BaseModal changer={changer}>

            <div className="modal__content">
                <span className="title">
                    Просмотр текущих воркеров
                    <button className="update" onClick={() => expressionsStore.getWorkersInfo()}>
                        <span className="material-symbols-rounded">replay</span>
                    </button>
                </span>
                <div className="workers scroll-bar">
                    { expressionsStore.workers.map((v: any) => (
                        <div key={v.id} className="worker">
                            <div className="left">
                                <div className="id">guid: {v.id}</div>
                                <div className="current_expression">Текущая задача: {v.currentExpression || "Нет"}</div>
                                <div className="last_expression">Предвыдущая задача: {v.lastExpression || "Нет"}</div>
                                <div className="last_touch">Последний запуск: {v.lastTouch.toString()}</div>
                            </div>
                            <div className={`state ${v.state === WorkerState.WORKING && "working"}`}>{v.state === WorkerState.WORKING ? "Работает" : "Отдыхает"}</div>
                        </div>
                    )) }
                </div>
            </div>

        </BaseModal>
    )
}