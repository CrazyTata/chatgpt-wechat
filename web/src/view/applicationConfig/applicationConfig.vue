<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始时间"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束时间"></el-date-picker>
      </el-form-item>
        <el-form-item label="应用名">
         <el-input v-model="searchInfo.agentName" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="model">
         <el-input v-model="searchInfo.model" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-popover v-model:visible="deleteVisible" placement="top" width="160">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="deleteVisible = false">取消</el-button>
                <el-button type="primary" @click="onDelete">确定</el-button>
            </div>
            <template #reference>
                <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="deleteVisible = true">删除</el-button>
            </template>
            </el-popover>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="应用" prop="agentId" width="120" />
        <el-table-column align="left" label="应用secret" prop="agentSecret" width="120" />
        <el-table-column align="left" label="应用名" prop="agentName" width="120" />
        <el-table-column align="left" label="model" prop="model" width="120" />
        <el-table-column align="left" label="发送请求的model" prop="postModel" width="120" />
        <el-table-column align="left" label="openai 基础设定（可选）" prop="basePrompt" width="120" />
        <el-table-column align="left" label="进入应用时的欢迎语" prop="welcome" width="120" />
        <el-table-column align="left" label="是否启用ChatGPT应用内部交流群" prop="groupEnable" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.groupEnable) }}</template>
        </el-table-column>
        <el-table-column align="left" label="是否启用embedding" prop="embeddingEnable" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.embeddingEnable) }}</template>
        </el-table-column>
        <el-table-column align="left" label="embedding的搜索模式" prop="embeddingMode" width="120" />
        <el-table-column align="left" label="分数" prop="score" width="120" />
        <el-table-column align="left" label="topK" prop="topK" width="120" />
        <el-table-column align="left" label="需要清理上下文的时间，按分配置，默认0不清理" prop="clearContextTime" width="120" />
        <el-table-column align="left" label="ChatGPT群名" prop="groupName" width="120" />
        <el-table-column align="left" label="ChatGPT应用内部交流群chat_id" prop="groupChatId" width="120" />
        <el-table-column align="left" label="按钮组">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateApplicationConfigFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="弹窗操作">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="应用:"  prop="agentId" >
          <el-input v-model.number="formData.agentId" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="应用secret:"  prop="agentSecret" >
          <el-input v-model="formData.agentSecret" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="应用名:"  prop="agentName" >
          <el-input v-model="formData.agentName" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="model:"  prop="model" >
          <el-input v-model="formData.model" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="发送请求的model:"  prop="postModel" >
          <el-input v-model="formData.postModel" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="openai 基础设定（可选）:"  prop="basePrompt" >
          <el-input v-model="formData.basePrompt" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="进入应用时的欢迎语:"  prop="welcome" >
          <el-input v-model="formData.welcome" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="是否启用ChatGPT应用内部交流群:"  prop="groupEnable" >
          <el-switch v-model="formData.groupEnable" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
        </el-form-item>
        <el-form-item label="是否启用embedding:"  prop="embeddingEnable" >
          <el-switch v-model="formData.embeddingEnable" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
        </el-form-item>
        <el-form-item label="embedding的搜索模式:"  prop="embeddingMode" >
          <el-input v-model="formData.embeddingMode" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="分数:"  prop="score" >
          <el-input-number v-model="formData.score"  style="width:100%" :precision="2" :clearable="true"  />
        </el-form-item>
        <el-form-item label="topK:"  prop="topK" >
          <el-input v-model.number="formData.topK" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="需要清理上下文的时间，按分配置，默认0不清理:"  prop="clearContextTime" >
          <el-input v-model.number="formData.clearContextTime" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="ChatGPT群名:"  prop="groupName" >
          <el-input v-model="formData.groupName" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="ChatGPT应用内部交流群chat_id:"  prop="groupChatId" >
          <el-input v-model="formData.groupChatId" :clearable="true"  placeholder="请输入" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'ApplicationConfig'
}
</script>

<script setup>
import {
  createApplicationConfig,
  deleteApplicationConfig,
  deleteApplicationConfigByIds,
  updateApplicationConfig,
  findApplicationConfig,
  getApplicationConfigList
} from '@/api/applicationConfig'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        agentId: 0,
        agentSecret: '',
        agentName: '',
        model: '',
        postModel: '',
        basePrompt: '',
        welcome: '',
        groupEnable: false,
        embeddingEnable: false,
        embeddingMode: '',
        score: 0,
        topK: 0,
        clearContextTime: 0,
        groupName: '',
        groupChatId: '',
        })

// 验证规则
const rule = reactive({
})

const elFormRef = ref()


// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  if (searchInfo.value.groupEnable === ""){
      searchInfo.value.groupEnable=null
  }
  if (searchInfo.value.embeddingEnable === ""){
      searchInfo.value.embeddingEnable=null
  }
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getApplicationConfigList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteApplicationConfigFunc(row)
        })
    }


// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async() => {
      const ids = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          ids.push(item.ID)
        })
      const res = await deleteApplicationConfigByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        deleteVisible.value = false
        getTableData()
      }
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateApplicationConfigFunc = async(row) => {
    const res = await findApplicationConfig({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.reapplicationConfig
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteApplicationConfigFunc = async (row) => {
    const res = await deleteApplicationConfig({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        agentId: 0,
        agentSecret: '',
        agentName: '',
        model: '',
        postModel: '',
        basePrompt: '',
        welcome: '',
        groupEnable: false,
        embeddingEnable: false,
        embeddingMode: '',
        score: 0,
        topK: 0,
        clearContextTime: 0,
        groupName: '',
        groupChatId: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createApplicationConfig(formData.value)
                  break
                case 'update':
                  res = await updateApplicationConfig(formData.value)
                  break
                default:
                  res = await createApplicationConfig(formData.value)
                  break
              }
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}
</script>

<style>
</style>
