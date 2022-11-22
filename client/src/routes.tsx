import { lazy, Suspense } from "react"
import { Route, Routes } from "react-router-dom"

import RequireAuth from "./components/Auth/RequireAuth"
import MainLayout from "./components/Layout/Main"

const Auth = lazy(() => import("./pages/Auth/Auth"))
const Home = lazy(() => import("./pages/Home/Home"))
const Position = lazy(() => import("./pages/Position/Position"))
const Orders = lazy(() => import("./pages/Orders/Orders"))

export const MyRoutes = () => {
    return (
        <Suspense fallback={null}>
            <Routes>
                <Route path='/auth' element={<Auth />} />

                <Route
                    path='/'
                    element={
                        <RequireAuth>
                            <MainLayout />
                        </RequireAuth>
                    }
                >
                    <Route index element={<Home />} />
                    <Route path='/positions/:id' element={<Position />} />
                    <Route path='/orders' element={<Orders />} />
                </Route>

                <Route path='*' element={<div />} />
            </Routes>
        </Suspense>
    )
}
