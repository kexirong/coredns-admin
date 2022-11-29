import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/login/index.vue'
import { RouteRecordRaw } from 'vue-router'
import { etcd } from './routes/etcd'
const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    redirect: '/etcd',
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  etcd,

]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes
})

export {
  routes,
  router
} 
