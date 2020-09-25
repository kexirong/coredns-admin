<template>
  <div>
    <el-button type="primary" size="small" @click="handleAdd">添加记录</el-button>
    <el-alert v-if="showAlert" title="提示" type="custom" style="margin:8px 0;">
      <template v-for="(v,i) in description">
        <p v-html="v" :key="'desc'+i"></p>
      </template>
    </el-alert>
    <el-table :data="tableData" ref="table">
      <el-table-column prop="type" label="Type" width="120">
        <template slot-scope="scope">
          <el-select
            v-model="scope.row.type"
            v-if="scope.row.action"
            size="small"
            placeholder="type"
          >
            <el-option label="A" value="A"></el-option>
            <el-option label="NS" value="NS" disabled></el-option>
            <el-option label="CNAME" value="CNAME"></el-option>
            <el-option label="PTR" value="PTR"></el-option>
            <el-option label="MX" value="MX"></el-option>
            <el-option label="TXT" value="TXT"></el-option>
            <el-option label="AAA" value="AAA"></el-option>
            <el-option label="SRV" value="SRV"></el-option>
          </el-select>
          <span v-else>{{ scope.row.type }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="ttl" label="TTL" width="110">
        <template slot-scope="scope">
          <el-input
            type="number"
            size="small"
            v-model.number="scope.row.ttl"
            placeholder="ttl"
            v-if="scope.row.action"
          />
          <span v-else>{{ formatter(scope.row.ttl) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="priority" label="Priority" width="100">
        <template slot-scope="scope">
          <el-input
            type="number"
            size="small"
            v-model.number="scope.row.priority"
            placeholder="priority"
            v-if="scope.row.action"
          />
          <span v-else>{{ formatter(scope.row.priority) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="Name" min-width="400">
        <template slot-scope="scope">
          <el-input
            size="small"
            v-model="scope.row.name"
            placeholder="name"
            v-if="scope.row.action"
          />
          <span v-else>{{ scope.row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="content" label="Content" min-width="700">
        <template slot-scope="scope">
          <el-input
            size="small"
            v-model="scope.row.content"
            placeholder="content"
            v-if="scope.row.action"
          />
          <span v-else>{{ scope.row.content }}</span>
        </template>
      </el-table-column>

      <el-table-column label width="200">
        <template slot-scope="scope">
          <template v-if="scope.row.action">
            <el-button size="mini" type="success" @click="handleSubmit(scope.$index, scope.row)">提交</el-button>
            <el-button size="mini" @click="tableData.splice(scope.$index,1)">取消</el-button>
          </template>
          <template v-else>
            <el-button size="mini" type="danger" @click="handleDelete(scope.$index, scope.row)">删除</el-button>
          </template>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { encodeURI } from 'js-base64'

export default {
  name: 'Records',
  data () {
    return {
      showAlert: false,
      tableData: [],
      description: [
        'Name需要填写完整域名（www.baidu.com）。Priority可选填写',
        'A: &emsp;&emsp;&emsp;&emsp;用于IPv4的A记录，Name填写域名，Content填写IP',
        'NS: &emsp;&emsp;&emsp;暂不支持',
        'CNAME: Content填写别名的值',
        'PTR: &emsp;&emsp;Name填写ip, 如127.0.0.1, Content填写域名',
        'MX: &emsp;&emsp;&emsp;Content可直接填写ip也可填写域名+该域名的A记录',
        'TXT: &emsp;&emsp;Content填写txt内容',
        'AAA: &emsp;&emsp;用于IPv6的A记录',
        'SRV: &emsp;&emsp;Content格式：weight port target/hostname, 如：10 8080 127.0.0.1/srv.baidu.com'
      ]
    }
  },
  methods: {
    handleAdd () {
      this.tableData.unshift({ action: 'add' })
      this.showAlert = true
    },
    handleDelete  (index, row) {
      const key = encodeURI(row.key)

      const url = this.$base_url + '/api/v1/records/' + key
      fetch(url, { method: 'DELETE' })
        .then(async (response) => {
          if (!response.ok) {
            const json = await response.json()
            this.$message.error(json.msg)
            return
          }
          this.$message.success('success')
          this.tableData.splice(index, 1)
        })

        .catch((error) => console.error('Error:', error))
    },
    handleSubmit (index, row) {
      const data = {}
      const url = this.$base_url + '/api/v1/records'
      for (const k in row) {
        if (row[k] === '') {
          continue
        }
        data[k] = row[k]
      }
      fetch(url, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data), // data can be `string` or {object}!
        headers: new Headers({
          'Content-Type': 'application/json'
        })
      })
        .then((response) => response.json())

        .then((json) => {
          if (json.msg !== 'success') {
            this.$message.error(json.msg)
            return
          }
          this.$message.success(json.msg)

          // this.$refs.table.doLayout()
          row.action = undefined
        })
        .catch((error) => console.error('Error:', error))
    },
    formatter (cellValue) {
      return cellValue || '-'
    }
  },
  created () {
    fetch('http://localhost:8088/api/v1/records')
      .then(function (response) {
        return response.json()
      })
      .then((myJson) => {
        this.tableData = myJson.data
        console.log(myJson)
      })
  }
}
</script>
<style lang="scss" scoped>
</style>
