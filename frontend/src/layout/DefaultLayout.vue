<template>
  <a-layout position="absolute">
    <a-layout-header class="h-56px" bordered>
      <Header />
    </a-layout-header>
    <a-layout has-sider position="absolute" style="top: 56px;">
      <a-layout-sider bordered :width="240" collapse-mode="width" :collapsed-width="64"
        v-model:collapsed="store.collapsed" :native-scrollbar="false">
        <a-menu :collapsed="store.collapsed" @update:value="handleUpdateValue">
          <MenuItems :routes="routes" />
        </a-menu>
      </a-layout-sider>
      <a-layout-content content-style="padding: 16px; height: 100%" embedded :native-scrollbar="false">
        <!-- <router-view :key="$route.fullPath" /> -->
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in" appear>
            <component :is="Component" :key="route.fullPath" />
          </transition>
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup lang="ts">
import { h, ref } from 'vue'
import type { SetupContext, VNode, RendererNode, RendererElement } from 'vue'

import { useRouter, useRoute } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

import { MenuItem, SubMenu } from '@arco-design/web-vue'
import { useMainStore } from '@/stores'

import Header from './common/Header.vue'
import icons from '@/assets/icons'
import { routes } from "@/router"


// const expandedKeys: string[] = []

const router = useRouter()
const route = useRoute()

const active = ref<string | null>(null)
active.value = route.name as string
const store = useMainStore()



function handleUpdateValue(key: string) {
  router.push({ name: key })
}

function hasChildren(item: RouteRecordRaw) {
  return Boolean(item.children && item.children.length)
}

function checkPermission(item: RouteRecordRaw, userLevel: string) {
  if (!item.meta?.roles || item.meta?.roles?.includes(userLevel)) {
    return true
  }
  return false
}

function setIcon(iconName: string) {
  const icon = icons[iconName]
  if (icon) {
    return () => h(icon)
  }
}

interface Props {
  routes: RouteRecordRaw[]
  role?: string
}
type Node = VNode<RendererNode, RendererElement, { [key: string]: any; }>
function MenuItems(props: Readonly<Props>, context: SetupContext) {
  const menuItems: Node[] = []
  if (!props.routes) {
    return menuItems
  }
  function travel(_route: RouteRecordRaw[], nodes: Node[]) {
    for (let i of props.routes) {
      if (i.meta?.hideInMenu || !checkPermission(i, props.role as string)) {
        continue
      }

      if (hasChildren(i)) {
        if (i.meta?.flatChildrenInMenu) {
          for (let child of i.children as RouteRecordRaw[]) {
            const item = h(MenuItem, { key: child.name as string }, { icon: () => setIcon(<string>child.meta?.icon), default: () => { child.meta?.locale } })
            nodes.push(item)
          }
        } else {
          const children: never[] = []
          travel(i.children as RouteRecordRaw[], children)
          const item = h(SubMenu, { key: i.name as string }, { icon: () => setIcon(<string>i.meta?.icon), default: () => children })
          nodes.push(item)
        }

      }
      else {
        const item = h(MenuItem, { key: i.name as string }, { icon: () => setIcon(<string>i.meta?.icon), default: () => { i.meta?.locale } })
        nodes.push(item)
      }
    }
  }

  return menuItems
}



</script>


<style scoped lang="less">
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
