import { useState } from "react";
import "./index.scss"
import SpeedModal from "../modals/SpeedModal";
import { ExpressionsModal } from "../modals/ExpressionsModal";
import { useExpressionsStore } from "src/store/expressionsStrore";
import { useSettingsStore } from "src/store/settingsStore";
import { ObserveWorkersModal } from "../modals/ObserveWorkersModal";

export const App = () => {
    const [expression, setExpression] = useState("");

    const [viewSpeedModal, setViewSpeedModal] = useState(false);
    const [viewWorkersModal, setViewWorkersModal] = useState(false);
    const [viewTrakingModal, setViewTrakingModal] = useState(false);

    const expressionsStore = useExpressionsStore()
    const settingsStore = useSettingsStore()
    const speedSettings = settingsStore

    return (
        <div className="container">
            <form action="#" className="container__form">
                <h1>Calculator</h1>
                <h3>Enter yout expression</h3>

                <div className="form__inputbox">
                    <input 
                        value={expression} 
                        onChange={(e) => setExpression(e.target.value)}
                        className="expression" placeholder="For example: 2 + 2" 
                    />
                </div>
                
                <div className="options">
                    <button 
                        onClick={() => setViewSpeedModal(!viewSpeedModal)}
                    >
                        <span className="material-symbols-rounded">speed</span>
                    </button>
                    <button
                        onClick={() => setViewWorkersModal(!viewWorkersModal)}
                    >
                        <span className="material-symbols-rounded">settings_backup_restore</span>
                    </button>
                    <button
                        onClick={() => setViewTrakingModal(!viewTrakingModal)}
                    >
                        <span className="material-symbols-rounded">show_chart</span>
                    </button>
                </div>

                <div className="form__btn">
                    <button
                        onClick={async () => (await expressionsStore.sendExpression(expression, speedSettings), setExpression(""))}
                        className="form__submit"
                    >Submit</button>
                </div>
            </form>
            { viewSpeedModal && <SpeedModal changer={setViewSpeedModal} /> }
            { viewWorkersModal && <ExpressionsModal changer={setViewWorkersModal} /> }
            { viewTrakingModal && <ObserveWorkersModal changer={setViewTrakingModal} /> }

        </div>
    )
}