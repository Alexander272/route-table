import axios from "axios"

const api = axios.create({
    withCredentials: true,
    baseURL: "/api/v1",
})

export default api
