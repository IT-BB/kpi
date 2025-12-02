import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('../views/DashboardView.vue')
  },
  {
    path: '/users',
    name: 'users',
    component: () => import('../views/UsersView.vue')
  },
  {
    path: '/departments',
    name: 'departments',
    component: () => import('../views/DepartmentsView.vue')
  },
  {
    path: '/kpi/templates',
    name: 'kpi-templates',
    component: () => import('../views/KPITemplatesView.vue')
  },
  {
    path: '/kpi/assignments',
    name: 'kpi-assignments',
    component: () => import('../views/AssignmentsView.vue')
  },
  {
    path: '/reports',
    name: 'reports',
    component: () => import('../views/ReportsView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const store = useUserStore()
  if (to.meta.public) {
    if (to.name === 'login' && store.isAuthenticated) {
      return next({ name: 'dashboard' })
    }
    return next()
  }

  if (!store.isAuthenticated) {
    return next({ name: 'login', query: { redirect: to.fullPath } })
  }
  return next()
})

export default router
