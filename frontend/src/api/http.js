import axios from 'axios'
import Qs from 'qs'
import { useMainStore } from '../store'

import { router } from '../router'

// axios 配置
// axios.defaults.timeout = 60000
// const axios = axios.create();
axios.defaults.baseURL = import.meta.env.VITE_BASE_URL


axios.defaults.paramsSerializer = (params) => Qs.stringify(params, { arrayFormat: 'repeat' })

// http request 拦截器
axios.interceptors.request.use(
    config => {
        const store = useMainStore()
        if (store.token) {
            if (!config.headers) {
                config.headers = {};
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
                        window.$message?.error(data?.detail)
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
                    window.$message?.error('permission denied');
                    break
                default:
                    let detail = data?.detail
                    let message = ''
                    if (!detail) {
                        window.$message?.error('Unknown Error');
                        break
                    }
                    if (typeof detail == 'string') {
                        message = detail
                    } else {
                        message = detail[0]['msg']
                    }
                    window.$message?.error(message)
            }
            return Promise.reject(error.response.data)
        } else {
            window.$message?.error('Request Error');
        }
        return Promise.reject(error)
    },
)
