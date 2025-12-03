import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;

export const departmentApi = {
  list: () => apiClient.get('/departments'),
  get: (id: number) => apiClient.get(`/departments/${id}`),
  create: (data: any) => apiClient.post('/departments', data),
  update: (id: number, data: any) => apiClient.put(`/departments/${id}`, data),
  delete: (id: number) => apiClient.delete(`/departments/${id}`),
};

export const kpiIndicatorApi = {
  list: (params?: any) => apiClient.get('/kpi-indicators', { params }),
  get: (id: number) => apiClient.get(`/kpi-indicators/${id}`),
  create: (data: any) => apiClient.post('/kpi-indicators', data),
  update: (id: number, data: any) => apiClient.put(`/kpi-indicators/${id}`, data),
  delete: (id: number) => apiClient.delete(`/kpi-indicators/${id}`),
  publish: (id: number) => apiClient.post(`/kpi-indicators/${id}/publish`),
};

export const assessmentApi = {
  listPlans: (params?: any) => apiClient.get('/assessment-plans', { params }),
  getPlan: (id: number) => apiClient.get(`/assessment-plans/${id}`),
  createPlan: (data: any) => apiClient.post('/assessment-plans', data),
  updatePlan: (id: number, data: any) => apiClient.put(`/assessment-plans/${id}`, data),
  confirmPlan: (id: number) => apiClient.post(`/assessment-plans/${id}/confirm`),
  updateProgress: (data: any) => apiClient.post('/kpi-progress', data),
  getProgress: (planItemId: number) => apiClient.get(`/kpi-progress/${planItemId}`),
};

export const performanceApi = {
  listReviews: (params?: any) => apiClient.get('/performance-reviews', { params }),
  getReview: (id: number) => apiClient.get(`/performance-reviews/${id}`),
  initiateReview: (planId: number) => apiClient.post(`/performance-reviews/initiate/${planId}`),
  submitSelfReview: (id: number, data: any) => apiClient.post(`/performance-reviews/${id}/self-review`, data),
  submitManagerReview: (id: number, data: any) => apiClient.post(`/performance-reviews/${id}/manager-review`, data),
  adjustScore: (itemId: number, data: any) => apiClient.post(`/performance-reviews/items/${itemId}/adjust-score`, data),
  finalizeReview: (id: number) => apiClient.post(`/performance-reviews/${id}/finalize`),
};
