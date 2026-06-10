# 乒乓球赛事全流程管理平台 — 架构说明文档

## 1. 系统架构概览

```
┌──────────────────────────────────────────────────┐
│                    Nginx / :3000                  │
│              (Svelte SPA + 反向代理)              │
└──────────────────┬───────────────────────────────┘
                   │ /api/* → proxy
┌──────────────────▼───────────────────────────────┐
│              Go Gin / :8080                       │
│  ┌─────────┐  ┌──────────┐  ┌────────────────┐  │
│  │ Handler │→│ Service   │→│ Algorithm       │  │
│  │  Layer  │  │  Layer    │  │ (Snake/Random) │  │
│  └────┬────┘  └─────┬────┘  └────────────────┘  │
│       │             │                            │
│  ┌────▼────┐  ┌─────▼────┐                      │
│  │Middleware│  │  Model   │                      │
│  │(JWT/RBAC)│  │  Layer   │                      │
│  └─────────┘  └─────┬────┘                      │
└──────────────────┬───────────────────────────────┘
                   │
┌──────────────────▼───────────────────────────────┐
│             MongoDB / :27017                      │
│   Database: tournament                            │
└──────────────────────────────────────────────────┘
```

技术栈: SvelteKit + Go Gin + MongoDB + Docker

---

## 2. MongoDB 集合设计

| 集合 | 用途 | 关键索引 |
|------|------|----------|
| `users` | 用户（选手/裁判/组委会） | `username` (unique), `role` |
| `tournaments` | 赛事 | `status`, `created_by` |
| `events` | 比赛项目（男单/女单/混双） | `tournament_id` |
| `registrations` | 报名记录 | `(event_id, player_id)` (unique), `status` |
| `brackets` | 对阵表 | `event_id` (unique) |
| `matches` | 比赛场次 | `bracket_id`, `(event_id, round, position)`, `referee_id` |
| `audit_logs` | 审计日志 | `operator_id`, `(target_type, target_id)`, `created_at` |
| `walkover_records` | 弃权记录 | `player_id`, `event_id` |

### 集合关系图

```
tournaments ──1:N──> events ──1:N──> registrations
                   │                     │
                   └──1:1──> brackets ──1:N──> matches
                                                  │
                              walkover_records <──┘
                              audit_logs <────────┘
```

### 核心字段说明

**users**: `id, username, password(bcrypt), display_name, role(player|referee|committee), ranking`

**tournaments**: `id, name, start_date, end_date, location, status(draft→open→registration_closed→drawn→in_progress→completed→published)`

**events**: `id, tournament_id, name, type(ms|ws|xd), draw_method(snake|random), seed_count, bracket_size, status, registration_open`

**matches**: `id, bracket_id, event_id, round, position, player1_id, player2_id, score1, score2, best_of, games[], winner_id, status(pending|completed|walkover), referee_id, next_match_id, next_slot, walkover`

---

## 3. Gin 路由与中间件结构

### 中间件

| 中间件 | 功能 |
|--------|------|
| `cors` | 跨域支持（开发环境 localhost:3000/5173） |
| `AuthRequired()` | JWT Bearer Token 认证 |
| `RoleRequired(roles...)` | RBAC 角色控制 |

### 路由表

```
POST   /api/v1/register                     注册
POST   /api/v1/login                        登录
GET    /api/v1/profile                      当前用户信息 [Auth]
GET    /api/v1/users                         用户列表 [Auth, Committee]
PUT    /api/v1/users/:id/ranking             修改排名 [Auth, Committee]

POST   /api/v1/tournaments                   创建赛事 [Auth, Committee]
GET    /api/v1/tournaments                   赛事列表 [Auth]
GET    /api/v1/tournaments/:id              赛事详情 [Auth]
PUT    /api/v1/tournaments/:id/status        更新状态 [Auth, Committee]
POST   /api/v1/tournaments/:id/publish       发布成绩 [Auth, Committee]

POST   /api/v1/tournaments/:tournament_id/events        创建项目 [Auth, Committee]
GET    /api/v1/tournaments/:tournament_id/events        项目列表 [Auth]
POST   /api/v1/events/:event_id/close-registration      关闭报名 [Auth, Committee]

POST   /api/v1/events/:event_id/register                选手报名 [Auth]
GET    /api/v1/events/:event_id/registrations            报名列表 [Auth]
POST   /api/v1/events/:event_id/batch-approve            批量审核 [Auth, Committee]
PUT    /api/v1/registrations/:id/reject                   拒绝报名 [Auth, Committee]

POST   /api/v1/events/:event_id/draw                     一键抽签 [Auth, Committee]
GET    /api/v1/events/:event_id/bracket                   查看对阵 [Auth]

POST   /api/v1/matches/:match_id/score                    录入比分 [Auth, Referee/Committee]
POST   /api/v1/matches/:match_id/walkover                 弃权处理 [Auth, Committee/Referee]
PUT    /api/v1/matches/:match_id/override                  修改比分 [Auth, Committee]
GET    /api/v1/matches/:match_id                          场次详情 [Auth]
POST   /api/v1/matches/:match_id/assign-referee            指派裁判 [Auth, Committee]
GET    /api/v1/referee/matches                             裁判场次 [Auth, Referee]
GET    /api/v1/audit-logs                                  审计日志 [Auth, Committee]

GET    /api/v1/tournaments/:id/medal-board                 奖牌榜
GET    /api/v1/certificate/export                          导出PDF证书
```

### 请求流程

```
Client → CORS → AuthRequired → RoleRequired → Handler → Service → MongoDB
                                                        ↓
                                                   Algorithm (抽签)
                                                        ↓
                                                   AuditLog (修改)
```

---

## 4. Svelte 页面清单

| 路由 | 页面 | 角色权限 |
|------|------|----------|
| `/` | 首页（赛事概览） | 公开 |
| `/login` | 登录页 | 公开 |
| `/register` | 注册页（选择角色） | 公开 |
| `/tournaments` | 赛事列表 | 登录用户 |
| `/tournaments/[id]` | 赛事详情 + 项目列表 | 登录用户 |
| `/events/[event_id]/bracket` | 对阵表树形图 | 登录用户 |
| `/events/[event_id]/registrations` | 报名审核 | 组委会 |
| `/referee` | 裁判工作台（待执裁场次） | 裁判 |
| `/matches/[match_id]/score` | 比分录入页（逐局比分 + 弃权） | 裁判 |
| `/my-registrations` | 选手报名记录 | 选手 |
| `/admin` | 管理后台首页 | 组委会 |
| `/admin/medal-board` | 奖牌榜 | 组委会 |
| `/admin/certificates` | 证书导出 | 组委会 |
| `/audit` | 审计日志 | 组委会 |

### 关键组件

- **BracketTree.svelte**: 淘汰赛对阵树形图，按轮次纵向排列，显示选手名、比分、晋级路径
- **全局导航**: 根据用户角色动态显示菜单项

---

## 5. 抽签算法说明

### 5.1 蛇形分组（Snake Seeding）

按选手积分排名进行蛇形排列，确保高排位选手均匀分布。

**算法步骤**:

1. 将选手按排名从高到低排序
2. 计算对阵表大小 `bracketSize = NextPowerOf2(playerCount)`
3. 使用标准蛇形种子位置分配：
   - 1号种子 vs 末号种子（决赛半区分离）
   - 2号种子 vs 倒数第二号种子（另一半区）
   - 递归细分每个子区间
4. 排名未进入种子位的选手填充至空位

**示例**（8人签位）:

```
签位:  1    2    3    4    5    6    7    8
种子: #1   #8   #4   #5   #2   #7   #3   #6
```

这确保了:
- #1 和 #2 在决赛前不会相遇（分属上下半区）
- #1/#2 和 #3/#4 在半决赛前不会相遇

### 5.2 随机抽签（Random Draw）

1. 将所有审核通过的选手ID随机洗牌（Fisher-Yates）
2. 依次填入签位
3. 空签位设为 BYE（轮空），轮空选手自动晋级

### 5.3 BYE（轮空）处理

- 当选手数量不足以填满 2^n 签位时，空位为 BYE
- BYE 选手自动晋级到下一轮
- 系统在生成对阵时自动标记 `walkover=true` 并将对手设为胜者

### 5.4 对阵表构建

```
第一轮 (n/2场) → 第二轮 (n/4场) → ... → 半决赛 → 决赛
     ↓               ↓                    ↓
  胜者晋级       胜者晋级             胜者=冠军
```

每场比赛记录 `next_match_id` 和 `next_slot`（0=上半区/1=下半区），比赛结束后自动将胜者填入下一轮对应位置。

---

## 6. docker-compose 配置

```yaml
services:
  mongo:       # MongoDB 7, 端口 27017, 持久化数据卷
  backend:     # Go Gin, 端口 8080, 依赖 mongo
  frontend:    # Nginx + Svelte SPA, 端口 3000, 依赖 backend
```

启动方式: `docker-compose up -d`

---

## 7. 核心业务流程

### 7.1 赛事全流程

```
创建赛事 → 添加项目 → 选手报名 → 组委会审核 → 关闭报名 → 一键抽签
   → 生成对阵表 → 指派裁判 → 录入比分 → 晋级更新 → 决赛结束
   → 发布成绩 → 奖牌榜 + 证书导出
```

### 7.2 权限控制

| 操作 | 选手 | 裁判 | 组委会 |
|------|------|------|--------|
| 报名 | ✅ | - | - |
| 审核报名 | - | - | ✅ |
| 抽签 | - | - | ✅ |
| 录入比分 | - | ✅(仅自己场次) | ✅ |
| 修改比分 | - | - | ✅(留审计日志) |
| 弃权处理 | - | ✅ | ✅ |
| 查看审计日志 | - | - | ✅ |
| 发布成绩 | - | - | ✅ |

### 7.3 弃权流程

1. 裁判/组委会标记选手弃权
2. 对手自动晋级，比赛状态设为 `walkover`
3. 弃权记录写入 `walkover_records` 集合
4. 系统自动将晋级者填入下一轮对应位置
5. 递归处理连续 BYE/Walkover 场景

### 7.4 审计日志

组委会修改任意场次比分时:
1. 保存修改前的完整比赛数据（JSON）到 `old_value`
2. 执行修改操作
3. 保存修改后的数据到 `new_value`
4. 记录操作人、时间、目标类型和ID

---

## 8. 项目目录结构

```
6106/
├── docker-compose.yml
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── cmd/main.go                     # 入口
│   ├── internal/
│   │   ├── handler/                    # HTTP处理器
│   │   │   ├── user_handler.go
│   │   │   ├── tournament_handler.go
│   │   │   ├── event_handler.go
│   │   │   ├── registration_handler.go
│   │   │   ├── bracket_handler.go
│   │   │   ├── match_handler.go
│   │   │   └── ranking_handler.go
│   │   ├── middleware/auth.go          # JWT + RBAC
│   │   ├── model/models.go             # 数据模型
│   │   ├── service/                    # 业务逻辑
│   │   │   ├── user_service.go
│   │   │   ├── tournament_service.go
│   │   │   ├── event_service.go
│   │   │   ├── registration_service.go
│   │   │   ├── bracket_service.go
│   │   │   ├── match_service.go
│   │   │   └── ranking_service.go
│   │   ├── algorithm/draw.go           # 抽签算法
│   │   └── router/router.go            # 路由注册
│   └── pkg/response/response.go        # 统一响应
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── package.json
│   ├── vite.config.js
│   ├── svelte.config.js
│   └── src/
│       ├── lib/
│       │   ├── api.js                  # API客户端
│       │   ├── global.css              # 全局样式
│       │   └── components/
│       │       └── BracketTree.svelte   # 对阵树形图
│       └── routes/                     # SvelteKit页面
│           ├── +layout.svelte
│           ├── +page.svelte
│           ├── login/+page.svelte
│           ├── register/+page.svelte
│           ├── tournaments/+page.svelte
│           ├── tournaments/[id]/+page.svelte
│           ├── events/[event_id]/
│           │   ├── bracket/+page.svelte
│           │   └── registrations/+page.svelte
│           ├── matches/[match_id]/score/+page.svelte
│           ├── referee/+page.svelte
│           ├── my-registrations/+page.svelte
│           ├── admin/+page.svelte
│           ├── admin/medal-board/+page.svelte
│           ├── admin/certificates/+page.svelte
│           └── audit/+page.svelte
└── mongo/
    └── init/init.js                    # MongoDB初始化
```
