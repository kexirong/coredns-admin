
import type { RouteRecordRaw } from 'vue-router'

import DefaultLayout from '@/layout/DefaultLayout.vue'


export const etcd: RouteRecordRaw = {
    path: '/etcd',
    name: 'Etcd',
    redirect: {name:'EtcdDomains'},
    component: DefaultLayout,
    meta: {
        locale: 'Etcd',
        icon: 'Etcd'
    },
    children: [
      {
        path: 'domains',
        name: 'EtcdDomains',
        component: () => import('@/views/etcd-domain/index.vue'),
        meta: {
            locale: 'Domains',
            icon: 'Parse'
        },
      }
    ]
  }