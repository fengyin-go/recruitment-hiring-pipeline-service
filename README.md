# 招聘系统（recruitment）

一个纯 Go 标准库（`net/http`，零第三方依赖）实现的招聘管理后端服务。采用 `cmd/internal/pkg` 标准工程分层，可编译、可测试、可运行。

## 功能特性

- 部门与职位管理：职位归属部门，记录招聘人数与已入职数，薪资区间校验。
- 候选人管理：候选人状态机 `applied → interviewing → offered → hired / rejected`，全程追踪。
- 简历管理：候选人关联简历，支持技能标签、工作经历。
- 面试管理：安排面试、完成面试（通过/不通过），自动流转候选人状态。
- Offer 管理：发放 / 接受（入职）/ 拒绝，接受后职位入职数自动 +1。
- 内推：登记内推关系，候选人入职后标记内推成功并累计奖励。
- 批量操作：批量淘汰候选人。
- 多维度统计：候选人状态分布、招聘漏斗、职位进度、面试通过率、部门统计、按月趋势、内推统计。
- 数据导出：职位 / 候选人 JSON 导出。

> 薪资字段（`salary` / `salary_min` / `salary_max` / `bonus`）单位为「元/月」（整数，避免浮点精度问题）。

## 目录结构

```
origin/
├── cmd/server/main.go
├── internal/
│   ├── app/          # 依赖装配
│   ├── config/       # 环境变量配置
│   ├── model/        # 领域模型 + 校验 + 状态机
│   ├── store/        # 数据访问接口 + 内存实现
│   ├── service/      # 业务逻辑（招聘流程/统计/导出）
│   └── handler/      # HTTP 路由 + 处理器
└── pkg/
    ├── httpx/        # 统一响应/分页/JSON
    ├── idgen/        # ID 与短码生成
    └── logger/       # 分级日志
```

## 运行

```bash
cd origin
go run ./cmd/server
```

环境变量：`PORT`（默认 8080）、`ADDR`、`MAX_PAGE_SIZE`（默认 100）、`LOG_LEVEL`（默认 info）。

## API 一览

| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/PUT/DELETE | `/api/departments` | 部门 CRUD |
| POST/GET/PUT/DELETE | `/api/positions` | 职位 CRUD |
| POST/GET | `/api/candidates` | 投递/列表 |
| GET | `/api/candidates/{id}` | 候选人详情 |
| POST | `/api/candidates/{id}/reject` | 淘汰 |
| POST | `/api/candidates/batch-reject` | 批量淘汰 |
| POST/PUT/GET | `/api/resumes` | 简历管理 |
| GET | `/api/candidates/{id}/resume` | 候选人简历 |
| POST/GET | `/api/interviews` | 安排/列表面试 |
| POST | `/api/interviews/{id}/complete` | 完成面试 |
| POST/GET | `/api/offers` | 发放/列表 Offer |
| POST | `/api/offers/{id}/accept` | 接受 Offer |
| POST | `/api/offers/{id}/decline` | 拒绝 Offer |
| POST/GET | `/api/referrals` | 内推登记/列表 |
| POST | `/api/referrals/{id}/hired` | 内推成功 |
| GET | `/api/stats/candidates` | 候选人统计 |
| GET | `/api/stats/funnel` | 招聘漏斗 |
| GET | `/api/stats/positions` | 职位进度 |
| GET | `/api/stats/interviews` | 面试通过率 |
| GET | `/api/stats/departments` | 部门统计 |
| GET | `/api/stats/monthly` | 按月趋势 |
| GET | `/api/stats/referrals` | 内推统计 |
| GET | `/api/export/positions` | 导出职位 |
| GET | `/api/export/candidates` | 导出候选人 |

## 业务闭环示例

```bash
# 1. 建部门 + 职位
curl -s -X POST localhost:8080/api/departments -d '{"name":"研发部"}'
curl -s -X POST localhost:8080/api/positions -d '{"title":"后端工程师","department_id":"<dept>","headcount":3,"salary_min":20000,"salary_max":40000}'

# 2. 候选人投递 + 简历
curl -s -X POST localhost:8080/api/candidates -d '{"name":"张三","email":"z@x.com","position_id":"<pos>"}'
curl -s -X POST localhost:8080/api/resumes -d '{"candidate_id":"<cand>","summary":"5 年后端","skills":["Go","MySQL"]}'

# 3. 安排并完成面试（通过）
curl -s -X POST localhost:8080/api/interviews -d '{"candidate_id":"<cand>","position_id":"<pos>","interviewer":"王经理","scheduled_at":"2026-08-20T10:00:00Z"}'
curl -s -X POST localhost:8080/api/interviews/<iv>/complete -d '{"passed":true,"feedback":"通过"}'

# 4. 发 Offer 并接受（入职）
curl -s -X POST localhost:8080/api/offers -d '{"candidate_id":"<cand>","position_id":"<pos>","salary":35000}'
curl -s -X POST localhost:8080/api/offers/<offer>/accept
```

## 测试

```bash
go test ./...
```
