import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, message, Space, Tag, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, CheckOutlined } from '@ant-design/icons';
import { kpiIndicatorApi } from '../services/api';
import { KPIIndicator } from '../types';

const KPIIndicatorPage: React.FC = () => {
  const [indicators, setIndicators] = useState<KPIIndicator[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingIndicator, setEditingIndicator] = useState<KPIIndicator | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchIndicators();
  }, []);

  const fetchIndicators = async () => {
    setLoading(true);
    try {
      const response = await kpiIndicatorApi.list();
      setIndicators(response.data.data || []);
    } catch (error) {
      message.error('获取指标列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingIndicator(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (indicator: KPIIndicator) => {
    setEditingIndicator(indicator);
    form.setFieldsValue(indicator);
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await kpiIndicatorApi.delete(id);
      message.success('删除成功');
      fetchIndicators();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handlePublish = async (id: number) => {
    try {
      await kpiIndicatorApi.publish(id);
      message.success('发布成功');
      fetchIndicators();
    } catch (error: any) {
      message.error(error.response?.data?.error || '发布失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingIndicator) {
        await kpiIndicatorApi.update(editingIndicator.id, values);
        message.success('更新成功');
      } else {
        await kpiIndicatorApi.create(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchIndicators();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '指标名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      width: 120,
    },
    {
      title: '数据类型',
      dataIndex: 'data_type',
      key: 'data_type',
      width: 120,
    },
    {
      title: '计算方式',
      dataIndex: 'calculation_type',
      key: 'calculation_type',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'is_published',
      key: 'is_published',
      width: 100,
      render: (published: boolean) => (
        <Tag color={published ? 'green' : 'default'}>
          {published ? '已发布' : '草稿'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: KPIIndicator) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          {!record.is_published && (
            <Button
              type="link"
              icon={<CheckOutlined />}
              onClick={() => handlePublish(record.id)}
            >
              发布
            </Button>
          )}
          <Popconfirm
            title="确认删除？"
            onConfirm={() => handleDelete(record.id)}
            okText="确认"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h1>KPI指标库管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新建指标
        </Button>
      </div>

      <Table
        dataSource={indicators}
        columns={columns}
        loading={loading}
        rowKey="id"
      />

      <Modal
        title={editingIndicator ? '编辑指标' : '新建指标'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="指标名称"
            rules={[{ required: true, message: '请输入指标名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item
            name="category"
            label="分类"
            rules={[{ required: true, message: '请选择分类' }]}
          >
            <Select>
              <Select.Option value="business">业务指标</Select.Option>
              <Select.Option value="technical">技术指标</Select.Option>
              <Select.Option value="management">管理指标</Select.Option>
              <Select.Option value="quality">质量指标</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="data_type"
            label="数据类型"
            rules={[{ required: true, message: '请选择数据类型' }]}
          >
            <Select>
              <Select.Option value="numeric">数值型</Select.Option>
              <Select.Option value="percentage">百分比</Select.Option>
              <Select.Option value="boolean">布尔型</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="unit" label="单位">
            <Input />
          </Form.Item>
          <Form.Item
            name="calculation_type"
            label="计算方式"
            rules={[{ required: true, message: '请选择计算方式' }]}
          >
            <Select>
              <Select.Option value="linear">线性计算</Select.Option>
              <Select.Option value="stepped">阶梯计算</Select.Option>
              <Select.Option value="custom">自定义公式</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="formula" label="计算公式">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default KPIIndicatorPage;
