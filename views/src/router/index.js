import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '@/views/home'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    redirect: '/records',
    component: Home,
    children: [
      {
        path: 'records',
        name: 'Records',
        component: () => import(/* webpackChunkName: "records" */ '../views/records.vue')
      },
      {
        path: 'domains',
        name: 'Domains',
        component: () => import(/* webpackChunkName: "domains" */ '../views/domains.vue')
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import(/* webpackChunkName: "login" */ '../views/login.vue')

  }

]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})
router.beforeEach((to, from, next) => {
  window.from = from

  if (to.meta.title) {
    document.title = to.meta.title
  } else {
    document.title = to.name
  }
  if (to.matched.length === 0) {
    from.path ? next({ path: from.path }) : next('/')
    return
  }
  if (to.path !== '/login' && !(localStorage.jwtToken)) {
    next({ name: 'Login', query: { next: to.path } })
  } else {
    next()
  }
})
export default router
