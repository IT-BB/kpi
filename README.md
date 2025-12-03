# KPI 考核系统

面向互联网云计算企业的关键绩效指标考核系统，支持多维度指标设定、目标分解、实时追踪、自动评分、360度评估等功能。

## 技术栈

### 后端
- Golang 1.21+
- Gin (Web框架)
- GORM (ORM)
- MySQL 8.0 / SQLite
- Redis

### 前端
- React
- TypeScript
- Ant Design

### 测试
- Go testing
- rapid (属性测试库)

### 文档
- OpenAPI/Swagger (swaggo)

## 功能特性

- **组织架构管理**: 部门层级、用户角色、权限体系
- **KPI指标库**: 标准化指标、多维度分类、评分规则
- **考核方案**: 目标设定、权重分配、周期管理
- **进度追踪**: 实时更新、完成率计算、预警提醒
- **绩效评分**: 自动计算、手动调整、评分审核
- **360度评估**: 多维度评价、匿名化处理、结果汇总
- **数据分析**: 绩效报表、趋势分析、人才盘点
- **审计安全**: 操作日志、权限控制、数据加密

## 快速开始

### 后端

```bash
cd backend
go mod download
go run cmd/server/main.go
```

### 前端

```bash
cd frontend
npm install
npm start
```

### Docker

```bash
docker-compose up -d
```

## 项目结构

```
.
├── backend/           # Go 后端
│   ├── cmd/          # 主程序入口
│   ├── internal/     # 内部包
│   └── pkg/          # 公共包
├── frontend/         # React 前端
│   └── src/          # 源代码
└── docs/             # 文档
```

## API 文档

启动服务后访问: http://localhost:8080/swagger/index.html

## 许可证

MIT
