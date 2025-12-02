<template>
  <div class="login-wrapper">
    <form class="card" @submit.prevent="handleSubmit">
      <h1>企业 KPI 管理平台</h1>
      <p class="subtitle">请使用工作邮箱登录</p>
      <label>邮箱</label>
      <input v-model="form.email" type="email" placeholder="name@company.com" required />
      <label>密码</label>
      <input v-model="form.password" type="password" placeholder="••••••••" required />
      <button type="submit" :disabled="userStore.loading">{{ userStore.loading ? '登录中...' : '登录' }}</button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const error = ref('')
const form = reactive({
  email: 'admin@enterprise.local',
  password: 'ChangeMe123!'
})

const handleSubmit = async () => {
  error.value = ''
  try {
    await userStore.login(form)
    const redirect = route.query.redirect || '/dashboard'
    router.replace(redirect)
  } catch (err) {
    error.value = err?.response?.data?.message || '登录失败，请稍后再试'
  }
}
</script>

<style scoped>
.login-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(circle at top, #dbeafe, #e2e8f0);
}
.card {
  width: 360px;
  background: white;
  padding: 32px;
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.1);
  display: flex;
  flex-direction: column;
  gap: 12px;
}
button {
  background: #2563eb;
  border: none;
  border-radius: 12px;
  padding: 12px;
  color: white;
  font-size: 1rem;
}
input {
  padding: 12px;
  border-radius: 10px;
  border: 1px solid #cbd5f5;
}
.error {
  color: #dc2626;
  text-align: center;
}
.subtitle {
  color: #64748b;
  margin-top: -8px;
}
</style>
