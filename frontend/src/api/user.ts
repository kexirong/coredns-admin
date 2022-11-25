import axios from 'axios'

export function login(userInfo: { username: string, password: string }) {
    return axios.post('/login', userInfo, { headers: { 'Content-Type': 'multipart/form-data' } })
}