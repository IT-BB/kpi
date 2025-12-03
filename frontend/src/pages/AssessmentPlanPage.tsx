import React, { useState, useEffect } from 'react';
import { Table, Button, message, Space, Tag } from 'antd';
import { PlusOutlined, EyeOutlined } from '@ant-design/icons';
import { assessmentApi } from '../services/api';
import { AssessmentPlan } from '../types';

const AssessmentPlanPage: React.FC = () => {
  const [plans, setPlans] = useState<AssessmentPlan[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchPlans();
  }, []);

  const fetchPlans = async () => {
    setLoading(true);
    try {
      const response = await assessmentApi.listPlans();
      setPlans(response.data.data || []);
    } catch (error) {
      message.error('获取考核方案列表失败');
    } finally {
      setLoading(false);
    }
  };

  const getStatusTag = (status: string) => {
    const statusConfig: Record<string, { color: string; text: string }> = {
      draft: { color: 'default', text: '草稿' },
      pending: { color: 'orange', text: '待确认' },
      confirmed: { color: 'blue', text: '已确认' },
      in_progress: { color: 'processing', text: '进行中' },
      completed: { color: 'success', text: '已完成' },
      cancelled: { color: 'error', text: '已取消' },
    };
    const config = statusConfig[status] || { color: 'default', text: status };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '方案名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '考核周期',
      dataIndex: 'cycle_type',
      key: 'cycle_type',
      width: 120,
      render: (type: string) => {
        const typeMap: Record<string, string> = {
          monthly: '月度',
          quarterly: '季度',
          yearly: '年度',
        };
        return typeMap[type] || type;
      },
    },
    {
      title: '员工',
      dataIndex: ['employee', 'full_name'],
      key: 'employee',
      width: 120,
    },
    {
      title: '考核人',
      dataIndex: ['manager', 'full_name'],
      key: 'manager',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '开始日期',
      dataIndex: 'start_date',
      key: 'start_date',
      width: 120,
      render: (date: string) => date?.split('T')[0],
    },
    {
      title: '结束日期',
      dataIndex: 'end_date',
      key: 'end_date',
      width: 120,
      render: (date: string) => date?.split('T')[0],
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      render: (_: any, record: AssessmentPlan) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />}>
            查看
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h1>考核方案管理</h1>
        <Button type="primary" icon={<PlusOutlined />}>
          创建考核方案
        </Button>
      </div>

      <Table
        dataSource={plans}
        columns={columns}
        loading={loading}
        rowKey="id"
      />
    </div>
  );
};

export default AssessmentPlanPage;
