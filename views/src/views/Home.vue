<template>
  <el-container>
    <el-header>
      <div class="container">
        <h1>CoreDNS</h1>
      </div>
    </el-header>
    <el-container style="  margin-top: 24px;">
      <el-aside width="200px">
        <el-menu>
          <el-menu-item-group>
            <template slot="title">
              <span style="font-size:20px">etcd</span>
            </template>
            <el-menu-item index="1">etcd</el-menu-item>
            <el-menu-item index="1-2">选项2</el-menu-item>
          </el-menu-item-group>
        </el-menu>
      </el-aside>
      <el-main>
        <el-button type="primary" size="small">添加解析</el-button>

        <el-table :data="tableData" style="width: 100%">
          <template slot="append">
            <el-form v-if="false" :inline="true" :model="formInline" class="demo-form-inline">
              <el-form-item label="Type">
                <el-select v-model="formInline.region" placeholder="活动区域">
                  <el-option label="区域一" value="shanghai"></el-option>
                  <el-option label="区域二" value="beijing"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="TTL">
                <el-input v-model="formInline.user" placeholder="审批人"></el-input>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="onSubmit">查询</el-button>
              </el-form-item>
            </el-form>
          </template>
          <el-table-column prop="type" label="Type" width="180"></el-table-column>
          <el-table-column prop="ttl" label="TTL" width="180"></el-table-column>
          <el-table-column prop="priority" label="Priority" width="180"></el-table-column>
          <el-table-column prop="name" label="Name" width="180"></el-table-column>
          <el-table-column prop="content" label="Content"></el-table-column>

          <el-table-column label>
            <template slot-scope="scope">
              <el-button size="mini" type="danger" @click="handleDelete(scope.$index, scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
// @ is an alias to /src

export default {
  name: 'Home',
  data () {
    return {
      formInline: {
        user: '',
        region: ''
      },
      tableData: []
    }
  },
  methods: {
    handleDelete (index, row) {
      console.log(index, row)
    },
    onSubmit () {
      console.log('submit!')
    }
  },
  beforeMount () {
    fetch('http://localhost:8088/v1/records')
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
.container {
  height: 100%;

  box-sizing: border-box;
  border-bottom: 1px solid #dcdfe6;
  h1 {
    line-height: 60px;
    margin: 0;
    float: left;
    color: #409eff;
    font-size: 30px;
    font-weight: 500;
  }
}
.el-menu {
  height: 100%;
}
</style>
