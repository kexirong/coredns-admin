import Vue from 'vue'
import App from './App.vue'
import router from './router'
import '@/assets/cutom.scss'
import './plugins/element.js'
import { decode } from 'js-base64'

Vue.config.productionTip = false
Vue.prototype.$base_url = 'http://localhost:8088'
Vue.prototype.$jwtToken = localStorage.jwtToken

Vue.prototype.$payloadDecode = function (authToken) {
  console.log('================================')
  console.log(this.$jwtToken)
  console.log('================================')
  if (!authToken || authToken.constructor !== String) {
    return null
  }
  const payload = authToken.split('.')[1]
  if (!payload) {
    return null
  }
  return JSON.parse(decode(payload))
}

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
