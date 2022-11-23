import axios from 'axios'

function login(userInfo) {
    return axios.post('/auth', userInfo, { headers: { 'Content-Type': 'multipart/form-data' } })
}