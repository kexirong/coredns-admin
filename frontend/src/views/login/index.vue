<template>
    <div class="w-full h-full flex-center">
        <a-card class="w-480px" hoverable :bordered="true">
            <template #cover>
                <div class="flex-x-center mt-8px">
                    <img alt="dessert" src="@/assets/CoreDNS_Colour_Horizontal.png" />
                </div>
            </template>
            <a-card-meta title="用户登录">
                <template #description>
                    <a-form :model="formModel" :layout="layout" auto-label-width @submit="handleSubmit" class="mt-8px">
                        <a-form-item field="username" label="Username:">
                            <a-input v-model="formModel.username" placeholder="please enter your username..." />
                        </a-form-item>
                        <a-form-item field="password" label="Password:">
                            <a-input v-model="formModel.password" type="password"
                                placeholder="please enter your password..." />
                        </a-form-item>
                        <a-form-item>
                            <a-button type="primary" html-type="submit" long :loading="loading">登录</a-button>
                        </a-form-item>
                    </a-form>
                </template>
            </a-card-meta>
        </a-card>
    </div>
</template>
  
<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ValidatedError, Message } from '@arco-design/web-vue'
import { useMainStore } from '@/stores'
import { useRoute, useRouter } from 'vue-router'

const layout = ref<"inline" | "horizontal" | "vertical">('horizontal')
const formModel = reactive({
    username: '',
    password: ''
})
const loading = ref(false)
const { Login } = useMainStore()

const router = useRouter()
const route = useRoute()

interface Data {
    values: Record<string, any>
    errors: Record<string, ValidatedError> | undefined
}

function handleSubmit({ values, errors }: Data, ev: Event): any {
    if (errors) return false;

    loading.value = true
    Login(values.username, values.password)
        .then((ret) => {
            if (ret) {
                Message.success("登录成功")
                router.push(route.query.redirect as string || "/")
            }
        })
        .finally(() => {
            loading.value = false
        })


}

</script>