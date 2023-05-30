<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="客服ID:" prop="kfId">
          <el-input v-model="formData.kfId" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="kfName字段:" prop="kfName">
          <el-input v-model="formData.kfName" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="prompt字段:" prop="prompt">
          <el-input v-model="formData.prompt" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="发送请求的model:" prop="postModel">
          <el-input v-model="formData.postModel" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="是否启用embedding:" prop="embeddingEnable">
          <el-switch v-model="formData.embeddingEnable" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
        </el-form-item>
        <el-form-item label="embedding的搜索模式:" prop="embeddingMode">
          <el-input v-model="formData.embeddingMode" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="分数:" prop="score">
          <el-input-number v-model="formData.score" :precision="2" :clearable="true"></el-input-number>
        </el-form-item>
        <el-form-item label="topK:" prop="topK">
          <el-input v-model.number="formData.topK" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="需要清理上下文的时间，按分配置，默认0不清理:" prop="clearContextTime">
          <el-input v-model.number="formData.clearContextTime" :clearable="true" placeholder="请输入" />
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
  name: 'CustomerConfig'
}
</script>

<script setup>
import {
  createCustomerConfig,
  updateCustomerConfig,
  findCustomerConfig
} from '@/api/customerConfig'

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            kfId: '',
            kfName: '',
            prompt: '',
            postModel: '',
            embeddingEnable: false,
            embeddingMode: '',
            score: 0,
            topK: 0,
            clearContextTime: 0,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findCustomerConfig({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.recustomerConfig
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
               res = await createCustomerConfig(formData.value)
               break
             case 'update':
               res = await updateCustomerConfig(formData.value)
               break
             default:
               res = await createCustomerConfig(formData.value)
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
