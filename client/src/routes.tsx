import { lazy, Suspense } from "react"
import { Route, Routes } from "react-router-dom"

const Home = lazy(() => import("./pages/Home/Home"))
const Position = lazy(() => import("./pages/Position/Position"))
const Orders = lazy(() => import("./pages/Orders/Orders"))

export const MyRoutes = () => {
    return (
        <Suspense fallback={null}>
            <Routes>
                <Route path='/' element={<Home />} />
                <Route path='/positions/:id' element={<Position />} />
                <Route path='/orders' element={<Orders />} />

                <Route path='*' element={<div />} />
            </Routes>
        </Suspense>
    )
}
