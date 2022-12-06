import axios from 'axios'

import { useMainStore } from '../stores'

import { router } from '../router'
import { Message } from '@arco-design/web-vue'

// axios 配置
// axios.defaults.timeout = 60000
// const axios = axios.create()
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL

// http request 拦截器

axios.interceptors.request.use(
    config => {
        const store = useMainStore()

        if (store.token) {
            if (!config.headers) {
                config.headers = {}
            }
            config.headers.Authorization = `Bearer ${store.token}`
        }
        return config
    },
    error => {
        return Promise.reject(error)
    },
)



axios.interceptors.response.use(
    response => {
        return response
    },
    error => {
        if (error.response) {

            const { status, data } = error.response

            switch (status) {
                case 401:
                    if (router.currentRoute.value.path == '/login') {
                        Message?.error(data?.reason)
                    } else {
                        router.replace({
                            name: 'Login',
                            query: { redirect: router.currentRoute.value.path }
                        })
                    }

                    break

                case 403:
                    const params = router.resolve(window.location.pathname).params

                    if (params && Object.keys(params).length) {
                        router.back()
                    }
                    Message?.error('permission denied')
                    break
                default:
                    let reason = data?.reason
                    Message?.error(reason || 'Unknown Error')
            }
            return Promise.reject(error.response.data)
        } else {
            Message?.error('Request Error')
        }
        return Promise.reject(error)
    },
)
