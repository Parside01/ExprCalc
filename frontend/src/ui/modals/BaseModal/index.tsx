import { createPortal } from "react-dom"
import "./index.scss"

interface BaseModalProps {
    changer: React.Dispatch<React.SetStateAction<boolean>>,
    children: React.ReactNode
}

export const BaseModal = ({ changer, children }: BaseModalProps) => {
    return (
        <>
            {createPortal(
                (
                    <div onClick={() => changer(false)} className="modal">
                        <div onClick={(e) => e.stopPropagation()}>
                            { children }
                        </div>
                    </div>
                ),
                document.body
            )}
        </>
    )
}