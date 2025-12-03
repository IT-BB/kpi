# KPI 系统 API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: Bearer Token (JWT)

## API 端点

### 1. 部门管理 (Departments)

#### 创建部门
```
POST /departments
Content-Type: application/json

{
  "name": "研发部",
  "description": "负责产品研发",
  "parent_id": 1
}
```

#### 获取部门列表
```
GET /departments
```

#### 获取部门详情
```
GET /departments/:id
```

#### 更新部门
```
PUT /departments/:id
Content-Type: application/json

{
  "name": "研发一部",
  "description": "更新后的描述"
}
```

#### 删除部门
```
DELETE /departments/:id
```

### 2. KPI 指标管理 (KPI Indicators)

#### 创建指标
```
POST /kpi-indicators
Content-Type: application/json

{
  "name": "代码质量",
  "description": "衡量代码质量的指标",
  "category": "technical",
  "data_type": "percentage",
  "unit": "%",
  "calculation_type": "linear",
  "formula": "bug_count / total_lines * 100"
}
```

#### 获取指标列表
```
GET /kpi-indicators?category=technical&is_published=true
```

#### 发布指标
```
POST /kpi-indicators/:id/publish
```

### 3. 考核方案 (Assessment Plans)

#### 创建考核方案
```
POST /assessment-plans
Content-Type: application/json

{
  "name": "2024 Q1 考核",
  "description": "第一季度考核方案",
  "cycle_type": "quarterly",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-03-31T23:59:59Z",
  "employee_id": 10,
  "manager_id": 5,
  "items": [
    {
      "indicator_id": 1,
      "weight": 40,
      "target_value": 95,
      "challenge_value": 98,
      "baseline_value": 90
    },
    {
      "indicator_id": 2,
      "weight": 60,
      "target_value": 100,
      "challenge_value": 120,
      "baseline_value": 80
    }
  ]
}
```

#### 确认考核方案
```
POST /assessment-plans/:id/confirm
```

#### 更新 KPI 进度
```
POST /kpi-progress
Content-Type: application/json

{
  "plan_item_id": 1,
  "current_value": 92,
  "notes": "本周进展顺利"
}
```

### 4. 绩效评估 (Performance Reviews)

#### 发起评估
```
POST /performance-reviews/initiate/:planId
```

#### 提交自评
```
POST /performance-reviews/:id/self-review
Content-Type: application/json

{
  "score": 85,
  "comment": "本季度完成了主要目标"
}
```

#### 提交经理评估
```
POST /performance-reviews/:id/manager-review
Content-Type: application/json

{
  "score": 88,
  "comment": "表现优秀，超额完成目标"
}
```

#### 调整分数
```
POST /performance-reviews/items/:itemId/adjust-score
Content-Type: application/json

{
  "adjusted_score": 90,
  "reason": "考虑到特殊情况，给予额外加分"
}
```

#### 完成评估
```
POST /performance-reviews/:id/finalize
```

## 响应格式

### 成功响应
```json
{
  "data": {
    "id": 1,
    "name": "...",
    ...
  }
}
```

### 错误响应
```json
{
  "error": "错误描述信息"
}
```

## 状态码

- `200 OK` - 请求成功
- `201 Created` - 创建成功
- `400 Bad Request` - 请求参数错误
- `401 Unauthorized` - 未授权
- `403 Forbidden` - 禁止访问
- `404 Not Found` - 资源不存在
- `500 Internal Server Error` - 服务器错误
