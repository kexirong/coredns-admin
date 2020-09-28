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
        name: 'Record',
        component: () => import(/* webpackChunkName: "records" */ '../views/records.vue')
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import(/* webpackChunkName: "records" */ '../views/login.vue')

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
  }

  console.log(Vue.prototype.$jwtToken)
  if (to.path !== '/login' && (Vue.prototype.$jwtToken === undefined || Vue.prototype.$jwtToken === 'undefined')) {
    next({ name: 'Login', query: { next: to.path } })
  } else {
    next()
  }
})
export default router
