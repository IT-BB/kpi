<template>
  <LayoutShell>
    <div class="cards">
      <div class="card" v-for="card in kpiCards" :key="card.label">
        <p class="label">{{ card.label }}</p>
        <h2>{{ card.value }}</h2>
        <small>{{ card.hint }}</small>
      </div>
    </div>

    <div class="panels">
      <section class="panel">
        <header>
          <h3>流程进度</h3>
          <small>实时 KPI 状态</small>
        </header>
        <table>
          <thead>
            <tr>
              <th>阶段</th>
              <th>数量</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(value, key) in summary.statusBreakdown" :key="key">
              <td>{{ labels[key] || key }}</td>
              <td>{{ value }}</td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="panel">
        <header>
          <h3>部门表现</h3>
          <small>管理者评分</small>
        </header>
        <table>
          <thead>
            <tr>
              <th>部门</th>
              <th>自评</th>
              <th>主管评分</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="dept in summary.departments" :key="dept.department">
              <td>{{ dept.department }}</td>
              <td>{{ dept.averageSelfScore?.toFixed(1) || '-' }}</td>
              <td>{{ dept.averageManagerScore?.toFixed(1) || '-' }}</td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>
  </LayoutShell>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import LayoutShell from '../components/LayoutShell.vue'
import api from '../api/http'

const summary = ref({
  totalAssignments: 0,
  averageManagerScore: 0,
  averageSelfScore: 0,
  statusBreakdown: {},
  departments: []
})
const loading = ref(true)
const labels = {
  draft: '草稿',
  submitted: '已提交',
  reviewed: '已评审',
  approved: '已归档'
}

const kpiCards = computed(() => [
  {
    label: '考核任务',
    value: summary.value.totalAssignments,
    hint: '当前周期总数'
  },
  {
    label: '平均自评',
    value: summary.value.averageSelfScore?.toFixed(1) || '0.0',
    hint: '所有 KPI 自评平均值'
  },
  {
    label: '平均主管评分',
    value: summary.value.averageManagerScore?.toFixed(1) || '0.0',
    hint: '所有 KPI 主管评分平均值'
  }
])

const loadSummary = async () => {
  loading.value = true
  try {
    const { data } = await api.get('/kpi/reports/summary')
    summary.value = data.data
  } finally {
    loading.value = false
  }
}

onMounted(loadSummary)
</script>

<style scoped>
.cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}
.card {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.05);
}
.card .label {
  margin: 0;
  color: #64748b;
}
.card h2 {
  margin: 8px 0;
}
.panels {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 24px;
}
.panel {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
}
header h3 {
  margin: 0;
}
table {
  width: 100%;
  border-spacing: 0;
}
th,
td {
  text-align: left;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
}
</style>
