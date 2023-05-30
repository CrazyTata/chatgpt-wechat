<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="用户Id:" prop="user">
          <el-input v-model="formData.user" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="消息ID:" prop="message_id">
          <el-input v-model="formData.message_id" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="客服标识:" prop="open_kf_id">
          <el-input v-model="formData.open_kf_id" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="应用ID:" prop="agent_id">
          <el-input v-model.number="formData.agent_id" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="用户发送内容:" prop="req_content">
          <el-input v-model="formData.req_content" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="openai响应内容:" prop="res_content">
          <el-input v-model="formData.res_content" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
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
  updateChat,
  findChat
} from '@/api/chat'

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
const route = useRoute()
const router = useRouter()

const type = ref('')
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

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findChat({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.rechat
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createChat(formData.value)
               break
             case 'update':
               res = await updateChat(formData.value)
               break
             default:
               res = await createChat(formData.value)
               break
           }
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
