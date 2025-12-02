# 企业级 KPI 考核平台 (Go + MySQL + Vue 3)

该项目提供一个面向企业的绩效考核平台，采用 **Golang + Gin + Gorm + MySQL** 实现服务端，前端使用 **Vue 3 + Vite + Pinia + Vue Router**，覆盖 KPI 模板配置、周期管理、指派与审批流程、数据报表等核心能力。

## 功能概览

- **组织 & 人员管理**：部门维护、员工/主管/管理员三种角色、默认管理员自动初始化。
- **KPI 模板**：支持多指标、权重、目标值配置，可启停 & 版本化编辑。
- **考核周期**：自定义开始/结束时间，绑定指派。
- **任务指派**：将模板在某周期内分派给指定员工 & 主管，自动生成指标得分条目。
- **审批流程**：自评 → 主管评审 → 管理员归档，流程节点均可添加意见与审计日志。
- **报表驾驶舱**：统计不同状态任务量、均分、各部门热力图，便于 HRBP / 管理层监控。
- **前后端分离**：JWT 认证、角色鉴权、CORS、前端导航守卫，支持 SPA 体验。

## 项目结构

```
.
├── backend            # Go 服务，Gin + Gorm
│   ├── cmd/server     # 程序入口
│   ├── internal       # 业务逻辑（config/database/handler/...）
│   └── pkg            # 通用库（JWT、密码）
├── frontend           # Vue 3 + Vite 前端代码
├── docker-compose.yml # 一键启动 MySQL + API
├── .env.example       # 环境变量示例
└── README.md
```

## 后端快速启动

```bash
cd backend
cp ../.env.example ../.env   # 或自定义

# 本地开发
export $(grep -v '^#' ../.env | xargs) # 可选
go run ./cmd/server
```

- 默认端口：`8080`
- 默认管理员：`admin@enterprise.local / ChangeMe123!`
- 需要可用的 MySQL 实例（本地或 docker-compose）。

### Docker 启动

```bash
cp .env.example .env
docker compose up --build
```

Compose 将启动 `mysql:8.0` 与构建后的 API 容器（热加载/编译环境请直接使用 `go run`）。

## 前端启动

```bash
cd frontend
npm install
npm run dev # 默认 http://localhost:5173
```

- Vite dev server 自带代理，将 `/api` 请求转发到 `http://localhost:8080`。
- 生产构建：`npm run build`。

## 关键 API

| 模块 | 方法 | 路径 | 说明 | 权限 |
| ---- | ---- | ---- | ---- | ---- |
| Auth | POST | `/api/v1/auth/login` | 登录获取 JWT | 公开 |
| Auth | GET | `/api/v1/auth/profile` | 当前登录人信息 | 登录 |
| Users | CRUD | `/api/v1/users` | 员工管理 | 管理员/主管（仅查询）|
| Departments | CRUD | `/api/v1/departments` | 部门维护 | 管理员 |
| KPI 模板 | GET/POST/PUT | `/api/v1/kpi/templates` | 模板与指标 | 管理员/主管（写） |
| KPI 周期 | GET/POST | `/api/v1/kpi/cycles` | 考核周期 | 管理员 |
| KPI 指派 | GET/POST | `/api/v1/kpi/assignments` | 任务指派 | 管理员/主管 |
| KPI 提交流程 | POST | `/api/v1/kpi/assignments/:id/(submit|review|finalize)` | 自评/评审/归档 | 角色受限 |
| 报表 | GET | `/api/v1/kpi/reports/summary` | 汇总统计 | 管理员/主管 |

> 所有受保护接口均需在 `Authorization: Bearer <token>` 中携带 JWT。

## 数据模型

- `users`：包含角色、部门、上级关系。
- `departments`：组织架构。
- `kpi_templates / kpi_metrics`：模板与指标。
- `kpi_cycles`：时间窗口。
- `kpi_assignments / kpi_scores / kpi_comments`：指派、具体得分、流程意见。
- `audit_logs`：所有流程动作留痕。

## 前端亮点

- Pinia 全局状态保存 token & profile，刷新后自动恢复。
- Router 守卫限制未登录访问后台页面。
- KPI 任务抽屉提供在线自评、主管评分、管理员归档操作。
- 可视化驾驶舱展示状态柱状、部门均分、关键指标卡片。

## 下一步可扩展方向

- 绩效校准/强制分布算法。
- KPI 模板版本管理与草稿箱。
- 目标拆分（Okr-style）与评论提及。
- 数据仓库/BI 对接，导出 PDF/Excel。
- 与企业 IM/邮件系统的提醒集成。

欢迎根据业务需要继续扩展！
