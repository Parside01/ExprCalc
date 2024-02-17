import "./index.scss"
import { BaseModal } from "../BaseModal"
import { ExpressionState, useExpressionsStore } from "src/store/expressionsStrore"
import { useEffect } from "react"

interface ExpressionsModalProps {
    changer: React.Dispatch<React.SetStateAction<boolean>>
}

export const ExpressionsModal = ({ changer }: ExpressionsModalProps) => {
    const expressionsStore = useExpressionsStore()

    useEffect(() => {
        expressionsStore.getExpressions({take: 10, skip: 0})
    }, [])

    return (
        <BaseModal changer={changer}>
            <div className="modal__content">
                <span className="title">
                    Все выполняющиеся задачи
                    <button className="update" onClick={() => expressionsStore.getExpressions({take: 10, skip: 0})}>
                        <span className="material-symbols-rounded">replay</span>
                    </button>
                </span>
                <div className="expressions scroll-bar">
                    { expressionsStore.expressions.map((v: any) => (
                        <div key={v.id} className="expression">
                            <div className="id">{v.id}</div>
                            <div className="expression_string">{v.expression} = {v.result}</div>
                            <div className="created">{v.executionTime} ms</div>
                            <div className={`state ${v.state === ExpressionState.COMPLETED && "completed"}`}>{v.state === ExpressionState.COMPLETED ? "Выполнено" : "Выполняется"}</div>
                        </div>
                    )) }
                </div>
            </div>
        </BaseModal>
    )
}