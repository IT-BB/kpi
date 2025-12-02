<template>
  <LayoutShell>
    <div class="split">
      <section class="panel">
        <header>
          <h3>KPI 模板</h3>
          <small>定义指标、权重与标准</small>
        </header>
        <article v-for="template in templates" :key="template.id" class="template">
          <div class="template-head">
            <div>
              <h4>{{ template.name }}</h4>
              <small>{{ template.description }}</small>
            </div>
            <button class="ghost" @click="edit(template)">编辑</button>
          </div>
          <table>
            <thead>
              <tr>
                <th>指标</th>
                <th>权重</th>
                <th>目标</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="metric in template.metrics" :key="metric.id">
                <td>{{ metric.name }}</td>
                <td>{{ metric.weight }}%</td>
                <td>{{ metric.targetValue }} {{ metric.unit }}</td>
              </tr>
            </tbody>
          </table>
        </article>
      </section>

      <section class="panel">
        <header>
          <h3>{{ editingId ? '更新模板' : '创建模板' }}</h3>
        </header>
        <form class="form" @submit.prevent="save">
          <label>名称</label>
          <input v-model="form.name" required />
          <label>说明</label>
          <textarea v-model="form.description" rows="3"></textarea>
          <label>状态</label>
          <select v-model="form.isActive">
            <option :value="true">启用</option>
            <option :value="false">停用</option>
          </select>
          <div class="metric" v-for="(metric, index) in form.metrics" :key="index">
            <div class="metric-head">
              <strong>指标 {{ index + 1 }}</strong>
              <button type="button" class="ghost" v-if="form.metrics.length > 1" @click="removeMetric(index)">删除</button>
            </div>
            <label>名称</label>
            <input v-model="metric.name" required />
            <label>描述</label>
            <input v-model="metric.description" />
            <div class="grid">
              <div>
                <label>权重 %</label>
                <input v-model.number="metric.weight" type="number" min="0" max="100" />
              </div>
              <div>
                <label>目标值</label>
                <input v-model.number="metric.targetValue" type="number" />
              </div>
              <div>
                <label>单位</label>
                <input v-model="metric.unit" />
              </div>
            </div>
          </div>
          <button type="button" class="ghost" @click="addMetric">新增指标</button>
          <div class="actions">
            <button type="submit">保存</button>
            <button type="button" class="ghost" @click="reset">重置</button>
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

const templates = ref([])
const editingId = ref(null)
const emptyMetric = () => ({ name: '', description: '', weight: 0, targetValue: 0, unit: '%' })
const form = reactive({
  name: '',
  description: '',
  isActive: true,
  metrics: [emptyMetric()]
})

const load = async () => {
  const { data } = await api.get('/kpi/templates')
  templates.value = data.data
}

const addMetric = () => {
  form.metrics.push(emptyMetric())
}

const removeMetric = (index) => {
  form.metrics.splice(index, 1)
}

const save = async () => {
  const payload = JSON.parse(JSON.stringify(form))
  if (editingId.value) {
    await api.put(`/kpi/templates/${editingId.value}`, payload)
  } else {
    await api.post('/kpi/templates', payload)
  }
  await load()
  reset()
}

const edit = (template) => {
  editingId.value = template.id
  Object.assign(form, {
    name: template.name,
    description: template.description,
    isActive: template.isActive,
    metrics: template.metrics.map((metric) => ({
      name: metric.name,
      description: metric.description,
      weight: metric.weight,
      targetValue: metric.targetValue,
      unit: metric.unit
    }))
  })
}

const reset = () => {
  editingId.value = null
  Object.assign(form, { name: '', description: '', isActive: true, metrics: [emptyMetric()] })
}

onMounted(load)
</script>

<style scoped>
.split {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}
.panel {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
}
.template {
  border-bottom: 1px solid #f1f5f9;
  padding-bottom: 16px;
  margin-bottom: 16px;
}
.template-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.metric {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 12px;
}
.metric-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}
input,
textarea,
select {
  padding: 10px;
  border-radius: 10px;
  border: 1px solid #d0d7ff;
}
button {
  border: none;
  border-radius: 10px;
  padding: 10px 16px;
  background: #2563eb;
  color: #fff;
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
