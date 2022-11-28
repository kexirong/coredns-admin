<template>
    <div class="header-container flex-y-center h-full">
        <div class="w-240px flex-center text-primary">

            <CrisisAlert />

            <span class="font-600 ml-8px">通知分发器
            </span>
        </div>
        <div class="flex-1 flex-y-center cursor-pointer">
            <div @click="store.collapsed = !store.collapsed">
                <MenuFoldRight v-if="store.collapsed" />
                <MenuFoldLeft v-else />
            </div>
        </div>
        <div class="flex-center h-full">
            <dark-mode-switch v-model:dark="store.darkMode" class="w-40px flex-x-center" />
            <a-dropdown @select="handleDropdown" trigger="click">
                <div class="flex-center px-12px cursor-pointer">
                    <DefaultAvatar />
                    <span class="pl-8px text-18px font-500">{{ store.userClaims.name }}</span>
                </div>
                <template #content>
                    <a-doption value="logout">
                        <template #icon>
                            <Logout />
                        </template>
                        <template #default>退出登录</template>
                    </a-doption>
                </template>
            </a-dropdown>
        </div>
    </div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'

import DarkModeSwitch from '@/components/DarkModeSwitch.vue'

import CrisisAlert from '@/assets/icons/CrisisAlert.vue'
import DefaultAvatar from '@/assets/icons/DefaultAvatar.vue'
import MenuFoldLeft from '@/assets/icons/MenuFoldLeft.vue'
import MenuFoldRight from '@/assets/icons/MenuFoldRight.vue'
import Logout from '@/assets/icons/Logout.vue'

import { useMainStore } from '@/stores'


const store = useMainStore()



// const UserAvatar = icons['UserAvatar']
// const Logout = icons['Logout']
// const Theme = icons['Theme']

const router = useRouter()

function handleDropdown(key: string) {
    switch (key) {
        case 'change-password':
            router.push({ name: "UserChangePassword" })
            break
        case 'logout':
            router.push({ name: "Login" })
            break
    }
}



</script>

<style lang="scss" scoped>
.header-container {
    box-shadow: 0 1px 2px rgb(0 21 41 / 8%)
}
</style>