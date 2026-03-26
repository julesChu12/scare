# sCare 社区养老信息分发平台 - 系统当前状态文档

> **版本**: v1.0.0  
> **状态**: 开发中  
> **最后更新**: 2026-03-24  
> **文档类型**: 系统完整状态快照

---

## 📋 目录

1. [项目概述](#项目概述)
2. [技术架构](#技术架构)
3. [后端架构详解](#后端架构详解)
4. [前端架构详解](#前端架构详解)
5. [核心功能实现](#核心功能实现)
6. [数据库设计](#数据库设计)
7. [安全机制](#安全机制)
8. [部署架构](#部署架构)
9. [开发规范](#开发规范)
10. [已知问题与改进建议](#已知问题与改进建议)

---

## 项目概述

### 基本信息

- **项目名称**: 昌平区霍营街道社区养老信息分发平台
- **项目类型**: 社区养老服务信息分发系统（毕业设计项目）
- **Go 模块名**: `community-elderly-care-platform`
- **架构模式**: B/S 架构 + 地理围栏智能匹配

### 核心特性

1. 🎯 **地理围栏智能匹配**
   - 基于射线法（Ray Casting）的点在多边形算法
   - BoundingBox 预过滤优化性能
   - 多围栏优先级排序 + 兜底规则（最近站点）

2. 🔐 **五表细粒度权限管理**
   - `permissions` - 权限定义（API 级别）
   - `roles` - 角色定义
   - `role_permissions` - 角色-权限关联
   - `user_identities` - 用户-身份关联（支持多身份）
   - `menus` - 前端菜单配置

3. 📱 **双端隔离系统**
   - B 端：管理门户（工作人员 + 管理后台）
   - C 端：用户端（PWA，支持离线）
   - JWT Token 端类型隔离（`type: "b_end" | "c_end"`）
   - 独立的身份体系与权限模型

4. 🚀 **PWA 支持**
   - C 端支持安装到桌面
   - 离线访问能力
   - NetworkFirst 缓存策略（API 请求缓存 1 小时）

5. 👴 **适老化设计**
   - 三档字体大小调节（14px/16px/18px）
   - CSS 变量系统驱动的全局切换
   - 持久化用户偏好设置

---

## 技术架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         客户端层                              │
├─────────────────────────────┬───────────────────────────────┤
│       B 端管理门户            │         C 端用户端（PWA）       │
│   Vue 3 + TypeScript        │      Vue 3 + TypeScript       │
│   Element Plus + ECharts    │      Element Plus             │
│   Playwright (E2E 测试)     │      Vite PWA                 │
└──────────────┬──────────────┴──────────────┬────────────────┘
               │                              │
               ▼                              ▼
┌───────────────────────────────────────────────────────────────┐
│                      Nginx 反向代理                            │
│              (静态资源 + API 请求转发)                         │
└───────────────────────────────┬───────────────────────────────┘
                                │
                                ▼
┌───────────────────────────────────────────────────────────────┐
│                        Go 后端服务                            │
│                     Gin Web Framework                        │
├───────────────────────────────────────────────────────────────┤
│  中间件链：                                                    │
│  AuthMiddleware → RequireEndType → PermissionMiddleware       │
│                                                              │
│  分层架构：                                                    │
│  Handler → Service → Repository → Model (GORM Gen)            │
└───────────────────────────────┬───────────────────────────────┘
                                │
                ┌───────────────┼───────────────┐
                ▼               ▼               ▼
         ┌──────────┐    ┌──────────┐    ┌──────────┐
         │  MySQL   │    │  Redis   │    │ 阿里云OSS │
         │  8.0+    │    │  7.0+    │    │  (可选)  │
         └──────────┘    └──────────┘    └──────────┘
```

### 技术栈总览

#### 后端技术栈

| 类别 | 技术 | 版本 | 说明 |
|-----|------|------|------|
| **语言** | Go | 1.25 | - |
| **Web 框架** | Gin | 1.11 | HTTP 路由与中间件 |
| **ORM** | GORM | 1.31 | 数据库操作 |
| **代码生成** | GORM Gen | 0.3.27 | 模型与查询自动生成 |
| **数据库驱动** | MySQL Driver | 1.6 | 生产环境 |
| | SQLite Driver | 1.6 | 开发/测试环境 |
| **缓存** | Redis | 9.5 (go-redis/v9) | 会话、黑名单、缓存 |
| **认证** | JWT | 5.2 (golang-jwt/jwt/v5) | 双端 Token 认证 |
| **日志** | Zap | 1.27 | 结构化日志 |
| | Lumberjack | - | 日志轮转 |
| **CLI** | Cobra | 1.10 | 命令行工具 |
| | Viper | 1.18 | 配置管理 |
| **API 文档** | Swagger | 1.16 (swaggo/swag) | 自动生成 API 文档 |
| **测试** | Testify | 1.11 | 单元测试 |
| **其他** | Miniredis | 2.36.1 | Redis 模拟（测试） |
| | UUID | 1.6 | 唯一标识生成 |
| | Excelize | 2.10.0 | Excel 报表生成 |
| | Crypto | 0.47.0 | 加密工具 |

#### 前端技术栈

**B 端管理门户** (`frontend/management-portal/`):

| 类别 | 技术 | 版本 |
|-----|------|------|
| **框架** | Vue | 3.4.15 |
| **语言** | TypeScript | 5.3.3 |
| **构建工具** | Vite | 6.0 |
| **UI 组件库** | Element Plus | 2.5.2 |
| **状态管理** | Pinia | 2.1.7 |
| **路由** | Vue Router | 4.2.5 |
| **图表** | ECharts | 5.4.3 |
| | vue-echarts | 6.6.8 |
| **HTTP 客户端** | Axios | 1.6.5 |
| **富文本编辑器** | @wangeditor/editor | 5.1.23 |
| | @wangeditor/editor-for-vue | 5.1.12 |
| **地图** | 高德地图 (@amap/amap-jsapi-loader) | 1.0.1 |
| **图标** | @element-plus/icons-vue | 2.3.1 |
| **日期处理** | dayjs | 1.11.19 |
| **E2E 测试** | Playwright | 1.58.1 |
| **代码质量** | ESLint | 8.56.0 |
| | Prettier | 3.2.4 |

**C 端用户端** (`frontend/c-end/`):

| 类别 | 技术 | 版本 |
|-----|------|------|
| **框架** | Vue | 3.4.0 |
| **语言** | TypeScript | 5.4 |
| **构建工具** | Vite | 5.0.12 |
| **UI 组件库** | Element Plus | 2.5.4 |
| **状态管理** | Pinia | 2.1.7 |
| **路由** | Vue Router | 4.2.5 |
| **HTTP 客户端** | Axios | 1.6.5 |
| **工具库** | @vueuse/core | 10.7.2 |
| **图标** | @element-plus/icons-vue | 2.3.1 |
| **PWA** | vite-plugin-pwa | 0.17.5 |
| **开发工具** | @vitejs/plugin-basic-ssl | 1.2.0 (HTTPS) |
| | @vitejs/plugin-vue | 5.0.3 |
| | canvas | 3.2.1 (图标生成) |
| | png-to-ico | 3.0.1 (图标生成) |

---

## 后端架构详解

### 分层架构

```
┌──────────────────────────────────────────────────────────┐
│  Handler 层 (HTTP 处理器)                                 │
│  - 参数校验 (Params)                                      │
│  - 响应格式化 (Response)                                  │
│  - 错误转换 (Error → JSON)                                │
└────────────────┬─────────────────────────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────────────────────────┐
│  Service 层 (业务逻辑)                                     │
│  - 业务规则实现                                           │
│  - 事务处理                                               │
│  - 跨表操作协调                                           │
│  - 外部服务调用 (邮件、存储等)                             │
└────────────────┬─────────────────────────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────────────────────────┐
│  Repository 层 (数据访问)                                  │
│  - GORM 查询封装                                          │
│  - 复杂查询构建                                           │
│  - 数据库交互                                             │
└────────────────┬─────────────────────────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────────────────────────┐
│  Model 层 (数据模型)                                       │
│  - GORM Gen 自动生成 (禁止手动修改)                        │
│  - 数据库表结构映射                                       │
└──────────────────────────────────────────────────────────┘
```

### 依赖注入容器

**位置**: `internal/router/deps.go`

**初始化顺序**（严格遵守依赖关系）:

```go
type Deps struct {
    // 1. 基础设施
    DB           *gorm.DB
    Redis        *redis.Client
    Config       *config.Config
    JWTManager   *jwt.Manager
    
    // 2. Repository 层
    UserRepo             *repository.UserRepo
    UserIdentityRepo     *repository.UserIdentityRepo
    ZoneRepo             *repository.ZoneRepo
    StationRepo          *repository.StationRepo
    RequestRepo          *repository.RequestRepo
    TaskRepo             *repository.TaskRepo
    RoleRepo             *repository.RoleRepo
    PermissionRepo       *repository.PermissionRepo
    RolePermissionRepo   *repository.RolePermissionRepo
    MenuRepo             *repository.MenuRepo
    NotificationRepo     *repository.NotificationRepo
    
    // 3. Service 层
    AuthService          *service.AuthService
    GeofenceService      *service.GeofenceService
    RequestService       *service.RequestService
    TaskService          *service.TaskService
    PermissionService    *service.PermissionService
    MenuService          *service.MenuService
    NotificationService  *service.NotificationService
    StorageService       *service.StorageService
    
    // 4. Handler 层
    AuthHandler          *handler.AuthHandler
    RequestHandler       *handler.RequestHandler
    TaskHandler          *handler.TaskHandler
    UploadHandler        *handler.UploadHandler
}
```

### 路由与中间件设计

#### 中间件链

```
公开路由 (无中间件)
├─ /api/v1/b/auth/login
├─ /api/v1/c/auth/*
└─ /swagger/*

受保护路由 (完整中间件链)
├─ /api/v1/b/*
│   └─ AuthMiddleware → RequireEndType("b_end") → PermissionMiddleware
└─ /api/v1/c/*
    └─ AuthMiddleware → RequireEndType("c_end")
```

#### AuthMiddleware 实现

**位置**: `internal/middleware/auth.go`

**功能**:
1. 解析 JWT Token → Claims
2. 检查 Token 黑名单（单 Token）
3. 检查用户级撤销（角色变更）
4. 将用户信息存入 `gin.Context`

```go
// 设置到 Context 的字段
c.Set("user_id", claims.UserID)
c.Set("user_identities", claims.Identities)
c.Set("user_primary", claims.Primary)
c.Set("user_type", claims.Type)          // "b_end" | "c_end"
c.Set("station_id", claims.StationID)
c.Set("token_id", claims.ID)
c.Set("token_expires_at", claims.ExpiresAt)
```

#### RequireEndType 实现

**位置**: `internal/middleware/end_type.go`

**功能**: 防止跨端访问

```go
func RequireEndType(expectedType string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userType := c.GetString("user_type")
        if userType != expectedType {
            c.AbortWithStatusJSON(403, gin.H{
                "msg": "token end type mismatch",
                "data": gin.H{"expected": expectedType, "got": userType},
            })
            return
        }
        c.Next()
    }
}
```

#### PermissionMiddleware 实现

**位置**: `internal/middleware/permission.go`

**功能**: RBAC 权限检查

**检查流程**:
1. C 端用户跳过检查
2. Admin 身份跳过所有检查
3. 公共 API（`is_public=1`）跳过检查
4. 调用 `CheckAPIPermission(identities, path, method)`

---

## 前端架构详解

### B 端管理门户架构

#### 目录结构

```
src/
├── pages/              # 页面组件（按功能模块组织）
│   ├── dashboard/     # 仪表盘
│   ├── tasks/         # 任务管理
│   ├── requests/      # 需求管理
│   ├── users/         # 用户管理
│   └── ...
├── layouts/           # 布局组件
│   └── BasicLayout.vue # 侧边栏 + 顶栏 + 内容区
├── components/        # 共享组件
│   ├── MapEditor/     # 地理围栏编辑器
│   └── ImageUpload/   # 图片上传组件
├── router/            # 路由配置
│   └── guards/        # 路由守卫
│       └── permission.guard.ts
├── store/             # Pinia 状态管理
│   └── modules/
│       └── auth.ts    # 认证状态
├── api/               # API 层（集中管理）
│   └── index.ts       # 所有接口单文件定义
├── directives/        # 自定义指令
│   └── permission.ts  # v-permission 权限指令
├── composables/       # 组合式函数
│   └── usePermission.ts
├── types/             # TypeScript 类型定义
│   └── api.ts         # API 相关类型
├── utils/             # 工具函数
│   └── request.ts     # HTTP 客户端
└── config/            # 配置文件
    └── index.ts
```

#### 权限体系（四层防护）

1. **路由守卫** (`permission.guard.ts`)
   - 检查 `meta.permission_code`
   - 无权限 → 403 或跳转首页

2. **v-permission 指令** (`directives/permission.ts`)
   - DOM 级别权限控制
   - 无权限 → 移除元素

3. **usePermission composable**
   ```typescript
   const { hasPermission, hasAnyPermission } = usePermission()
   ```

4. **动态菜单** (`/b/menus/user`)
   - 从后端获取用户可见菜单
   - 根据 `permission_code` 过滤

#### HTTP 客户端

**位置**: `src/utils/request.ts`

**特性**:
- 自动附加 JWT Token
- 401 自动登出并跳转登录
- 详细的错误分类处理（401/403/404/409/500）
- 不自动解包响应（保留 `{ msg, data }` 格式）

### C 端用户端架构

#### 目录结构

```
src/
├── views/             # 页面组件
│   ├── Login.vue      # 登录页
│   ├── Quick.vue      # 快速开通
│   ├── Home.vue       # 首页
│   └── ...
├── api/               # API 层（模块化拆分）
│   ├── client.ts      # HTTP 客户端
│   ├── auth.ts        # 认证 API
│   ├── requests.ts    # 需求 API
│   ├── profile.ts     # 资料 API
│   ├── news.ts        # 新闻 API
│   └── index.ts       # 统一导出
├── store/             # Pinia 状态管理
│   ├── tokenStore.ts  # Token 状态
│   └── userStore.ts   # 用户状态
├── router/            # 路由配置
│   └── index.ts
├── composables/       # 组合式函数
│   └── useFontSize.ts # 适老化字体调节
├── styles/            # 样式文件
│   └── variables.css  # CSS 变量系统
└── config/            # 配置文件
    └── serviceTypes.ts # 服务类型
```

#### PWA 配置

**位置**: `vite.config.ts`

**Workbox 缓存策略**:

| 资源类型 | 策略 | 过期时间 |
|---------|------|---------|
| API 请求 (`/api/v1/*`) | NetworkFirst | 1 小时 |
| 图片 | CacheFirst | 30 天 |
| 字体 | CacheFirst | 365 天 |

**离线能力**:
- `navigateFallback: '/offline.html'` - 离线回退页面
- `cleanupOutdatedCaches: true` - 自动清理旧缓存

#### API 层架构

**HTTP 客户端特性**:

1. **队列化 Token 刷新**
   - 并发 401 请求时只刷新一次
   - 其他请求排队等待新 Token

2. **自动解包响应**
   ```typescript
   // 后端格式 { msg, data } → 自动提取 data
   if (body && typeof body === 'object' && 'data' in body) {
       return body.data
   }
   ```

3. **错误处理**
   - 网络错误 → "网络连接失败"
   - 401 → 刷新 Token 或跳转登录
   - 其他 → 显示后端 `msg`

#### 适老化特性

**三档字体大小**:

```typescript
type FontSize = 'small' | 'medium' | 'large'

const fontSizeConfig = {
  small: { label: '小', value: '14px' },
  medium: { label: '中', value: '16px' },
  large: { label: '大', value: '18px' }
}
```

**CSS 变量系统**:

```css
:root {
  --font-size-small: 14px;
  --font-size-medium: 16px;
  --font-size-large: 18px;
  --font-size-base: var(--font-size-medium);
  
  /* 所有字体基于基准相对计算 */
  --font-size-title: calc(var(--font-size-base) + 8px);
  --font-size-sm: calc(var(--font-size-base) - 2px);
}

[data-font-size="small"] { --font-size-base: var(--font-size-small); }
[data-font-size="medium"] { --font-size-base: var(--font-size-medium); }
[data-font-size="large"] { --font-size-base: var(--font-size-large); }
```

---

## 核心功能实现

### 1. 地理围栏智能匹配

#### 核心算法

**位置**: `pkg/geo/polygon.go`

**射线法实现**:

```go
func PointInPolygon(point Point, polygon []Point) bool {
    if len(polygon) < 3 {
        return false
    }
    
    inside := false
    j := len(polygon) - 1
    
    for i := 0; i < len(polygon); i++ {
        pi := polygon[i]
        pj := polygon[j]
        
        // 边界处理：点在边上
        if pointOnSegment(point, pj, pi) {
            return true
        }
        
        // 射线法核心
        intersects := ((pi.Lat > point.Lat) != (pj.Lat > point.Lat)) &&
            (point.Lng < (pj.Lng-pi.Lng)*(point.Lat-pi.Lat)/(pj.Lat-pi.Lat)+pi.Lng)
        if intersects {
            inside = !inside
        }
        j = i
    }
    
    return inside
}
```

#### 性能优化

**BoundingBox 预过滤**:

```go
type BoundingBox struct {
    MinLat, MaxLat float64
    MinLng, MaxLng float64
}

// O(1) 快速排除
func (b BoundingBox) Contains(p Point) bool {
    return p.Lat >= b.MinLat && p.Lat <= b.MaxLat &&
           p.Lng >= b.MinLng && p.Lng <= b.MaxLng
}
```

**引擎创建时预计算**:

```go
func NewEngine(zones []Zone) *Engine {
    filtered := make([]Zone, 0, len(zones))
    
    for _, zone := range zones {
        if len(zone.Points) < 3 {
            continue
        }
        zone.Box = NewBoundingBox(zone.Points)  // 预计算外包矩形
        filtered = append(filtered, zone)
    }
    
    // 按优先级降序排列
    sort.SliceStable(filtered, func(i, j int) bool {
        return filtered[i].Priority > filtered[j].Priority
    })
    
    return &Engine{zones: filtered}
}
```

#### 内存缓存机制

**位置**: `internal/service/geofence_service.go`

```go
type GeofenceService struct {
    zoneRepo *repository.ZoneRepository
    mu       sync.RWMutex      // 读写锁
    engine   *geo.Engine       // 内存中的围栏引擎
}

// 从数据库重新加载
func (s *GeofenceService) Reload() error {
    zones, _ := s.zoneRepo.ListActive()
    engineZones := convertToEngineZones(zones)
    
    engine := geo.NewEngine(engineZones)
    
    s.mu.Lock()
    s.engine = engine
    s.mu.Unlock()
    
    return nil
}

// 并发安全读取
func (s *GeofenceService) Match(lat, lng float64) (int64, bool) {
    s.mu.RLock()
    engine := s.engine
    s.mu.RUnlock()
    
    return engine.Match(geo.Point{Lat: lat, Lng: lng})
}
```

**缓存刷新触发点**:
- 启动时预加载
- 围栏创建/更新/删除时自动刷新

#### 兜底规则

**Haversine 距离计算**（`internal/service/geo.go`）:

```go
const earthRadiusMeters = 6371000

func HaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
    lat1Rad := lat1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    deltaLat := (lat2 - lat1) * math.Pi / 180
    deltaLng := (lng2 - lng1) * math.Pi / 180
    
    a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
         math.Cos(lat1Rad)*math.Cos(lat2Rad)*
         math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return earthRadiusMeters * c
}
```

### 2. 五表权限管理系统

#### 表结构关系

```
permissions (权限定义)
    ↑
    │ permission_code
    │
menus (前端菜单)

roles (角色定义)
    ↑
    │ role_id
    │
role_permissions (角色-权限关联)
    ↑
    │ code
    │
user_identities (用户-身份关联)
    ↑
    │ user_id
    │
users (用户表)
```

#### 权限检查流程

```
请求进入
    ↓
Admin 身份？
    ├─ 是 → 跳过所有检查
    └─ 否 ↓
公共 API (is_public=1)？
    ├─ 是 → 放行
    └─ 否 ↓
检查 API 权限
    ├─ 遍历用户所有角色
    ├─ 获取角色的所有权限
    ├─ 匹配 API 路径和方法
    └─ 有权限 → 继续 | 无权限 → 403
```

#### 多身份系统

**user_identities 表**:

```sql
CREATE TABLE user_identities (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    identity_type VARCHAR(20) NOT NULL,  -- 直接对应 roles.code
    is_primary TINYINT(1) DEFAULT 0,      -- 主身份（每用户有且仅有一个）
    station_id BIGINT,                    -- 所属站点（B端身份）
    status VARCHAR(20) DEFAULT 'active',
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    granted_by BIGINT,
    revoked_at DATETIME,
    
    UNIQUE KEY(user_id, identity_type, deleted_at)
);
```

**多角色权限并集**:

```go
// 获取用户的所有权限码（多角色取并集）
func (s *PermissionService) GetUserPermissionCodes(roleCodes []string) ([]string, error) {
    roles, _ := s.roleRepo.GetByCodes(roleCodes)
    
    permIDs, _ := s.rolePermRepo.GetPermissionIDsByRoleIDs(roleIDs)
    perms, _ := s.permRepo.GetByIDs(permIDs)
    
    codes := make([]string, len(perms))
    for i, p := range perms {
        codes[i] = p.Code
    }
    return codes, nil
}
```

### 3. 双端认证与隔离

#### JWT Token 结构

**位置**: `pkg/jwt/manager.go`

```go
type Claims struct {
    UserID     int64    `json:"uid"`
    Identities []string `json:"identities,omitempty"` // 用户的所有身份
    Primary    string   `json:"primary,omitempty"`    // 主身份
    Type       string   `json:"type"`                 // "b_end" | "c_end"
    StationID  int64    `json:"station_id"`            // B端有值，C端为0
    jwtlib.RegisteredClaims  // 包含 ID(jti)、ExpiresAt、IssuedAt
}
```

#### B 端 Token 示例

```json
{
  "uid": 4,
  "identities": ["staff"],
  "primary": "staff",
  "type": "b_end",
  "station_id": 1,
  "jti": "token-uuid",
  "exp": 1709123456
}
```

#### C 端 Token 示例

```json
{
  "uid": 9,
  "identities": ["elderly", "family"],
  "primary": "elderly",
  "type": "c_end",
  "station_id": 0,
  "jti": "token-uuid",
  "exp": 1709123456
}
```

#### Token 黑名单机制

**Redis 存储结构**:

```
# 单 Token 黑名单（登出）
token:blacklist:{jti} = "1"
TTL: Token 剩余有效期

# 用户级 Token 撤销（角色变更）
user:revoked:{uid}:{endType} = {撤销时间戳}
TTL: 24 小时
```

**检查逻辑** (`internal/middleware/auth.go`):

```go
// 1. 检查单个 Token 是否在黑名单
isBlacklisted, _ := blacklistService.IsBlacklisted(ctx, claims.ID)
if isBlacklisted {
    c.AbortWithStatusJSON(401, "token revoked")
    return
}

// 2. 检查用户级 Token 撤销
isRevoked, _ := blacklistService.IsUserTokenRevoked(
    ctx, claims.UserID, claims.Type, claims.IssuedAt.Unix(),
)
if isRevoked {
    c.AbortWithStatusJSON(401, "user tokens revoked due to role change")
    return
}
```

### 4. 通知系统

#### 架构设计

```
业务触发层
    ├─ RequestService (需求创建/取消)
    ├─ TaskService (任务认领/完成/转派)
    └─ ...
         │
         │ go func() { } (异步，无队列)
         ▼
    NotificationService
         ├─ SendInApp() → 创建通知记录
         └─ SendEmail() → SMTPSender 发送
                │
                ▼
         ┌────────┴────────┐
         ▼                 ▼
    NotificationRepo   SMTPSender
         │                 │
         ▼                 ▼
    MySQL           SMTP 服务器
```

#### 邮件发送实现

**位置**: `internal/notify/smtp_sender.go`

**特性**:
- 支持 TLS 加密连接
- 连接超时控制（10 秒）
- 纯文本格式（`text/plain; charset=UTF-8`）

**配置**:

```go
type MailConfig struct {
    SMTPHost string
    SMTPPort int
    SMTPUser string
    SMTPPass string
    SMTPTLS  bool
}
```

#### 通知触发时机

| 业务事件 | 触发位置 | 通知对象 | 渠道 |
|---------|---------|---------|------|
| 创建需求 | RequestService.Create() | 站点管理员 | 站内信 |
| 任务认领 | TaskService.Claim() | 需求用户 | 站内信 |
| 任务完成 | TaskService.Complete() | 需求用户 | 站内信 + 邮件 |
| 任务转派 | TaskService.Transfer() | 新工作人员 | 站内信 |

**当前限制**:
- 无消息队列（goroutine 异步发送）
- 无重试机制
- 发送状态未更新（`send_status` 始终为 `pending`）

### 5. 文件存储系统

#### 存储抽象接口

**位置**: `internal/storage/storage.go`

```go
type Provider interface {
    Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error)
    SignedURL(ctx context.Context, key string, expire time.Duration) (string, error)
    Delete(ctx context.Context, key string) error
}
```

#### 本地存储实现

**位置**: `internal/storage/local/local.go`

**路径生成规则**:

```go
// 最终格式: {module}/{YYYYMMDD}/{timestamp}.{ext}
// 示例: c_end/20260228/1709123456789012345.jpg
func buildObjectKey(module, filename string) string {
    ext := filepath.Ext(filename)
    if ext == "" {
        ext = ".bin"
    }
    date := time.Now().Format("20060102")
    name := fmt.Sprintf("%d", time.Now().UnixNano())
    
    return fmt.Sprintf("%s/%s/%s%s", module, date, name, ext)
}
```

**安全特性**:
- 路径遍历攻击防护（`filepath.Clean` + 前缀检查）
- 模块名称清洗（移除 `..`、`\`、`//`）

#### 阿里云 OSS 实现

**位置**: `internal/storage/oss/oss.go`

**签名 URL 生成**:

```go
func (p *Provider) SignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
    seconds := int64(expire.Seconds())
    if seconds <= 0 {
        seconds = 3600  // 默认 1 小时
    }
    return p.bucket.SignURL(key, aliyunoss.HTTPGet, seconds)
}
```

#### 工厂模式

**位置**: `internal/storage/factory.go`

```go
func NewProvider(cfg config.StorageConfig) (Provider, error) {
    switch strings.ToLower(cfg.Driver) {
    case "local":
        return local.New(cfg.Local.BasePath, cfg.Local.BaseURL)
    case "oss":
        return oss.New(oss.Config{
            Endpoint: cfg.OSS.Endpoint,
            Bucket: cfg.OSS.Bucket,
            // ...
        })
    default:
        return nil, errors.New("unsupported storage driver")
    }
}
```

**环境变量配置**:

```bash
# 存储驱动
STORAGE_DRIVER=local  # local | oss

# 本地存储
STORAGE_LOCAL_BASE_PATH=./storage
STORAGE_LOCAL_BASE_URL=http://localhost:8080/static

# OSS 存储
STORAGE_OSS_ENDPOINT=oss-cn-beijing.aliyuncs.com
STORAGE_OSS_BUCKET=your-bucket
STORAGE_OSS_ACCESS_KEY_ID=your-key
STORAGE_OSS_ACCESS_KEY_SECRET=your-secret
```

#### B/C 端上传限制

| 限制项 | C 端 | B 端 |
|--------|-----|-----|
| 文件大小 | 5MB | 50MB |
| 文件类型 | 仅图片 | 多种文件 |
| 默认模块 | `c_end` | `b_end` |

---

## 数据库设计

### 核心业务表

#### 1. 用户与身份

**users** (用户表):
```sql
id BIGINT PRIMARY KEY
phone VARCHAR(20) UNIQUE
password_hash VARCHAR(255)
name VARCHAR(50)
avatar VARCHAR(255)
status VARCHAR(20) DEFAULT 'active'
created_at, updated_at, deleted_at
```

**user_identities** (用户-身份关联):
```sql
id BIGINT PRIMARY KEY
user_id BIGINT
identity_type VARCHAR(20)  -- admin/station_manager/staff/elderly/family
is_primary TINYINT(1)       -- 主身份
station_id BIGINT           -- 所属站点（B端）
status VARCHAR(20)          -- active/inactive
granted_at DATETIME
granted_by BIGINT
revoked_at DATETIME
```

#### 2. 权限系统

**permissions** (权限定义):
```sql
id BIGINT PRIMARY KEY
code VARCHAR(100) UNIQUE      -- 格式: module:resource:action
name VARCHAR(100)
type VARCHAR(20)             -- menu/button/resource
parent_id BIGINT              -- 树形结构
api_path VARCHAR(200)        -- API 路径（支持通配符 *）
api_method VARCHAR(20)       -- GET/POST/PUT/DELETE
is_public TINYINT(1)          -- 公共权限
module VARCHAR(50)           -- 所属模块
status VARCHAR(20)
sort INT
```

**roles** (角色定义):
```sql
id BIGINT PRIMARY KEY
code VARCHAR(50) UNIQUE      -- admin/station_manager/staff/elderly/family
name VARCHAR(100)
is_system TINYINT(1)         -- 系统内置角色
status VARCHAR(20)
sort INT
```

**role_permissions** (角色-权限关联):
```sql
id BIGINT PRIMARY KEY
role_id BIGINT
permission_id BIGINT
UNIQUE(role_id, permission_id)
```

**menus** (前端菜单):
```sql
id BIGINT PRIMARY KEY
parent_id BIGINT              -- 父菜单，0 表示顶级
name VARCHAR(50)
path VARCHAR(200)             -- 路由路径
component VARCHAR(200)        -- 组件路径
icon VARCHAR(50)
permission_code VARCHAR(100)  -- 关联权限码
sort INT
hidden TINYINT(1)
status VARCHAR(20)
```

#### 3. 服务站点与围栏

**service_stations** (服务站点):
```sql
id BIGINT PRIMARY KEY
name VARCHAR(100)
code VARCHAR(50) UNIQUE
address VARCHAR(200)
latitude DECIMAL(10,7)
longitude DECIMAL(10,7)
capacity INT                  -- 服务容量
work_hours VARCHAR(100)
status VARCHAR(20)
```

**service_zones** (地理围栏):
```sql
id BIGINT PRIMARY KEY
station_id BIGINT
name VARCHAR(100)
points JSON                   -- 多边形顶点 [[lng,lat], ...]
priority INT                  -- 优先级（数值越大越优先）
status VARCHAR(20)
```

#### 4. 服务请求与任务

**service_requests** (服务需求):
```sql
id BIGINT PRIMARY KEY
request_no VARCHAR(50) UNIQUE  -- 需求编号
user_id BIGINT
service_type VARCHAR(20)       -- daily_care/medical_care/meal_service/...
status VARCHAR(20)             -- pending/dispatched/claimed/completed/cancelled
description TEXT
submit_location_lat DECIMAL(10,7)
submit_location_lng DECIMAL(10,7)
contact_name VARCHAR(50)
contact_phone VARCHAR(20)
address VARCHAR(200)
appointment_time DATETIME
station_id BIGINT               -- 匹配的站点
reject_reason TEXT
images JSON                    -- 图片URL列表
```

**task_assignments** (任务分配):
```sql
id BIGINT PRIMARY KEY
request_id BIGINT
station_id BIGINT
staff_id BIGINT                -- 认领工作人员
status VARCHAR(20)             -- pending/claimed/processing/completed/cancelled
claimed_at DATETIME
completed_at DATETIME
rating INT                     -- 服务评分 1-5
feedback TEXT
staff_notes TEXT
images JSON
```

#### 5. 通知系统

**notifications** (通知):
```sql
id BIGINT PRIMARY KEY
user_id BIGINT
title VARCHAR(100)
body TEXT
type VARCHAR(20)
related_id BIGINT
related_type VARCHAR(20)
channel VARCHAR(20)            -- in_app/email
send_status VARCHAR(20)        -- pending/sent/failed
sent_at DATETIME
is_read TINYINT(1)
read_at DATETIME
retry_count BIGINT
```

### 数据库约定

1. **逻辑外键**: 所有关联字段使用逻辑外键，不使用数据库外键约束
2. **软删除**: 所有表支持软删除（`deleted_at` 字段）
3. **字符集**: `utf8mb4` + `utf8mb4_unicode_ci`
4. **索引**: 关键查询字段建立索引
5. **JSON 存储**: 地理围栏、图片列表使用 JSON 类型

---

## 安全机制

### 1. 认证安全

#### 密码存储

- 使用 bcrypt 哈希算法
- `password_hash` 字段存储哈希值
- 原始密码**不存储**

#### Token 安全

- Access Token 过期时间：可配置（通常较短）
- Refresh Token 过期时间：可配置（通常较长）
- Token 唯一标识（`jti`）支持黑名单

### 2. 端隔离安全

#### JWT Token 类型区分

```json
{
  "type": "b_end"  // 或 "c_end"
}
```

#### 中间件强制检查

```
B 端路由 → RequireEndType("b_end") → 拒绝 C 端 Token
C 端路由 → RequireEndType("c_end") → 拒绝 B 端 Token
```

### 3. 权限安全

#### RBAC 权限模型

- 三表关联：`permissions` ↔ `roles` ↔ `role_permissions`
- 最小权限原则：用户仅拥有其角色所需的权限
- Admin 特殊处理：硬编码跳过所有权限检查

#### API 级别权限控制

- 每个 API 端点可配置权限码
- 支持通配符匹配（如 `/api/v1/b/tasks/*/claim`）
- 公共 API（`is_public=1`）所有登录用户可访问

### 4. 数据安全

#### 敏感字段加密

- `id_card_hmac`: 身份证号 HMAC 摘要
- `id_card_masked`: 身份证号脱敏值

#### 软删除

- 所有表支持 `deleted_at` 字段
- 查询时自动过滤已删除数据

### 5. 网络安全

#### CORS 配置

```go
// 允许的前端域名
AllowOrigins: []string{"http://localhost:3001", "http://localhost:5174"}
AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}
```

#### HTTPS 生产环境

- 生产环境强制使用 HTTPS
- TLS 加密连接（数据库、SMTP）

---

## 部署架构

### 开发环境

```bash
# 启动基础设施
docker compose up -d  # MySQL + Redis

# 启动后端
cd backend
air  # 热重载（端口 8080）

# 启动 B 端
cd frontend/management-portal
npm run dev  # 端口 3001

# 启动 C 端
cd frontend/c-end
npm run dev  # 端口 5174
```

### 生产环境

#### Docker Compose 部署

```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: scare_db
  
  redis:
    image: redis:7.0
  
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis
  
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./deployment/nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./storage:/app/storage
```

#### Kubernetes 部署（可选）

- `deployment/kubernetes/` 目录包含 K8s 配置
- 支持 Service、Deployment、Ingress
- ConfigMap 和 Secret 管理配置

### 环境变量配置

```bash
# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_USER=scare_user
DB_PASSWORD=scare_pass
DB_NAME=scare_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h
JWT_REFRESH_EXPIRES_IN=720h

# SMTP
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=noreply@example.com
SMTP_PASS=your-password
SMTP_TLS=true

# 存储
STORAGE_DRIVER=local  # local | oss
STORAGE_LOCAL_BASE_PATH=./storage
STORAGE_LOCAL_BASE_URL=http://localhost:8080/static

# OSS（可选）
STORAGE_OSS_ENDPOINT=oss-cn-beijing.aliyuncs.com
STORAGE_OSS_BUCKET=your-bucket
STORAGE_OSS_ACCESS_KEY_ID=your-key
STORAGE_OSS_ACCESS_KEY_SECRET=your-secret
```

---

## 开发规范

### 后端开发规范

#### 分层职责

| 层 | 职责 | 禁止 |
|---|------|------|
| Handler | 参数校验、响应格式化 | 直接访问数据库 |
| Service | 业务逻辑、事务处理 | 包含 HTTP 或 GORM 细节 |
| Repository | 数据库查询封装 | 包含业务逻辑 |
| Model | 数据模型（GORM Gen 生成） | 手动修改 `.gen.go` 文件 |

#### 命名约定

- 文件名：小写下划线（`user_service.go`）
- 接口名：大驼峰（`UserService`）
- 方法名：大驼峰（`GetByID`、`Create`）
- 常量：大写下划线（`MAX_UPLOAD_SIZE`）

#### 错误处理

```go
// 统一错误响应
{ "msg": "错误描述", "data": null }

// Service 层返回错误
return errors.New("用户不存在")

// Handler 层转换
RespondError(c, 400, err.Error())
```

#### 数据库查询

```go
// 使用 GORM Gen 生成的查询
q := query.Use(db).User

// 简单查询
user, err := q.Where(q.ID.Eq(userID)).First()

// 复杂查询
users, err := q.Where(
    q.Status.Eq("active"),
    q.StationID.Neq(0),
).Order(q.CreatedAt.Desc()).Find()
```

### 前端开发规范

#### 组件风格

- 统一使用 Composition API（`<script setup lang="ts">`）
- TypeScript 严格模式（`strict: true`）
- 单文件组件（`.vue`）

#### 目录组织

- B 端页面：`src/pages/`
- C 端页面：`src/views/`
- 共享组件：`src/components/`
- 路径别名：`@` → `src/`

#### 状态管理

```typescript
// Composition API 风格
export const useUserStore = defineStore('user', () => {
  const state = ref({ ... })
  const getters = { ... }
  const actions = { ... }
  return { state, ...getters, ...actions }
})
```

#### API 调用

```typescript
// B 端：集中管理
import { authApi, taskApi } from '@/api'

// C 端：模块化
import { authAPI, requestsAPI } from '@/api'
```

---

## 已知问题与改进建议

### 1. 通知系统

**当前问题**:
- 无消息队列（goroutine 异步，服务重启丢失通知）
- 无重试机制
- 发送状态未更新（`send_status` 始终为 `pending`）

**改进建议**:
- 引入 Redis 队列实现持久化
- 实现定时任务扫描失败通知重发
- 发送成功后更新 `send_status`、`sent_at`
- 引入 HTML 邮件模板

### 2. API 层

**当前问题**:
- B 端所有接口集中在单文件（572 行）
- C 端与 B 端有大量重复代码

**改进建议**:
- B 端拆分为模块化文件（参考 C 端）
- 抽取共享代码到 `frontend/shared/`

### 3. 测试覆盖

**当前问题**:
- 前端无自动化测试（B 端仅 Playwright 配置）
- 后端部分模块缺少单元测试

**改进建议**:
- 补充后端 Service 层单元测试
- 实现前端关键组件的单元测试（Vitest）
- 扩展 Playwright E2E 测试覆盖

### 4. 文档

**当前问题**:
- 缺少 CHANGELOG.md
- 缺少 CONTRIBUTING.md

**改进建议**:
- 建立 CHANGELOG 记录版本变更
- 编写贡献指南（PR 流程、代码规范）

---

## 附录

### 测试账号

| 角色 | 手机号 | 密码 | 说明 |
|------|--------|------|------|
| Admin | 13800000001 | Test@123 | 系统管理员（全局权限） |
| Station Manager | 13800000002 | Test@123 | 站点负责人 |
| Staff | 13800000004 | Test@123 | 工作人员 |

### API 文档

- Swagger UI: http://localhost:8080/swagger/index.html
- 文档位置: `backend/docs/04-API接口设计.md`

### 相关文档

- **项目级**: `docs/README.md`
- **后端**: `backend/docs/01-开发指南.md`
- **前端 B 端**: `frontend/management-portal/docs/README.md`
- **前端 C 端**: `frontend/c-end/docs/README.md`
- **数据库**: `backend/database/docs/README.md`

---

**文档维护者**: sCare 开发团队  
**最后更新**: 2026-02-28
