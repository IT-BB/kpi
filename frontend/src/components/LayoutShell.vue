<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="logo">KPI Control Tower</div>
      <nav>
        <RouterLink v-for="item in menu" :key="item.path" :to="item.path" class="nav-link">
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
    </aside>
    <main class="main">
      <header class="top-bar">
        <div class="greeting">
          <p class="title">{{ title }}</p>
          <small>欢迎回来，{{ userStore.fullName }}</small>
        </div>
        <div class="actions">
          <span class="role-tag">{{ userStore.profile?.role }}</span>
          <button class="ghost" @click="handleLogout">退出</button>
        </div>
      </header>
      <section class="content">
        <slot />
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

onMounted(() => {
  if (!userStore.profile && userStore.token) {
    userStore.fetchProfile()
  }
})

const menu = [
  { path: '/dashboard', label: '驾驶舱' },
  { path: '/users', label: '员工管理' },
  { path: '/departments', label: '组织架构' },
  { path: '/kpi/templates', label: 'KPI 模板' },
  { path: '/kpi/assignments', label: '考核任务' },
  { path: '/reports', label: '分析报表' }
]

const title = computed(() => {
  const current = menu.find((item) => item.path === route.path)
  return current?.label || '概览'
})

const handleLogout = () => {
  userStore.logout()
  router.replace({ name: 'login' })
}
</script>

<style scoped>
.layout {
  display: flex;
  min-height: 100vh;
}
.sidebar {
  width: 240px;
  background: #0f172a;
  color: #e2e8f0;
  padding: 24px;
  display: flex;
  flex-direction: column;
}
.logo {
  font-size: 1.2rem;
  font-weight: 600;
  margin-bottom: 32px;
}
.nav-link {
  display: block;
  padding: 12px 16px;
  border-radius: 8px;
  color: inherit;
  margin-bottom: 8px;
}
.nav-link.router-link-active {
  background: #1d4ed8;
}
.main {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.top-bar {
  display: flex;
  justify-content: space-between;
  padding: 24px;
  background: #fff;
  border-bottom: 1px solid #e2e8f0;
}
.content {
  padding: 24px;
}
.ghost {
  background: transparent;
  border: 1px solid #cbd5f5;
  padding: 8px 16px;
  border-radius: 8px;
}
.role-tag {
  margin-right: 12px;
  text-transform: uppercase;
  font-size: 0.8rem;
  color: #475569;
}
.title {
  margin: 0;
  font-size: 1.2rem;
}
</style>
