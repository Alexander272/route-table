import { BrowserRouter } from "react-router-dom"
import { AuthProvider } from "./context/AuthProvider"
import { OrderContext } from "./context/order"
import { useOrder } from "./hooks/order"
import Main from "./Main"
import "./index.scss"

function App() {
    return (
        <BrowserRouter>
            <AuthProvider>
                <OrderContext.Provider value={useOrder()}>
                    <Main />
                </OrderContext.Provider>
            </AuthProvider>
        </BrowserRouter>
    )
}

export default App
