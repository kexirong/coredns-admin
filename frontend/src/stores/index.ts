import { defineStore } from 'pinia'
import { AxiosResponse } from 'axios'
import { Base64 } from 'js-base64'
import { login } from '../api/user'


function payloadDecode(authToken: string) {
  if (!authToken) {
    return null
  }
  let payload = authToken.split('.')[1]
  if (!payload) {
    return null
  }
  return JSON.parse(Base64.decode(payload))
}

export const useMainStore = defineStore('main', {
  state: () => {
    return {
      token: "",
      collapsed: false,
      darkMode: false,
    }
  },
  getters: {
    account: state => {
      let user = payloadDecode(state.token)
      return user?.username
    },
    userClaims: state => {
      return payloadDecode(state.token)
    },
    checkAuthToken(state) {
      let ts = new Date().getTime()
      let payload = payloadDecode(state.token)
      if (!payload || !Object.keys(payload).includes('exp')) {
        return false
      }
      return payload.exp * 1000 > ts
    },
  },
  actions: {
    async Login(username: string, password: string) {
      return await login({ username: username.trim(), password })
        .then((response: AxiosResponse) => {
          const data = response.data
          this.token = data.access_token
          return true
        })

    },
  },
  persist: {
    enabled: true,
    strategies: [
      {
        storage: localStorage,
      },
    ],
  }
})