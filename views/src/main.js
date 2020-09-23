import Vue from 'vue'
import App from './App.vue'
import router from './router'
import '@/assets/cutom.scss'
import './plugins/element.js'

Vue.config.productionTip = false
Vue.prototype.$base_url = 'http://localhost:8088'
new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
