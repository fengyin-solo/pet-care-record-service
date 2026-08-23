# 宠物管理系统（petsmanagement）

纯 Go 标准库实现的宠物管理后端服务，零第三方依赖。用于管理宠物主人、宠物档案、就诊记录、疫苗接种与回访计划。

## 运行

```bash
# 启动服务（默认监听 :8080）
go run ./cmd/server

# 自定义端口与配置
PORT=9090 MAX_PAGE_SIZE=100 VACCINE_DUE_DAYS=30 go run ./cmd/server
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 8080 | 监听端口 |
| ADDR | 空（覆盖 PORT） | 完整监听地址 |
| MAX_PAGE_SIZE | 100 | 分页单页最大条数 |
| VACCINE_DUE_DAYS | 30 | 疫苗到期提醒默认天数 |
| LOG_LEVEL | info | 日志级别（debug/info/warn/error） |

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`；错误时 `code` 非 0，HTTP 状态码对应 400/404/409/500。

### 主人 Owner

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/owners | 新增主人 |
| GET | /api/owners?keyword=&page=&size= | 分页列表（keyword 匹配姓名/电话/邮箱） |
| GET | /api/owners/{id} | 查询主人 |
| PUT | /api/owners/{id} | 更新主人 |
| DELETE | /api/owners/{id} | 删除主人 |

### 宠物 Pet

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/pets | 新增宠物（校验 owner_id 存在） |
| GET | /api/pets?owner_id=&species=&health_status=&status=&keyword=&page=&size= | 分页列表 |
| GET | /api/pets/{id} | 查询宠物 |
| PUT | /api/pets/{id} | 更新宠物 |
| DELETE | /api/pets/{id} | 删除宠物 |
| POST | /api/pets/{id}/archive | 归档宠物（active→archived，不可逆） |
| POST | /api/pets/batch-archive | 批量归档，body `{"ids":[...]}` |

宠物字段枚举：
- gender：`male` / `female` / `unknown`
- health_status：`healthy` / `sick` / `recovering` / `treatment`
- status：`active` / `archived`

### 就诊记录 MedicalRecord

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/medical-records | 新增就诊记录（校验 pet_id 存在） |
| GET | /api/medical-records?pet_id=&clinic=&page=&size= | 分页列表 |
| GET | /api/medical-records/{id} | 查询就诊记录 |
| PUT | /api/medical-records/{id} | 更新就诊记录 |
| DELETE | /api/medical-records/{id} | 删除就诊记录 |

### 疫苗 Vaccine

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/vaccines | 新增疫苗记录（校验 pet_id 存在） |
| GET | /api/vaccines?pet_id=&page=&size= | 分页列表 |
| GET | /api/vaccines/due?days=30 | 查询 N 天内到期的疫苗 |
| GET | /api/vaccines/{id} | 查询疫苗记录 |
| PUT | /api/vaccines/{id} | 更新疫苗记录 |
| DELETE | /api/vaccines/{id} | 删除疫苗记录 |

### 回访 FollowUp

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/follow-ups | 新增回访（校验 pet_id 存在） |
| GET | /api/follow-ups?pet_id=&status=&page=&size= | 分页列表 |
| GET | /api/follow-ups/{id} | 查询回访 |
| PUT | /api/follow-ups/{id} | 更新回访 |
| DELETE | /api/follow-ups/{id} | 删除回访 |
| POST | /api/follow-ups/{id}/complete | 完成回访（pending→done） |
| POST | /api/follow-ups/{id}/cancel | 取消回访（pending→cancelled） |

回访字段枚举：
- method：`phone` / `visit`
- status：`pending` / `done` / `cancelled`

### 统计 Stats

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats/overview | 整体概览统计 |
| GET | /healthz | 健康检查 |

## 项目结构

```
origin/
├── cmd/server/main.go       # 入口：配置加载、依赖装配、优雅关闭
├── internal/
│   ├── app/app.go           # 依赖装配 store -> service -> handler
│   ├── config/config.go     # 环境变量配置
│   ├── model/               # 领域模型 + 校验 + 状态机
│   ├── store/               # Store 接口 + 内存实现
│   ├── service/             # 业务逻辑
│   └── handler/             # HTTP 路由 + 处理器
└── pkg/
    ├── httpx/               # 统一响应、分页、JSON 解析
    ├── idgen/               # ID 生成
    └── logger/              # 分级日志
```

## 测试

```bash
go test ./...
```
