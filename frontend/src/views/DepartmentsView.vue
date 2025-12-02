<template>
  <LayoutShell>
    <div class="split">
      <section class="panel">
        <header>
          <h3>部门列表</h3>
        </header>
        <table>
          <thead>
            <tr>
              <th>名称</th>
              <th>负责人</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="dept in departments" :key="dept.id">
              <td>{{ dept.name }}</td>
              <td>{{ dept.manager?.fullName || '-' }}</td>
              <td>
                <button class="ghost" @click="edit(dept)">编辑</button>
              </td>
            </tr>
          </tbody>
        </table>
      </section>
      <section class="panel">
        <header>
          <h3>{{ editingId ? '更新部门' : '创建部门' }}</h3>
        </header>
        <form class="form" @submit.prevent="save">
          <label>名称</label>
          <input v-model="form.name" required />
          <label>描述</label>
          <textarea v-model="form.description" rows="3"></textarea>
          <label>负责人</label>
          <select v-model.number="form.managerId">
            <option :value="null">暂未指定</option>
            <option v-for="manager in managers" :key="manager.id" :value="manager.id">
              {{ manager.fullName }}
            </option>
          </select>
          <div class="actions">
            <button type="submit">保存</button>
            <button type="button" class="ghost" @click="reset">取消</button>
          </div>
        </form>
      </section>
    </div>
  </LayoutShell>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import LayoutShell from '../components/LayoutShell.vue'
import api from '../api/http'

const departments = ref([])
const managers = ref([])
const editingId = ref(null)
const form = reactive({
  name: '',
  description: '',
  managerId: null
})

const load = async () => {
  const [deptRes, userRes] = await Promise.all([
    api.get('/departments'),
    api.get('/users', { params: { role: 'manager' } })
  ])
  departments.value = deptRes.data.data
  managers.value = userRes.data.data
}

const save = async () => {
  const payload = { ...form }
  if (!payload.managerId) payload.managerId = null
  if (editingId.value) {
    await api.put(`/departments/${editingId.value}`, payload)
  } else {
    await api.post('/departments', payload)
  }
  await load()
  reset()
}

const edit = (dept) => {
  editingId.value = dept.id
  Object.assign(form, {
    name: dept.name,
    description: dept.description,
    managerId: dept.managerId || null
  })
}

const reset = () => {
  editingId.value = null
  Object.assign(form, { name: '', description: '', managerId: null })
}

onMounted(load)
</script>

<style scoped>
.split {
  display: grid;
  grid-template-columns: 3fr 2fr;
  gap: 24px;
}
.panel {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
}
.form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
input,
textarea,
select {
  padding: 10px;
  border-radius: 10px;
  border: 1px solid #d0d7ff;
}
textarea {
  resize: vertical;
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
  gap: 10px;
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
</style>
