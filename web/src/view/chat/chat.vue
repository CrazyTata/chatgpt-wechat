<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始时间"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束时间"></el-date-picker>
      </el-form-item>
        <el-form-item label="用户">
         <el-input v-model="searchInfo.user" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="客服">
         <el-input v-model="searchInfo.open_kf_id" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="应用">
             <el-input v-model.number="searchInfo.agent_id" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="客服类型">
          <el-select v-model="searchInfo.chat_type" clearable placeholder="请选择">
            <el-option
              v-for="item in chatType"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button v-auth="btnAuth.chatExport" icon="download" @click="onExport">导出</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
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
        <el-table-column align="left" label="用户" prop="user" width="120" />
        <el-table-column align="left" label="客服" prop="open_kf_id" width="120" />
        <el-table-column align="left" label="应用" prop="agent_id" width="120" />
        <el-table-column align="left" label="发送内容" prop="req_content" width="120" />
        <el-table-column align="left" label="响应内容" prop="res_content" width="480" />
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
        <el-form-item label="用户:"  prop="user" >
          <el-input v-model="formData.user" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="客服:"  prop="open_kf_id" >
          <el-input v-model="formData.open_kf_id" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="应用:"  prop="agent_id" >
          <el-input v-model.number="formData.agent_id" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="发送内容:"  prop="req_content" >
          <el-input v-model="formData.req_content" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="响应内容:"  prop="res_content" >
          <el-input v-model="formData.res_content" :clearable="true"  placeholder="请输入" />
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'Chat'
}
</script>

<script setup>
import {
  createChat,
  deleteChat,
  deleteChatByIds,
  updateChat,
  findChat,
  getChatList,
  exportChatList
} from '@/api/chat'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useBtnAuth } from '@/utils/btnAuth'
const btnAuth = useBtnAuth()
// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        user: '',
        message_id: '',
        open_kf_id: '',
        agent_id: 0,
        req_content: '',
        res_content: '',
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
const chatType = ref([
  {
    value: 1,
    label: '机器人'
  },{
    value: 2,
    label: '客服',
  }
])
// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}
// 搜索
const onExport = () => {
  exportData()
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
  const table = await getChatList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
// 导出
const exportData = async() => {
  const table = await exportChatList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    const baseUrl = ""; // 修改为实际的服务器地址
    const fileUrl = `${baseUrl}${table.data.file}`; // 在下载地址前面加上服务器地址
    const filename = fileUrl.substring(fileUrl.lastIndexOf("/") + 1); // 获取文件名
    const elink = document.createElement('a'); // 创建a标签
    elink.download = filename;
    elink.style.display = 'none';
    elink.href = fileUrl;
    document.body.appendChild(elink);
    elink.click();
    document.body.removeChild(elink); //移除a标签
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


</script>

<style>
</style>
