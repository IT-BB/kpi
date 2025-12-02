<template>
  <LayoutShell>
    <div class="reports">
      <section class="panel">
        <header>
          <h3>整体表现</h3>
          <small>来自 /kpi/reports/summary</small>
        </header>
        <div class="grid">
          <div class="card">
            <p>总任务</p>
            <h2>{{ summary.totalAssignments }}</h2>
          </div>
          <div class="card">
            <p>平均自评</p>
            <h2>{{ summary.averageSelfScore.toFixed(1) }}</h2>
          </div>
          <div class="card">
            <p>平均主管评分</p>
            <h2>{{ summary.averageManagerScore.toFixed(1) }}</h2>
          </div>
        </div>
      </section>

      <section class="panel">
        <header>
          <h3>审批进度</h3>
        </header>
        <ul class="status-list">
          <li v-for="(count, key) in summary.statusBreakdown" :key="key">
            <span>{{ labels[key] || key }}</span>
            <strong>{{ count }}</strong>
          </li>
        </ul>
      </section>

      <section class="panel">
        <header>
          <h3>部门得分</h3>
        </header>
        <table>
          <thead>
            <tr>
              <th>部门</th>
              <th>人数</th>
              <th>自评</th>
              <th>主管评分</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="dept in summary.departments" :key="dept.department">
              <td>{{ dept.department }}</td>
              <td>{{ dept.headcount }}</td>
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
import { onMounted, reactive } from 'vue'
import LayoutShell from '../components/LayoutShell.vue'
import api from '../api/http'

const summary = reactive({
  totalAssignments: 0,
  averageSelfScore: 0,
  averageManagerScore: 0,
  statusBreakdown: {},
  departments: []
})
const labels = {
  draft: '草稿',
  submitted: '已提交',
  reviewed: '已评审',
  approved: '已归档'
}

const loadSummary = async () => {
  const { data } = await api.get('/kpi/reports/summary')
  Object.assign(summary, data.data)
}

onMounted(loadSummary)
</script>

<style scoped>
.reports {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 24px;
}
.panel {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
}
.card {
  background: #f8fafc;
  border-radius: 12px;
  padding: 16px;
}
.status-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.status-list li {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
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
