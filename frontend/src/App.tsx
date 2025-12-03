import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider, Layout } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import DashboardPage from './pages/DashboardPage';
import DepartmentPage from './pages/DepartmentPage';
import KPIIndicatorPage from './pages/KPIIndicatorPage';
import AssessmentPlanPage from './pages/AssessmentPlanPage';
import PerformanceReviewPage from './pages/PerformanceReviewPage';
import AppLayout from './components/Layout';

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<AppLayout />}>
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<DashboardPage />} />
            <Route path="departments" element={<DepartmentPage />} />
            <Route path="kpi-indicators" element={<KPIIndicatorPage />} />
            <Route path="assessment-plans" element={<AssessmentPlanPage />} />
            <Route path="performance-reviews" element={<PerformanceReviewPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  );
};

export default App;
