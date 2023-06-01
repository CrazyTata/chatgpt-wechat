<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="应用:" prop="agentName">
          <el-input v-model="formData.agentName" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="model:" prop="postModel">
          <el-input v-model="formData.postModel" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="prompt:" prop="basePrompt">
          <el-input v-model="formData.basePrompt" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="欢迎语:" prop="welcome">
          <el-input v-model="formData.welcome" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="启用EM:" prop="embeddingEnable">
          <el-switch v-model="formData.embeddingEnable" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
        </el-form-item>
        <el-form-item label="搜索模式:" prop="embeddingMode">
          <el-input v-model="formData.embeddingMode" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="分数:" prop="score">
          <el-input-number v-model="formData.score" :precision="2" :clearable="true"></el-input-number>
        </el-form-item>
        <el-form-item label="topK:" prop="topK">
          <el-input v-model.number="formData.topK" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="清理时间:" prop="clearContextTime">
          <el-input v-model.number="formData.clearContextTime" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="启用群:" prop="groupEnable">
          <el-switch v-model="formData.groupEnable" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
        </el-form-item>
        <el-form-item label="群名:" prop="groupName">
          <el-input v-model="formData.groupName" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="群chat_id:" prop="groupChatId">
          <el-input v-model="formData.groupChatId" :clearable="true" placeholder="请输入" />
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
  name: 'ApplicationConfig'
}
</script>

<script setup>
import {
  createApplicationConfig,
  updateApplicationConfig,
  findApplicationConfig
} from '@/api/applicationConfig'

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
const route = useRoute()
const router = useRouter()

const type = ref('')
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

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findApplicationConfig({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.reapplicationConfig
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
