import { BrowserRouter } from "react-router-dom"
import { MyRoutes } from "./routes"
import "./index.css"

function App() {
    return (
        <BrowserRouter>
            <MyRoutes />
        </BrowserRouter>
    )
}

export default App
