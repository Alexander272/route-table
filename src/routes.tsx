import { lazy, Suspense } from "react"
import { Route, Routes } from "react-router-dom"

const Home = lazy(() => import("./pages/Home"))
const Position = lazy(() => import("./pages/Position"))

export const MyRoutes = () => {
    return (
        <Suspense fallback={null}>
            <Routes>
                <Route path='/' element={<Home />} />
                <Route path='/positions/:id' element={<Position />} />

                <Route path='*' element={<div />} />
            </Routes>
        </Suspense>
    )
}
