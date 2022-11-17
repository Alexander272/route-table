import { BrowserRouter } from "react-router-dom"
import { MyRoutes } from "./routes"
import { OrderContext } from "./context/order"
import { useOrder } from "./hooks/order"
import "./index.scss"

function App() {
    return (
        <BrowserRouter>
            <OrderContext.Provider value={useOrder()}>
                <MyRoutes />
            </OrderContext.Provider>
        </BrowserRouter>
    )
}

export default App
