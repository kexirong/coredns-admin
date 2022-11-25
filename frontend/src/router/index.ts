import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/login/index.vue'
import { RouteRecordRaw } from 'vue-router'
const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'login',
    component: Login
  },

]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes
})

export {
  routes,
  router
} 
