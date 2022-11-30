<template>
    <a-row>
        <a-col :span="6">
            <a-tree blockNode :data="treeData" :fieldNames="{ title: 'name', children: 'subdomain' }" size="large"
                @select="onTreeSelect" />
        </a-col>
        <a-col :span="18">
            <a-button type="primary" class="mb-8px" @click="handleAdd"> 添加记录 </a-button>
            <a-table :columns="columns" :bordered="false" :data="tableData" :pagination="pagination" :loading="loading">
                <template #type="{ record }">
                    <a-select v-if="record.action == 'add'" :options="recordTypeOptions" v-model="record.type"
                        class="w-88px" />
                    <template v-else>{{ record.type }}</template>
                </template>
                <template #ttl="{ record }">
                    <a-input-number v-if="record.action" v-model="record.ttl" hide-button class="w-56px" />
                    <template v-else>{{ record.ttl }}</template>
                </template>
                <template #priority="{ record }">
                    <a-input-number v-if="record.action" v-model="record.priority" hide-button class="w-56px" />
                    <template v-else>{{ record.priority }}</template>
                </template>
                <template #name="{ record }">
                    <a-input v-if="record.action == 'add'" v-model="record.name" />
                    <template v-else>{{ record.name }}</template>
                </template>
                <template #content="{ record }">
                    <a-input v-if="record.action" v-model="record.content" />
                    <template v-else>{{ record.content }}</template>
                </template>

                <template #optional="{ record, rowIndex }">
                    <a-space v-if="record.action">
                        <a-button @click="onRecordSave(record)" type="primary" status="warning">提交</a-button>
                        <a-button @click="onCancel(record, rowIndex)" type="primary">取消</a-button>
                    </a-space>
                    <a-space v-else>
                        <a-button @click="record.action = 'edit'" type="primary" status="success">编辑</a-button>
                        <a-popconfirm content="删除后不可恢复" type="warning" @ok="onRecordDelete(record, rowIndex)">
                            <a-button type="primary" status="danger">删除</a-button>
                        </a-popconfirm>
                    </a-space>
                </template>
            </a-table>
        </a-col>
    </a-row>
</template>
<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { getDomains, DomainsData, getRecord, getRecords, RecordData, postRecord, putRecord, deleteRecord } from '@/api/etcd'
import type { TreeNodeData, TableData } from '@arco-design/web-vue/es'
import { Message } from '@arco-design/web-vue'
import { encodeURI } from 'js-base64'
const treeData = ref<TreeNodeData[]>([])
const tableData = ref<TableData[]>([])
const pagination = { showPageSize: true }
const loading = ref(true)
const columns = [
    {
        title: 'TYPE',
        dataIndex: 'type',
        slotName: 'type'
    },
    {
        title: 'TTL',
        dataIndex: 'ttl',
        slotName: 'ttl'
    },
    {
        title: 'PRIORITY',
        dataIndex: 'priority',
        slotName: 'priority'
    },
    {
        title: 'NAME',
        dataIndex: 'name',
        slotName: 'name'
    },
    {
        title: 'CONTENT',
        dataIndex: 'content',
        slotName: 'content'
    },
    {
        title: '操作',
        slotName: 'optional'
    }
]

const recordTypeOptions = [
    { label: "A", value: "A" },
    { label: "NS", value: "NS", disabled: true },
    { label: "CNAME", value: "CNAME" },
    { label: "PTR", value: "PTR" },
    { label: "MX", value: "MX" },
    { label: "TXT", value: "TXT" },
    { label: "AAA", value: "AAA" },
    { label: "SRV", value: "SRV" },
]
function dataAddKey(data: DomainsData[], key: string) {
    for (const i of data) {
        if (!i.key) {
            i.key = key + '/' + i.name
        }

        if (i.subdomain) {
            dataAddKey(i.subdomain, i.key)
        }
    }
    return data
}
getDomains()
    .then(res => {
        const subdomain = res.data.subdomain
        if (subdomain) {
            treeData.value = dataAddKey(subdomain, '')
        }
    })
getRecords()
    .then(res => {
        tableData.value = res.data
        loading.value = false
    })
function handleAdd() {
    if (tableData.value[0]?.action) {
        return
    }
    tableData.value.unshift({ action: 'add' })
}
function onCancel(record: RecordData, rowIndex: number) {
    if (record.action === 'add') {
        tableData.value.splice(rowIndex, 1)
    } else {
        record.action = undefined
    }
}
function onTreeSelect(keys: (string | number)[]) {
    loading.value = true
    const key = encodeURI(keys[0] as string)
    getRecord(key)
        .then(res => {
            tableData.value = res.data
        })
        .finally(
            () => { loading.value = false }
        )
}
function onRecordSave(record: RecordData) {
    console.log(record)

    switch (record.action) {
        case 'edit':
            const key = encodeURI(<string>record.key)
            putRecord(key, record)
                .then((res) => {
                    record.key = res.data.key
                    Message.success('添加成功')
                    record.action = undefined
                })
            break
        case 'add':
            postRecord(record)
                .then((res) => {
                    record.key = res.data.key
                    Message.success('保存成功')
                    record.action = undefined
                })
            break
        default:
            Message.warning('数据异常')
    }

}
function onRecordDelete(record: RecordData, rowIndex: number) {
    const key = encodeURI(record.key as string)
    deleteRecord(key)
        .then(() => {
            tableData.value.splice(rowIndex, 1)
            Message.success('删除成功')
        })
}
</script>