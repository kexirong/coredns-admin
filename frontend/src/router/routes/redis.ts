
import type { RouteRecordRaw } from 'vue-router'

import DefaultLayout from '@/layout/DefaultLayout.vue'


export const redis: RouteRecordRaw = {
    path: '/redis',
    name: 'Redis',
    redirect: {name:'RedisDomains'},
    component: DefaultLayout,
    meta: {
        locale: 'Redis',
        icon: 'Redis'
    },
    children: [
      {
        path: 'domains',
        name: 'RedisDomains',
        component: () => import('@/views/redis-domain/index.vue'),
        meta: {
            locale: 'Domains',
            icon: 'Parse'
        },
      }
    ]
  }