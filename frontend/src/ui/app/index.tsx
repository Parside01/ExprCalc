import { useEffect, useState } from "react";
import "./index.scss"
import { api } from "src/utils/axios";

const sendExpression = (expr: string) => {
    api.post("/expr/calc", { expr }).catch(e => console.log(e))

    console.log(expr)
}

export const App = () => {
    const [expression, setExpression] = useState("");

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
                    <button><span className="material-symbols-rounded">speed</span></button>
                    <button><span className="material-symbols-rounded">settings</span></button>
                    <button><span className="material-symbols-rounded">settings_backup_restore</span></button>
                    <button><span className="material-symbols-rounded">show_chart</span></button>
                </div>

                <div className="form__btn">
                    <button
                        onClick={() => sendExpression(expression)}
                        className="form__submit"
                    >Submit</button>
                </div>
            </form>
        </div>
    )
}