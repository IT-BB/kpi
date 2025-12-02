<template>
  <LayoutShell>
    <div class="grid-layout">
      <section class="panel">
        <header class="panel-head">
          <div>
            <h3>考核任务</h3>
            <small>支持过滤和流程操作</small>
          </div>
          <button class="ghost" @click="loadAssignments">刷新</button>
        </header>
        <table>
          <thead>
            <tr>
              <th>对象</th>
              <th>模板</th>
              <th>周期</th>
              <th>状态</th>
              <th>截止日</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="assignment in assignments" :key="assignment.id">
              <td>
                <strong>{{ assignment.assignee.fullName }}</strong>
                <small>{{ assignment.manager.fullName }}</small>
              </td>
              <td>{{ assignment.template.name }}</td>
              <td>{{ assignment.cycle.name }}</td>
              <td><span class="pill" :class="assignment.status">{{ statusLabel(assignment.status) }}</span></td>
              <td>{{ formatDate(assignment.dueDate) }}</td>
              <td><button class="ghost" @click="openAssignment(assignment)">查看</button></td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="panel">
        <header>
          <h3>指派 KPI</h3>
        </header>
        <form class="form" @submit.prevent="createAssignment">
          <label>考核周期</label>
          <select v-model.number="createForm.cycleId" required>
            <option value="" disabled>选择周期</option>
            <option v-for="cycle in cycles" :key="cycle.id" :value="cycle.id">
              {{ cycle.name }}
            </option>
          </select>
          <label>模板</label>
          <select v-model.number="createForm.templateId" required>
            <option value="" disabled>选择模板</option>
            <option v-for="template in templates" :key="template.id" :value="template.id">
              {{ template.name }}
            </option>
          </select>
          <label>考核人</label>
          <select v-model.number="createForm.assigneeId" required>
            <option value="" disabled>选择员工</option>
            <option v-for="employee in employees" :key="employee.id" :value="employee.id">
              {{ employee.fullName }}
            </option>
          </select>
          <label>直属主管</label>
          <select v-model.number="createForm.managerId" required>
            <option value="" disabled>选择主管</option>
            <option v-for="manager in managers" :key="manager.id" :value="manager.id">
              {{ manager.fullName }}
            </option>
          </select>
          <label>截止日期</label>
          <input type="date" v-model="createForm.dueDate" required />
          <button type="submit">创建</button>
        </form>
      </section>
    </div>

    <div v-if="selected" class="drawer">
      <div class="drawer-content">
        <header>
          <div>
            <h3>{{ selected.template.name }} · {{ selected.assignee.fullName }}</h3>
            <small>状态：{{ statusLabel(selected.status) }}</small>
          </div>
          <button class="ghost" @click="closeDrawer">关闭</button>
        </header>
        <section class="metrics">
          <article v-for="score in scoreForm" :key="score.metricId">
            <div class="metric-head">
              <div>
                <h4>{{ score.metric.name }}</h4>
                <small>{{ score.metric.description }}</small>
              </div>
              <span>{{ score.metric.weight }}%</span>
            </div>
            <div class="metric-grid">
              <label>
                自评
                <input type="number" step="0.1" v-model.number="score.selfScore" :disabled="!canSubmit" />
              </label>
              <label>
                主管评分
                <input type="number" step="0.1" v-model.number="score.managerScore" :disabled="!canReview" />
              </label>
            </div>
            <label>
              备注
              <input v-model="score.comment" />
            </label>
          </article>
        </section>
        <textarea v-model="workflowComment" rows="3" placeholder="审批意见"></textarea>
        <div class="actions">
          <button v-if="canSubmit" @click="submitScores">提交自评</button>
          <button v-if="canReview" @click="reviewScores">主管评审</button>
          <button v-if="canFinalize" @click="finalizeAssignment">归档</button>
        </div>
      </div>
    </div>
  </LayoutShell>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import LayoutShell from '../components/LayoutShell.vue'
import api from '../api/http'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const assignments = ref([])
const templates = ref([])
const cycles = ref([])
const employees = ref([])
const managers = ref([])
const selected = ref(null)
const scoreForm = ref([])
const workflowComment = ref('')
const createForm = reactive({ cycleId: '', templateId: '', assigneeId: '', managerId: '', dueDate: '' })

const loadAssignments = async () => {
  const { data } = await api.get('/kpi/assignments')
  assignments.value = data.data
}

const loadRefs = async () => {
  const [tplRes, cycleRes, userRes] = await Promise.all([
    api.get('/kpi/templates'),
    api.get('/kpi/cycles'),
    api.get('/users')
  ])
  templates.value = tplRes.data.data
  cycles.value = cycleRes.data.data
  employees.value = userRes.data.data
  managers.value = userRes.data.data.filter((u) => u.role !== 'employee')
}

const createAssignment = async () => {
  await api.post('/kpi/assignments', {
    ...createForm,
    cycleId: Number(createForm.cycleId),
    templateId: Number(createForm.templateId),
    assigneeId: Number(createForm.assigneeId),
    managerId: Number(createForm.managerId),
    dueDate: new Date(createForm.dueDate)
  })
  await loadAssignments()
  Object.assign(createForm, { cycleId: '', templateId: '', assigneeId: '', managerId: '', dueDate: '' })
}

const openAssignment = (assignment) => {
  selected.value = assignment
  workflowComment.value = ''
  scoreForm.value = assignment.scores.map((score) => ({
    metricId: score.metricId,
    metric: score.metric,
    selfScore: score.selfScore,
    managerScore: score.managerScore,
    comment: score.comment || ''
  }))
}

const closeDrawer = () => {
  selected.value = null
  workflowComment.value = ''
  scoreForm.value = []
}

const canSubmit = computed(() => {
  if (!selected.value || !userStore.profile) return false
  if (selected.value.status !== 'draft') return false
  return userStore.profile.role === 'admin' || selected.value.assigneeId === userStore.profile.id
})

const canReview = computed(() => {
  if (!selected.value || !userStore.profile) return false
  if (selected.value.status !== 'submitted') return false
  return userStore.profile.role === 'admin' || selected.value.managerId === userStore.profile.id
})

const canFinalize = computed(() => {
  if (!selected.value || !userStore.profile) return false
  return userStore.profile.role === 'admin' && selected.value.status === 'reviewed'
})

const submitScores = async () => {
  await api.post(`/kpi/assignments/${selected.value.id}/submit`, {
    scores: scoreForm.value.map((score) => ({ metricId: score.metricId, score: score.selfScore, comment: score.comment })),
    comment: workflowComment.value
  })
  await refreshSelected()
}

const reviewScores = async () => {
  await api.post(`/kpi/assignments/${selected.value.id}/review`, {
    scores: scoreForm.value.map((score) => ({ metricId: score.metricId, score: score.managerScore, comment: score.comment })),
    comment: workflowComment.value
  })
  await refreshSelected()
}

const finalizeAssignment = async () => {
  await api.post(`/kpi/assignments/${selected.value.id}/finalize`)
  await refreshSelected()
}

const refreshSelected = async () => {
  await loadAssignments()
  if (selected.value) {
    const latest = assignments.value.find((item) => item.id === selected.value.id)
    if (latest) {
      openAssignment(latest)
    } else {
      closeDrawer()
    }
  }
}

const statusLabel = (status) => ({
  draft: '草稿',
  submitted: '自评完成',
  reviewed: '主管评审',
  approved: '已归档'
}[status] || status)

const formatDate = (date) => new Date(date).toLocaleDateString()

onMounted(async () => {
  await Promise.all([loadAssignments(), loadRefs()])
})
</script>

<style scoped>
.grid-layout {
  display: grid;
  grid-template-columns: 3fr 2fr;
  gap: 24px;
}
.panel {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
}
.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
table {
  width: 100%;
  border-spacing: 0;
}
th,
td {
  padding: 10px 0;
  border-bottom: 1px solid #f1f5f9;
  text-align: left;
}
.pill {
  border-radius: 999px;
  padding: 4px 12px;
  background: #e2e8f0;
}
.drawer {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.5);
  display: flex;
  justify-content: flex-end;
}
.drawer-content {
  width: 480px;
  background: #fff;
  height: 100%;
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.metrics {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.metric-head {
  display: flex;
  justify-content: space-between;
}
.metric-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}
.metric-grid input,
textarea,
select,
input {
  width: 100%;
  padding: 8px;
  border-radius: 8px;
  border: 1px solid #d0d7ff;
}
.actions {
  display: flex;
  gap: 12px;
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
</style>
