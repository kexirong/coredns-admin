<template>
  <el-row :gutter="5">
    <el-col :span="5">
      <el-input-number size="mini" :min="0" :max="10" v-model="num" @change="change" />
      <el-tree :data="treeData" default-expand-all :props="defaultProps" @node-click="handleNodeClick"
        :expand-on-click-node="false" />
    </el-col>
    <el-col :span="19">
      <el-button type="primary" size="mini" @click="handleAdd">
        添加记录
      </el-button>
      <el-table :data="tableData.filter(d => !search || d.name.includes(search))">
        <el-table-column prop="type" label="Type" width="120">
          <template slot-scope="scope">
            <el-select v-model="scope.row.type" v-if="scope.row.action == 'add'" size="small" placeholder="type">
              <el-option label="A" value="A" />
              <el-option label="NS" value="NS" disabled />
              <el-option label="CNAME" value="CNAME" />
              <el-option label="PTR" value="PTR" />
              <el-option label="MX" value="MX" />
              <el-option label="TXT" value="TXT" />
              <el-option label="AAA" value="AAA" />
              <el-option label="SRV" value="SRV" />
            </el-select>
            <span v-else>{{ scope.row.type }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ttl" label="TTL" width="110">
          <template slot-scope="scope">
            <el-input type="number" size="small" v-model.number="scope.row.ttl" placeholder="ttl"
              v-if="scope.row.action" />
            <span v-else>{{ formatter(scope.row.ttl) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="Priority" width="100">
          <template slot-scope="scope">
            <el-input type="number" size="small" v-model.number="scope.row.priority" placeholder="priority"
              v-if="scope.row.action" />
            <span v-else>{{ formatter(scope.row.priority) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="Name" min-width="200">
          <template slot-scope="scope">
            <el-input size="small" v-model="scope.row.name" placeholder="name" v-if="scope.row.action == 'add'" />
            <span v-else>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="Content" min-width="350">
          <template slot-scope="scope">
            <el-input size="small" v-model="scope.row.content" placeholder="content" v-if="scope.row.action" />
            <span v-else>{{ scope.row.content }}</span>
          </template>
        </el-table-column>

        <el-table-column label width="150">
          <!-- eslint-disable-next-line vue/no-unused-vars -->
          <template slot="header" slot-scope="scope">
            <el-input v-model="search" size="mini" placeholder="输入关键字搜索" />
          </template>
          <template slot-scope="scope">
            <template v-if="scope.row.action">
              <el-button size="mini" type="success" @click="handleSubmit(scope.$index, scope.row)">
                提交
              </el-button>
              <el-button size="mini" @click="handleCancel(scope.$index, scope.row)">
                取消
              </el-button>
            </template>
            <template v-else>
              <el-button size="mini" @click="scope.row.action = 'edit'">
                编辑
              </el-button>
              <i style="width: 8px; display: inline-block" />
              <el-popconfirm title="确认删除?" @onConfirm="handleDelete(scope.$index, scope.row)">
                <el-button size="mini" slot="reference" type="danger">
                  删除
                </el-button>
              </el-popconfirm>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </el-col>
  </el-row>
</template>

<script>
import { encodeURI } from 'js-base64'
function addKey(data, key) {
  for (const i of data) {
    i.key = key + '/' + i.name
    if (Object.prototype.hasOwnProperty.call(i, 'subdomain')) {
      addKey(i.subdomain, i.key)
    }
  }
  return data
}
export default {
  data() {
    return {
      num: 2,
      treeData: [],
      tableData: [],
      defaultProps: {
        children: 'subdomain',
        label: 'name'
      },
      search: ''
    }
  },
  created() {
    this.created()
  },

  methods: {
    handleNodeClick(data) {
      console.log(data)
      const key = encodeURI(data.key)

      const url = '/api/v1/record/' + key
      this.$ajax
        .get(url)
        .then((response) => {
          this.tableData = response.data.data.map((item) => {
            item.action = ''
            return item
          })
        })
        .catch((error) => console.error('Error:', error))
    },
    handleAdd() {
      this.tableData.unshift({ action: 'add' })
      this.showAlert = true
    },
    handleDelete(index, row) {
      const key = encodeURI(row.key)

      const url = '/api/v1/record/' + key
      this.$ajax
        .delete(url)
        .then((response) => {
          this.$message.success('success')
          this.tableData.splice(index, 1)
        })

        .catch((error) => {
          this.$message.error(error.response.data.msg)
        })
    },
    handleCancel(index, row) {
      if (row.action === 'add') {
        this.tableData.splice(index, 1)
      } else {
        row.action = ''
      }
    },
    handleSubmit(index, row) {
      const data = {}
      let url = '/api/v1/record'
      let method = 'post'
      for (const k in row) {
        if (row[k] === '' || k === 'action') {
          continue
        }
        data[k] = row[k]
      }
      if (row.action === 'edit') {
        const key = encodeURI(row.key)
        url += '/' + key
        method = 'put'
      }
      this.$ajax({ method, url, data })

        .then((response) => {
          this.$message.success(response.data.msg)

          // this.$refs.table.doLayout()
          row.action = ''
        })
        .catch((error) => console.error('Error:', error))
    },
    formatter(cellValue) {
      return cellValue || '-'
    },
    change() {
      this.created()
    },
    created() {
      this.$ajax
        .get(`/api/v1/domains?deep=${this.num}`)
        .then((response) => {
          const data = response.data.data
          if (
            data.constructor === Object &&
            Object.prototype.hasOwnProperty.call(data, 'subdomain')
          ) {
            this.treeData = addKey(data.subdomain, '')
          }
        })
        .catch((error) => console.error('Error:', error))
    }
  }
}
</script>
