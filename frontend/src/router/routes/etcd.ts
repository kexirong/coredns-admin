
import type { RouteRecordRaw } from 'vue-router'

import DefaultLayout from '@/layout/DefaultLayout.vue'


export const etcd: RouteRecordRaw = {
    path: '/etcd',
    name: 'Etcd',
    redirect: {name:'Domains'},
    component: DefaultLayout,
    meta: {
        locale: 'Etcd',
        icon: 'Prometheus'
    },
    children: [
      {
        path: 'domains',
        name: 'Domains',
        component: () => import('@/views/domain/index.vue'),
        meta: {
            locale: 'Domains',
            icon: 'Prometheus'
        },
      }
    ]
  }