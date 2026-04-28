# Go ERP

面向 ERP 的 Go 后端骨架，采用“模块化单体 + 可演进微服务”路线，当前已落地 MVP 主链路。

## 技术栈
- Web: Gin
- ORM: GORM（MySQL / PostgreSQL）
- 配置: Viper
- 日志: Zap
- 缓存: Redis（可选）
- 消息队列: RabbitMQ（可选）
- 鉴权: JWT

## 已实现能力（MVP）
- 统一启动装配：配置、日志、DB、Redis、MQ、自动迁移
- 中间件链：`Logger -> Recovery -> JWT -> Tenant -> RBAC -> Idempotency`
- 统一响应与错误映射：`code/message/data/trace_id`
- 业务模块（分层）：用户、商品、库存、订单、支付回调
- Redis 预扣库存（防超卖基础能力）
- 订单创建事件发布（进程内事件总线 + RabbitMQ 发布器）
- 库存 Lua 原子扣减（防超卖）
- Outbox 事件派发（DB 事务与 MQ 发布最终一致）
- 订单超时取消延迟队列闭环（TTL + DLX + 消费取消）
- AI 对话接口与后台页面（后端代理 DeepSeekV4-pro）

## 快速启动
```bash
export APP_ENV=dev
go run ./cmd/server
```

## 启动后台管理页面

```bash
# 安装前端依赖
npm --prefix "./web" install

# 启动前端开发服务（默认 http://127.0.0.1:5173）
npm --prefix "./web" run dev
```

说明：
- 前端工程位于 `web/`，技术栈为 Vue3 + Vite + Element Plus。
- 路由模式为 `Vue Router Hash`（`#/...`）。
- Axios 全局开启 `withCredentials`，用于会话场景（同时保留 Bearer token 头）。
- 已配置 `/api` 代理到 `http://127.0.0.1:8080`，请先启动后端服务。

## 核心目录
```text
cmd/server/                       # 程序入口
internal/bootstrap/               # 启动装配（config/logger/db/redis/mq/router/migrate）
internal/domain/                  # 领域实体
internal/dto/                     # 请求/响应 DTO
internal/repository/              # 数据访问层
internal/service/                 # 业务服务层
internal/handler/http/            # HTTP 处理层
internal/middleware/              # 中间件
pkg/                              # 基础类库（errs/httpx/jwt/cache/lock/idempotency/mq/event/...）
configs/                          # 多环境配置
docs/                             # 架构与设计文档
```

## 示例接口
- `POST /api/v1/auth/login`
- `POST /api/v1/users`
- `GET /api/v1/users/me`
- `POST /api/v1/products/spu`
- `POST /api/v1/products/sku`
- `POST /api/v1/inventory/deduct`
- `POST /api/v1/order/create`
- `GET /api/v1/order/{id}`
- `POST /api/v1/payments/callback`

## 配置说明
- `configs/config.yaml`：默认配置（生产倾向）
- `configs/config.dev.yaml`：开发环境配置
- 关键开关：
  - `redis.enabled`：是否启用 Redis
  - `mq.enabled`：是否启用 RabbitMQ
  - `mq.order_timeout_minutes`：订单超时取消延迟分钟数
  - `mq.outbox_max_retry`：Outbox 最大重试次数
  - `mq.outbox_base_delay_seconds`：Outbox 指数退避基础秒数
  - `ai.enabled`：是否启用 AI 对话能力
  - `ai.api_key`：DeepSeek API Key（建议走环境变量覆盖）

## AI 对话配置（DeepSeekV4-pro）

在 `configs/config.dev.yaml` 中启用：

```yaml
ai:
  enabled: true
  base_url: "https://api.deepseek.com/chat/completions"
  api_key: "your_deepseek_api_key"
  model: "DeepSeekV4-pro"
  timeout_seconds: 60
  temperature: 0.7
  max_tokens: 1024
```

接口：
- `POST /api/v1/ai/chat`
