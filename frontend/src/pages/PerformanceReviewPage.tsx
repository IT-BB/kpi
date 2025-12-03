import React, { useState, useEffect } from 'react';
import { Table, Button, message, Space, Tag } from 'antd';
import { EyeOutlined } from '@ant-design/icons';
import { performanceApi } from '../services/api';
import { PerformanceReview } from '../types';

const PerformanceReviewPage: React.FC = () => {
  const [reviews, setReviews] = useState<PerformanceReview[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchReviews();
  }, []);

  const fetchReviews = async () => {
    setLoading(true);
    try {
      const response = await performanceApi.listReviews();
      setReviews(response.data.data || []);
    } catch (error) {
      message.error('获取绩效评估列表失败');
    } finally {
      setLoading(false);
    }
  };

  const getStatusTag = (status: string) => {
    const statusConfig: Record<string, { color: string; text: string }> = {
      pending: { color: 'orange', text: '待评估' },
      completed: { color: 'success', text: '已完成' },
    };
    const config = statusConfig[status] || { color: 'default', text: status };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  const getLevelTag = (level: string) => {
    const levelConfig: Record<string, { color: string; text: string }> = {
      outstanding: { color: 'red', text: '卓越' },
      excellent: { color: 'orange', text: '优秀' },
      good: { color: 'green', text: '良好' },
      fair: { color: 'blue', text: '合格' },
      poor: { color: 'default', text: '待改进' },
    };
    const config = levelConfig[level] || { color: 'default', text: level };
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
      title: '员工',
      dataIndex: ['employee', 'full_name'],
      key: 'employee',
      width: 120,
    },
    {
      title: '评估人',
      dataIndex: ['reviewer', 'full_name'],
      key: 'reviewer',
      width: 120,
    },
    {
      title: '自评分数',
      dataIndex: 'self_score',
      key: 'self_score',
      width: 100,
      render: (score?: number) => score?.toFixed(2) || '-',
    },
    {
      title: '经理评分',
      dataIndex: 'manager_score',
      key: 'manager_score',
      width: 100,
      render: (score?: number) => score?.toFixed(2) || '-',
    },
    {
      title: '最终得分',
      dataIndex: 'final_score',
      key: 'final_score',
      width: 100,
      render: (score?: number) => score?.toFixed(2) || '-',
    },
    {
      title: '绩效等级',
      dataIndex: 'performance_level',
      key: 'performance_level',
      width: 100,
      render: (level: string) => level ? getLevelTag(level) : '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      render: (_: any, record: PerformanceReview) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />}>
            查看详情
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <h1>绩效评估管理</h1>
      </div>

      <Table
        dataSource={reviews}
        columns={columns}
        loading={loading}
        rowKey="id"
      />
    </div>
  );
};

export default PerformanceReviewPage;
