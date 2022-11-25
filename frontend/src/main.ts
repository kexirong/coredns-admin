import { createApp } from 'vue'
import { createPinia } from 'pinia'

import piniaPersist from 'pinia-plugin-persist'

import App from './App.vue'
import { router } from './router'
import installArcoDesign from './plugin/arco-design'

import './api/interceptor'
import 'uno.css'


const app = createApp(App)

installArcoDesign(app)

const pinia = createPinia()
pinia.use(piniaPersist)
app.use(pinia).use(router).mount('#app')
