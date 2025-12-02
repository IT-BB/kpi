<template>
  <LayoutShell>
    <div class="split">
      <section class="panel">
        <header>
          <h3>员工列表</h3>
          <small>管理员与普通员工</small>
        </header>
        <table>
          <thead>
            <tr>
              <th>姓名</th>
              <th>角色</th>
              <th>部门</th>
              <th>邮箱</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>{{ user.fullName }}</td>
              <td>{{ roleLabel(user.role) }}</td>
              <td>{{ user.department?.name || '-' }}</td>
              <td>{{ user.email }}</td>
              <td>
                <button class="ghost" @click="startEdit(user)">编辑</button>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="panel">
        <header>
          <h3>{{ editingId ? '更新员工' : '创建新员工' }}</h3>
        </header>
        <form class="form" @submit.prevent="handleSave">
          <label>姓名</label>
          <input v-model="form.fullName" required />
          <label>邮箱</label>
          <input v-model="form.email" type="email" required />
          <label>职位</label>
          <input v-model="form.title" />
          <label>角色</label>
          <select v-model="form.role">
            <option value="admin">管理员</option>
            <option value="manager">主管</option>
            <option value="employee">员工</option>
          </select>
          <label>部门</label>
          <select v-model.number="form.departmentId">
            <option :value="null">未分配</option>
            <option v-for="dept in departments" :key="dept.id" :value="dept.id">
              {{ dept.name }}
            </option>
          </select>
          <label>初始密码</label>
          <input v-model="form.password" type="password" :required="!editingId" placeholder="默认密码" />
          <div class="actions">
            <button type="submit">保存</button>
            <button type="button" class="ghost" @click="resetForm">重置</button>
          </div>
        </form>
      </section>
    </div>
  </LayoutShell>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import LayoutShell from '../components/LayoutShell.vue'
import api from '../api/http'

const users = ref([])
const departments = ref([])
const editingId = ref(null)
const form = reactive({
  fullName: '',
  email: '',
  role: 'employee',
  departmentId: null,
  title: '',
  password: ''
})

const loadData = async () => {
  const [{ data: userRes }, { data: deptRes }] = await Promise.all([
    api.get('/users'),
    api.get('/departments')
  ])
  users.value = userRes.data
  departments.value = deptRes.data
}

const handleSave = async () => {
  const payload = { ...form }
  if (!payload.password) {
    delete payload.password
  }
  if (!payload.departmentId) payload.departmentId = null
  if (editingId.value) {
    await api.put(`/users/${editingId.value}`, payload)
  } else {
    await api.post('/users', payload)
  }
  await loadData()
  resetForm()
}

const startEdit = (user) => {
  editingId.value = user.id
  Object.assign(form, {
    fullName: user.fullName,
    email: user.email,
    role: user.role,
    departmentId: user.departmentId,
    title: user.title,
    password: ''
  })
}

const resetForm = () => {
  editingId.value = null
  Object.assign(form, {
    fullName: '',
    email: '',
    role: 'employee',
    departmentId: null,
    title: '',
    password: ''
  })
}

const roleLabel = (role) => ({
  admin: '管理员',
  manager: '主管',
  employee: '员工'
}[role] || role)

onMounted(loadData)
</script>

<style scoped>
.split {
  display: grid;
  grid-template-columns: 3fr 2fr;
  gap: 24px;
}
.panel {
  background: white;
  border-radius: 16px;
  padding: 20px;
}
table {
  width: 100%;
  border-spacing: 0;
}
th,
td {
  text-align: left;
  padding: 10px 0;
  border-bottom: 1px solid #f1f5f9;
}
.form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
input,
select {
  padding: 10px;
  border-radius: 10px;
  border: 1px solid #cbd5f5;
}
button {
  border: none;
  border-radius: 10px;
  padding: 10px 16px;
  background: #2563eb;
  color: white;
}
button.ghost {
  background: transparent;
  border: 1px solid #d0d7ff;
  color: #1d4ed8;
}
.actions {
  display: flex;
  gap: 12px;
}
</style>
