# sCare 多角色重构规格说明

**版本**: v1.1.0
**创建日期**: 2026-01-18
**状态**: 待执行

---

## 一、核心设计原则

### 1.1 角色分配规则

**端内角色唯一性**：
- 一个用户可以拥有多个角色，但**在同一端只能有一个角色**
- C端角色：`elderly`（老年人），`family`（家属-仅联系人）
- B端角色：`staff`（工作人员），`station_manager`（站长），`admin`（管理员）
- 跨端多角色示例：
  - 用户A：`['elderly', 'staff']` - C端是老年人，B端是工作人员
  - 用户B：`['family']` - 纯家属，仅作为备用联系人
  - 用户C：`['station_manager']` - 仅B端角色

**角色语义**：
- `family`: 纯家属角色，**不承担业务流转**，仅作为老年人的备用联系人
- 其他角色：承担具体业务流程（提交需求、处理任务、管理站点等）

### 1.2 权限模型（Casbin RBAC）

**无角色继承**：
- 不使用 Casbin 角色继承特性（`g` 关系）
- 每个角色单独分配具体权限（`p` 策略）
- 权限类型：
  - **接口权限**：API 路径 + HTTP 方法（后端校验）
  - **按钮权限**：前端 UI 元素显示控制（前端校验）

**权限策略示例**：
```
p, staff, /api/b/tasks/pool, GET
p, staff, /api/b/tasks/:id/claim, POST
p, station_manager, /api/b/dashboard, GET
p, station_manager, /api/b/tasks/*, *
p, admin, /api/b/users, POST
```

### 1.3 登录与认证

**不同登录接口**：
- **C端登录**：`POST /api/c/auth/login`
  - 查询用户在 C端 的角色（`elderly` 或 `family`）
  - 如果无C端角色，返回 `403 该账号无C端访问权限`
  - 颁发包含C端角色的 JWT Token

- **B端登录**：`POST /api/b/auth/login`
  - 查询用户在 B端 的角色（`staff`、`station_manager`、`admin`）
  - 如果无B端角色，返回 `403 该账号无B端访问权限`
  - 颁发包含B端角色的 JWT Token

**JWT Token 设计**（单角色，端独立）：
```json
{
  "uid": 4,
  "role": "staff",        // 该端的唯一角色
  "station_id": 1,
  "exp": 1737216000,
  "iat": 1737129600
}
```

**Token 特性**：
- C端 Token 和 B端 Token 相互独立
- 一个用户可以同时持有两个端的有效 Token
- Token 不包含"所有角色"，只包含当前端角色
- 无需角色切换（端内唯一角色）

### 1.4 实时角色撤销

**需求**：角色撤销立即生效

**实现方案**：JWT + Redis 黑名单 + 数据库缓存
```go
// 撤销角色时
func RevokeUserRole(userID uint, role string) error {
    // 1. 数据库更新
    db.Model(&UserRole{}).
       Where("user_id = ? AND role = ?", userID, role).
       Update("status", "inactive")

    // 2. 将用户加入黑名单（强制重新登录）
    redis.Set(fmt.Sprintf("token_blacklist:user:%d", userID), 1, 24*time.Hour)

    // 3. 清除角色缓存
    redis.Del(fmt.Sprintf("user_roles:%d", userID))

    return nil
}

// 权限检查中间件
func CheckPermission(c *gin.Context) {
    uid := c.GetUint("uid")
    jwtRole := c.GetString("role")

    // 1. 检查黑名单（立即生效）
    if redis.Exists(fmt.Sprintf("token_blacklist:user:%d", uid)) {
        c.JSON(401, gin.H{"msg": "token_revoked"})
        c.Abort()
        return
    }

    // 2. 验证角色仍然激活（缓存5分钟）
    activeRoles := getUserActiveRolesCache(uid) // Redis缓存
    if !contains(activeRoles, jwtRole) {
        c.JSON(403, gin.H{"msg": "role_revoked"})
        c.Abort()
        return
    }

    // 3. Casbin 权限检查
    ok, _ := enforcer.Enforce(jwtRole, c.Request.URL.Path, c.Request.Method)
    if !ok {
        c.JSON(403, gin.H{"msg": "forbidden"})
        c.Abort()
        return
    }

    c.Next()
}
```

**性能开销**：
- 每次请求：Redis 查询 1次（黑名单）+ 1次（角色缓存）≈ 2-3ms
- 可接受范围（大部分业务接口响应时间 50-200ms）

---

## 二、数据库设计

### 2.1 已有表改造

#### users 表
```sql
-- 保留原 role 字段（兼容性，可标记为 deprecated）
-- 新增 primary_role 字段
ALTER TABLE users
  ADD COLUMN primary_role VARCHAR(20) DEFAULT NULL COMMENT '主角色（最常用）' AFTER role;

-- primary_role 用于快速查询用户"主要身份"，减少 JOIN
```

#### user_roles 表（已创建）
```sql
CREATE TABLE user_roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL COMMENT '用户ID',
  role VARCHAR(20) NOT NULL COMMENT '角色',
  is_primary BOOLEAN DEFAULT FALSE COMMENT '是否为主角色',
  status VARCHAR(20) DEFAULT 'active' COMMENT '角色状态(active/inactive)',
  granted_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
  granted_by BIGINT DEFAULT NULL COMMENT '授权人ID',
  revoked_at DATETIME DEFAULT NULL COMMENT '撤销时间',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE KEY uk_user_role (user_id, role),
  KEY idx_user_id (user_id),
  KEY idx_role (role),
  KEY idx_status (status),

  CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE
) COMMENT='用户角色关联表';
```

**数据约束**（应用层保证）：
- 端内角色唯一：一个用户在C端或B端只能有一个激活角色
- C端角色：`elderly`, `family`
- B端角色：`staff`, `station_manager`, `admin`

### 2.2 Casbin 策略表

```sql
-- Casbin 标准表结构
CREATE TABLE casbin_rule (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  ptype VARCHAR(10) NOT NULL COMMENT '策略类型(p/g)',
  v0 VARCHAR(256) DEFAULT NULL COMMENT '角色',
  v1 VARCHAR(256) DEFAULT NULL COMMENT '资源路径',
  v2 VARCHAR(256) DEFAULT NULL COMMENT 'HTTP方法',
  v3 VARCHAR(256) DEFAULT NULL,
  v4 VARCHAR(256) DEFAULT NULL,
  v5 VARCHAR(256) DEFAULT NULL,

  KEY idx_ptype (ptype),
  KEY idx_v0 (v0),
  KEY idx_v1 (v1)
) COMMENT='Casbin权限策略表';

-- 初始化权限策略
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
-- B端 - staff 权限
('p', 'staff', '/api/b/tasks/pool', 'GET'),
('p', 'staff', '/api/b/tasks/my', 'GET'),
('p', 'staff', '/api/b/tasks/:id', 'GET'),
('p', 'staff', '/api/b/tasks/:id/claim', 'POST'),
('p', 'staff', '/api/b/tasks/:id/complete', 'POST'),
('p', 'staff', '/api/b/requests/:id', 'GET'),

-- B端 - station_manager 权限
('p', 'station_manager', '/api/b/tasks/*', '*'),
('p', 'station_manager', '/api/b/dashboard', 'GET'),
('p', 'station_manager', '/api/b/stations/:id', 'GET'),
('p', 'station_manager', '/api/b/stations/:id', 'PUT'),
('p', 'station_manager', '/api/b/requests/*', '*'),

-- B端 - admin 权限
('p', 'admin', '/api/b/*', '*'),
('p', 'admin', '/api/common/users', 'POST'),
('p', 'admin', '/api/common/users/:id', 'PUT'),
('p', 'admin', '/api/common/users/:id', 'DELETE'),

-- C端 - elderly 权限
('p', 'elderly', '/api/c/requests', 'POST'),
('p', 'elderly', '/api/c/requests/my', 'GET'),
('p', 'elderly', '/api/c/requests/:id', 'GET'),
('p', 'elderly', '/api/c/profile', 'GET'),
('p', 'elderly', '/api/c/profile', 'PUT'),

-- C端 - family 权限（仅查看）
('p', 'family', '/api/c/profile', 'GET');
```

---

## 三、后端改造

### 3.1 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── c_end/              # C端 MVC
│   │   ├── controller/
│   │   │   ├── request_controller.go
│   │   │   └── profile_controller.go
│   │   ├── service/
│   │   │   ├── request_service.go
│   │   │   └── profile_service.go
│   │   └── router/
│   │       └── router.go   # 注册 /api/c/* 路由
│   │
│   ├── b_end/              # B端 MVC
│   │   ├── controller/
│   │   │   ├── task_controller.go
│   │   │   ├── dashboard_controller.go
│   │   │   └── station_controller.go
│   │   ├── service/
│   │   │   ├── task_service.go
│   │   │   └── station_service.go
│   │   └── router/
│   │       └── router.go   # 注册 /api/b/* 路由
│   │
│   ├── common/             # 共享服务
│   │   ├── user/
│   │   │   └── user_service.go
│   │   ├── auth/
│   │   │   ├── auth_service.go
│   │   │   └── jwt.go
│   │   └── upload/
│   │       └── upload_service.go
│   │
│   ├── model/              # 数据模型
│   │   ├── user.go
│   │   ├── user_role.go
│   │   ├── task.go
│   │   └── request.go
│   │
│   └── middleware/         # 中间件
│       ├── casbin.go       # 权限检查
│       ├── jwt.go          # JWT 验证
│       └── cors.go
│
├── pkg/
│   ├── cache/              # Redis 缓存
│   └── db/                 # 数据库连接
│
└── configs/
    ├── casbin_model.conf
    └── config.yaml
```

### 3.2 数据模型改造

#### model/user.go
```go
package model

import "time"

type User struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    Phone       string     `gorm:"uniqueIndex;size:20" json:"phone"`
    PasswordHash string    `gorm:"size:255" json:"-"`
    Name        string     `gorm:"size:50" json:"name"`
    Email       string     `gorm:"size:100" json:"email"`
    Role        string     `gorm:"size:20;comment:已废弃" json:"role"` // 兼容旧代码
    PrimaryRole string     `gorm:"size:20" json:"primary_role"`
    StationID   *uint      `json:"station_id"`
    Status      string     `gorm:"size:20;default:active" json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `gorm:"index" json:"-"`

    // 关联
    Roles       []UserRole `gorm:"foreignKey:UserID" json:"roles,omitempty"`
}

// 获取用户在指定端的角色
func (u *User) GetRoleByType(roleType string) string {
    var targetRoles []string
    if roleType == "C端" {
        targetRoles = []string{"elderly", "family"}
    } else if roleType == "B端" {
        targetRoles = []string{"staff", "station_manager", "admin"}
    }

    for _, ur := range u.Roles {
        if ur.Status == "active" && contains(targetRoles, ur.Role) {
            return ur.Role
        }
    }
    return ""
}

// 检查是否拥有某个角色
func (u *User) HasRole(role string) bool {
    for _, ur := range u.Roles {
        if ur.Role == role && ur.Status == "active" {
            return true
        }
    }
    return false
}
```

#### model/user_role.go
```go
package model

import "time"

type UserRole struct {
    ID        uint       `gorm:"primaryKey" json:"id"`
    UserID    uint       `gorm:"index" json:"user_id"`
    Role      string     `gorm:"index;size:20" json:"role"`
    IsPrimary bool       `gorm:"default:false" json:"is_primary"`
    Status    string     `gorm:"default:active;size:20" json:"status"`
    GrantedAt time.Time  `json:"granted_at"`
    GrantedBy *uint      `json:"granted_by"`
    RevokedAt *time.Time `json:"revoked_at"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
```

### 3.3 共享服务实现

#### common/user/user_service.go
```go
package user

import (
    "gorm.io/gorm"
    "scare/internal/model"
)

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

// 获取用户信息（预加载角色）
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
    var user model.User
    err := s.db.Preload("Roles", "status = ?", "active").
        First(&user, id).Error
    return &user, err
}

// 获取用户在指定端的角色
func (s *UserService) GetUserRoleByType(userID uint, roleType string) (string, error) {
    user, err := s.GetUserByID(userID)
    if err != nil {
        return "", err
    }
    return user.GetRoleByType(roleType), nil
}

// 获取用户所有激活角色
func (s *UserService) GetUserActiveRoles(userID uint) ([]string, error) {
    var roles []model.UserRole
    err := s.db.Where("user_id = ? AND status = ?", userID, "active").
        Find(&roles).Error

    roleNames := make([]string, len(roles))
    for i, r := range roles {
        roleNames[i] = r.Role
    }
    return roleNames, err
}

// 验证用户凭证
func (s *UserService) ValidateCredentials(phone, password string) (*model.User, error) {
    var user model.User
    err := s.db.Preload("Roles", "status = ?", "active").
        Where("phone = ?", phone).First(&user).Error
    if err != nil {
        return nil, err
    }

    // 验证密码（bcrypt）
    if !checkPasswordHash(password, user.PasswordHash) {
        return nil, errors.New("invalid_credentials")
    }

    return &user, nil
}
```

#### common/auth/auth_service.go
```go
package auth

import (
    "errors"
    "scare/internal/common/user"
    "scare/pkg/cache"
)

type AuthService struct {
    userService *user.UserService
    cache       *cache.RedisCache
}

// C端登录
func (s *AuthService) LoginCEnd(phone, password string) (string, error) {
    // 1. 验证凭证
    user, err := s.userService.ValidateCredentials(phone, password)
    if err != nil {
        return "", err
    }

    // 2. 获取C端角色
    role := user.GetRoleByType("C端")
    if role == "" {
        return "", errors.New("no_c_end_permission")
    }

    // 3. 生成 JWT Token
    token, err := GenerateJWT(user.ID, role, user.StationID)
    if err != nil {
        return "", err
    }

    return token, nil
}

// B端登录
func (s *AuthService) LoginBEnd(phone, password string) (string, error) {
    // 1. 验证凭证
    user, err := s.userService.ValidateCredentials(phone, password)
    if err != nil {
        return "", err
    }

    // 2. 获取B端角色
    role := user.GetRoleByType("B端")
    if role == "" {
        return "", errors.New("no_b_end_permission")
    }

    // 3. 生成 JWT Token
    token, err := GenerateJWT(user.ID, role, user.StationID)
    if err != nil {
        return "", err
    }

    return token, nil
}

// 撤销用户角色（立即生效）
func (s *AuthService) RevokeUserRole(userID uint, role string) error {
    // 1. 数据库更新
    err := s.db.Model(&model.UserRole{}).
        Where("user_id = ? AND role = ?", userID, role).
        Update("status", "inactive").Error
    if err != nil {
        return err
    }

    // 2. 加入黑名单（所有 Token 失效）
    s.cache.Set(fmt.Sprintf("token_blacklist:user:%d", userID), 1, 24*time.Hour)

    // 3. 清除角色缓存
    s.cache.Del(fmt.Sprintf("user_roles:%d", userID))

    return nil
}
```

#### common/auth/jwt.go
```go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UID       uint   `json:"uid"`
    Role      string `json:"role"`
    StationID *uint  `json:"station_id,omitempty"`
    jwt.RegisteredClaims
}

// 生成 JWT Token
func GenerateJWT(uid uint, role string, stationID *uint) (string, error) {
    claims := Claims{
        UID:       uid,
        Role:      role,
        StationID: stationID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}

// 验证 JWT Token
func ValidateJWT(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid_token")
}
```

### 3.4 认证接口实现

#### c_end/controller/auth_controller.go
```go
package controller

import (
    "github.com/gin-gonic/gin"
    "scare/internal/common/auth"
)

type CEndAuthController struct {
    authService *auth.AuthService
}

// C端登录
func (h *CEndAuthController) Login(c *gin.Context) {
    var req struct {
        Phone    string `json:"phone" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"msg": "invalid_request"})
        return
    }

    token, err := h.authService.LoginCEnd(req.Phone, req.Password)
    if err != nil {
        if err.Error() == "no_c_end_permission" {
            c.JSON(403, gin.H{"msg": "该账号无C端访问权限，请使用B端登录"})
            return
        }
        c.JSON(401, gin.H{"msg": "invalid_credentials"})
        return
    }

    c.JSON(200, gin.H{
        "msg": "ok",
        "data": gin.H{
            "token": token,
        },
    })
}
```

#### b_end/controller/auth_controller.go
```go
package controller

import (
    "github.com/gin-gonic/gin"
    "scare/internal/common/auth"
)

type BEndAuthController struct {
    authService *auth.AuthService
}

// B端登录
func (h *BEndAuthController) Login(c *gin.Context) {
    var req struct {
        Phone    string `json:"phone" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"msg": "invalid_request"})
        return
    }

    token, err := h.authService.LoginBEnd(req.Phone, req.Password)
    if err != nil {
        if err.Error() == "no_b_end_permission" {
            c.JSON(403, gin.H{"msg": "该账号无B端访问权限，请使用C端登录"})
            return
        }
        c.JSON(401, gin.H{"msg": "invalid_credentials"})
        return
    }

    c.JSON(200, gin.H{
        "msg": "ok",
        "data": gin.H{
            "token": token,
        },
    })
}
```

### 3.5 权限中间件改造

#### middleware/casbin.go
```go
package middleware

import (
    "fmt"
    "github.com/casbin/casbin/v2"
    "github.com/casbin/gorm-adapter/v3"
    "github.com/gin-gonic/gin"
    "scare/pkg/cache"
)

type CasbinMiddleware struct {
    enforcer *casbin.Enforcer
    cache    *cache.RedisCache
}

func NewCasbinMiddleware(db *gorm.DB, cache *cache.RedisCache) *CasbinMiddleware {
    // 使用 GORM Adapter 从数据库加载策略
    adapter, _ := gormadapter.NewAdapterByDB(db)
    enforcer, _ := casbin.NewEnforcer("configs/casbin_model.conf", adapter)

    return &CasbinMiddleware{
        enforcer: enforcer,
        cache:    cache,
    }
}

func (m *CasbinMiddleware) CheckPermission() gin.HandlerFunc {
    return func(c *gin.Context) {
        uid := c.GetUint("uid")
        role := c.GetString("role")

        // 1. 检查黑名单（Token 撤销）
        if m.cache.Exists(fmt.Sprintf("token_blacklist:user:%d", uid)) {
            c.JSON(401, gin.H{"msg": "token_revoked"})
            c.Abort()
            return
        }

        // 2. 验证角色仍然激活（缓存 5分钟）
        activeRoles := m.getUserActiveRolesCache(uid)
        if !contains(activeRoles, role) {
            c.JSON(403, gin.H{"msg": "role_revoked"})
            c.Abort()
            return
        }

        // 3. Casbin 权限检查
        ok, err := m.enforcer.Enforce(role, c.Request.URL.Path, c.Request.Method)
        if err != nil {
            c.JSON(500, gin.H{"msg": "permission_check_error"})
            c.Abort()
            return
        }

        if !ok {
            c.JSON(403, gin.H{
                "msg": "forbidden",
                "data": gin.H{
                    "role":   role,
                    "path":   c.Request.URL.Path,
                    "method": c.Request.Method,
                },
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// 获取用户激活角色（带缓存）
func (m *CasbinMiddleware) getUserActiveRolesCache(uid uint) []string {
    cacheKey := fmt.Sprintf("user_roles:%d", uid)

    // 尝试从缓存读取
    cached := m.cache.Get(cacheKey)
    if cached != "" {
        return strings.Split(cached, ",")
    }

    // 从数据库查询
    var roles []model.UserRole
    db.Where("user_id = ? AND status = ?", uid, "active").Find(&roles)

    roleNames := make([]string, len(roles))
    for i, r := range roles {
        roleNames[i] = r.Role
    }

    // 写入缓存（5分钟）
    m.cache.Set(cacheKey, strings.Join(roleNames, ","), 5*time.Minute)

    return roleNames
}
```

### 3.6 Casbin 模型配置

#### configs/casbin_model.conf
```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
```

**说明**：
- 不使用 `[role_definition]`（无角色继承）
- `keyMatch2` 支持路径通配符（如 `/api/b/tasks/*`）
- 每个角色单独定义权限策略

---

## 四、前端改造

### 4.1 类型定义更新

#### src/types/api.ts
```typescript
/**
 * 用户信息
 */
export interface User {
  user_id: number
  role: 'elderly' | 'family' | 'staff' | 'station_manager' | 'admin'  // 当前端的唯一角色
  name: string
  phone: string
  station_id: number | null
}

/**
 * 登录请求
 */
export interface LoginRequest {
  phone: string
  password: string
}

/**
 * 登录响应
 */
export interface LoginResponse {
  token: string
  // 注意：不再返回 refresh_token, roles 等多余字段
  // C端和B端分别登录，各自颁发独立 Token
}
```

### 4.2 API 接口（无需改造）

**说明**：
- C端和B端使用不同的 API Base URL
- 前端配置两套环境变量

#### .env.development
```bash
# C端 API
VITE_C_API_BASE_URL=http://localhost:8080/api/c

# B端 API（管理后台）
VITE_B_API_BASE_URL=http://localhost:8080/api/b

# 通用 API
VITE_COMMON_API_BASE_URL=http://localhost:8080/api/common
```

#### src/utils/request.ts（管理后台）
```typescript
// B端请求实例（管理后台使用）
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_B_API_BASE_URL || 'http://localhost:8080/api/b',
  timeout: 10000,
})

// 请求拦截器（添加 Token）
request.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

// 响应拦截器（处理错误）
request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response) {
      const { status, data } = error.response

      // Token 撤销
      if (data?.msg === 'token_revoked') {
        ElMessage.error('您的登录状态已失效，请重新登录')
        const authStore = useAuthStore()
        authStore.logout()
        router.push('/login')
        return Promise.reject(error)
      }

      // 角色撤销
      if (data?.msg === 'role_revoked') {
        ElMessage.error('您的角色权限已被撤销，请重新登录')
        const authStore = useAuthStore()
        authStore.logout()
        router.push('/login')
        return Promise.reject(error)
      }

      // 其他错误处理...
    }
    return Promise.reject(error)
  }
)
```

### 4.3 状态管理（无需改造）

**说明**：
- 管理后台（B端）和 C端是不同的前端应用
- 各自独立的 Store，无需角色切换逻辑
- 无需 `currentRole`, `switchRole` 等方法

#### src/store/modules/auth.ts（管理后台，保持现状）
```typescript
export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>('')
  const user = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(credentials: LoginRequest) {
    // 调用 B端登录接口
    const response = await authApi.login(credentials)
    const { token: accessToken } = response.data

    setToken(accessToken)

    // 解析 JWT 获取用户信息（或调用 /api/common/user/me）
    const userInfo = parseJWT(accessToken)
    setUser(userInfo)

    return userInfo
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // ... 其他方法保持不变
})
```

### 4.4 路由守卫（无需改造）

**说明**：
- 管理后台路由守卫只检查是否登录和角色权限
- 无需角色切换逻辑

#### src/router/guards/permission.guard.ts（保持现状）
```typescript
export function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // 公开路由
    if (to.meta.public) {
      return next()
    }

    // 登录检查
    if (!authStore.isLoggedIn) {
      ElMessage.warning('请先登录')
      return next({ path: '/login', query: { redirect: to.fullPath } })
    }

    // 角色权限检查
    const requiredRoles = to.meta.roles as string[] | undefined
    if (requiredRoles && requiredRoles.length > 0) {
      const userRole = authStore.user?.role
      if (!userRole || !requiredRoles.includes(userRole)) {
        ElMessage.error('权限不足，无法访问此页面')
        return next(from.path || '/')
      }
    }

    next()
  })
}
```

### 4.5 按钮权限指令（新增）

**功能**：根据角色隐藏/禁用按钮

#### src/directives/permission.ts
```typescript
import { Directive } from 'vue'
import { useAuthStore } from '@/store/modules/auth'

// v-permission="['staff', 'station_manager']"
export const permissionDirective: Directive = {
  mounted(el, binding) {
    const authStore = useAuthStore()
    const requiredRoles = binding.value as string[]

    if (!requiredRoles || requiredRoles.length === 0) {
      return
    }

    const userRole = authStore.user?.role
    if (!userRole || !requiredRoles.includes(userRole)) {
      // 移除元素或禁用
      el.style.display = 'none'
      // 或 el.disabled = true
    }
  },
}

// 注册指令
// main.ts
app.directive('permission', permissionDirective)
```

**使用示例**：
```vue
<template>
  <!-- 只有 station_manager 和 admin 能看到 -->
  <el-button v-permission="['station_manager', 'admin']" @click="handleDelete">
    删除
  </el-button>

  <!-- staff 可以看到 -->
  <el-button v-permission="['staff', 'station_manager', 'admin']" @click="handleClaim">
    认领任务
  </el-button>
</template>
```

---

## 五、测试计划

### 5.1 单元测试

**后端**：
```go
// user_service_test.go
func TestGetUserRoleByType(t *testing.T) {
    // 测试用户在C端的角色
    role, _ := userService.GetUserRoleByType(8, "C端")
    assert.Equal(t, "elderly", role)

    // 测试用户在B端的角色
    role, _ = userService.GetUserRoleByType(4, "B端")
    assert.Equal(t, "staff", role)

    // 测试用户在无权限端
    role, _ = userService.GetUserRoleByType(8, "B端")
    assert.Empty(t, role)
}

// auth_service_test.go
func TestLoginCEnd(t *testing.T) {
    // 测试 elderly 用户登录 C端
    token, err := authService.LoginCEnd("13800000008", "Test@123")
    assert.NoError(t, err)
    assert.NotEmpty(t, token)

    // 测试 staff 用户登录 C端（应该失败）
    token, err = authService.LoginCEnd("13800000004", "Test@123")
    assert.Error(t, err)
    assert.Equal(t, "no_c_end_permission", err.Error())
}

// casbin_test.go
func TestCasbinPermission(t *testing.T) {
    // 测试 staff 权限
    ok, _ := enforcer.Enforce("staff", "/api/b/tasks/pool", "GET")
    assert.True(t, ok)

    ok, _ = enforcer.Enforce("staff", "/api/b/users", "POST")
    assert.False(t, ok) // staff 无权限

    // 测试 admin 权限
    ok, _ = enforcer.Enforce("admin", "/api/b/users", "POST")
    assert.True(t, ok)
}
```

### 5.2 集成测试

**测试场景**：
1. **跨端登录**：
   - 用户A（elderly + staff）分别在 C端和B端登录
   - 验证两个 Token 独立且有效

2. **权限隔离**：
   - C端 Token 访问 B端接口 → 403
   - B端 Token 访问 C端接口 → 403

3. **角色撤销实时性**：
   - 撤销用户 staff 角色
   - 1分钟内该用户的 B端 Token 失效

4. **黑名单生效**：
   - 用户加入黑名单
   - 所有端的 Token 立即失效

### 5.3 端到端测试

**C端测试**：
```
1. elderly 用户登录 C端
2. 提交服务需求
3. 查看我的需求列表
4. 查看个人档案

5. family 用户登录 C端
6. 只能查看个人信息（无法提交需求）
```

**B端测试**：
```
1. staff 用户登录管理后台
2. 查看任务池
3. 认领任务
4. 完成任务（上传图片）
5. 查看我的任务

6. station_manager 用户登录
7. 访问 Dashboard
8. 查看站点统计
9. 查看所有任务（含已认领）

10. admin 用户登录
11. 访问用户管理
12. 创建/编辑用户
13. 分配/撤销角色
14. 查看权限配置
```

**角色撤销测试**：
```
1. staff 用户登录管理后台（获取 Token）
2. admin 撤销该用户的 staff 角色
3. 等待 30 秒
4. staff 用户刷新页面或发起新请求
5. 验证：被踢出登录，提示"角色已被撤销"
```

---

## 六、数据迁移计划

### 6.1 现有数据处理

**已完成**：
- ✅ user_roles 表已创建
- ✅ 现有用户数据已迁移到 user_roles
- ✅ users.primary_role 已填充

**待执行**：
1. **Casbin 策略迁移**：
```sql
-- 从 CSV 迁移到数据库表
-- 执行脚本：backend/scripts/migrate_casbin_to_db.go
```

2. **验证数据一致性**：
```sql
-- 检查是否有用户在同一端有多个角色（应该为空）
SELECT user_id, GROUP_CONCAT(role) as roles, COUNT(*) as cnt
FROM user_roles
WHERE status = 'active'
  AND role IN ('staff', 'station_manager', 'admin')
GROUP BY user_id
HAVING cnt > 1;

-- 检查是否有用户没有角色（应该为空）
SELECT id, phone, name
FROM users
WHERE id NOT IN (SELECT DISTINCT user_id FROM user_roles WHERE status = 'active');
```

### 6.2 回滚方案

**紧急回滚步骤**：
```sql
-- 1. 恢复 users.role 字段
UPDATE users u
SET role = (
  SELECT role
  FROM user_roles
  WHERE user_id = u.id AND is_primary = TRUE
  LIMIT 1
)
WHERE role IS NULL OR role = '';

-- 2. Casbin 策略回滚到 CSV
-- 导出数据库策略到 CSV 文件
-- 修改代码使用 CSV Adapter

-- 3. 代码回滚
git revert <commit-hash>
```

---

## 七、部署与监控

### 7.1 部署步骤

**灰度发布**：
1. **阶段1（10% 流量）**：
   - 部署新后端（兼容新旧 JWT 格式）
   - 前端保持不变
   - 观察 24 小时

2. **阶段2（50% 流量）**：
   - 增加灰度比例
   - 监控错误率和性能指标
   - 观察 48 小时

3. **阶段3（100% 流量）**：
   - 全量发布
   - 持续监控 1 周

### 7.2 监控指标

**关键指标**：
```yaml
业务指标:
  - login_success_rate: 登录成功率 (> 99%)
  - login_c_end_count: C端登录次数
  - login_b_end_count: B端登录次数
  - cross_end_login_count: 跨端登录用户数

性能指标:
  - jwt_validation_latency: JWT验证延迟 (< 5ms)
  - casbin_check_latency: 权限检查延迟 (< 10ms)
  - redis_blacklist_latency: 黑名单查询延迟 (< 2ms)

错误指标:
  - jwt_decode_error_rate: JWT解析错误率 (< 0.01%)
  - permission_denied_rate: 权限拒绝率
  - token_revoked_count: Token撤销次数
  - role_revoked_count: 角色撤销次数
```

**告警规则**：
```yaml
告警级别P0（立即处理）:
  - login_success_rate < 95%: 登录成功率过低
  - permission_denied_rate > 10%: 权限拒绝率异常

告警级别P1（24h内处理）:
  - casbin_check_latency > 50ms: 权限检查过慢
  - jwt_decode_error_rate > 0.1%: Token解析异常
```

### 7.3 日志规范

**关键日志**：
```go
// 登录日志
log.Info("user_login",
    zap.String("端", "C端"),
    zap.Uint("user_id", uid),
    zap.String("role", role),
    zap.String("ip", clientIP))

// 权限拒绝日志
log.Warn("permission_denied",
    zap.Uint("user_id", uid),
    zap.String("role", role),
    zap.String("path", path),
    zap.String("method", method))

// 角色撤销日志
log.Info("role_revoked",
    zap.Uint("user_id", uid),
    zap.String("role", role),
    zap.Uint("operator_id", operatorID))
```

---

## 八、风险评估与缓解

### 8.1 已识别风险

| 风险项 | 影响 | 概率 | 缓解措施 |
|--------|------|------|----------|
| JWT 格式不兼容导致旧客户端无法登录 | 高 | 中 | 后端同时支持新旧格式，灰度发布 |
| Casbin 策略迁移丢失权限 | 高 | 低 | 迁移前备份 CSV，验证数据完整性 |
| 角色撤销延迟导致权限泄露 | 中 | 低 | Redis 缓存 TTL 设为 5分钟，黑名单立即生效 |
| 跨端登录导致 Token 冲突 | 低 | 低 | C端和B端 Token 独立存储（localStorage 不同 key） |
| 权限检查性能下降 | 中 | 中 | Redis 缓存 + Casbin 内存策略，监控延迟 |

### 8.2 应急预案

**场景1：权限检查服务异常**
```
现象：大量 403 错误
处理：
  1. 检查 Casbin enforcer 加载状态
  2. 检查 Redis 连接
  3. 临时关闭权限检查（紧急开关）
  4. 修复后重新加载策略
```

**场景2：角色撤销不生效**
```
现象：用户被撤销角色后仍能访问
处理：
  1. 检查 Redis 黑名单是否生效
  2. 检查缓存 TTL 设置
  3. 手动清除该用户的所有缓存
  4. 强制该用户重新登录
```

---

## 九、实施时间线

### Phase 1: 后端核心（3-4天）
- **Day 1**：
  - [ ] 数据模型改造（User, UserRole）
  - [ ] 共享服务实现（UserService, AuthService）
  - [ ] 单元测试

- **Day 2**：
  - [ ] 登录接口改造（C端、B端）
  - [ ] JWT 生成与验证
  - [ ] 集成测试

- **Day 3**：
  - [ ] Casbin 策略迁移到数据库
  - [ ] 权限中间件改造
  - [ ] Redis 缓存集成

- **Day 4**：
  - [ ] 端到端测试
  - [ ] 性能测试与优化

### Phase 2: 前端适配（1-2天）
- **Day 5**：
  - [ ] 类型定义更新
  - [ ] 登录流程测试（C端、B端）
  - [ ] 错误处理适配（token_revoked, role_revoked）

- **Day 6**：
  - [ ] 按钮权限指令实现
  - [ ] 前端集成测试

### Phase 3: 部署与监控（2-3天）
- **Day 7**：
  - [ ] 灰度发布 10%
  - [ ] 监控指标观察

- **Day 8-9**：
  - [ ] 灰度 50% → 100%
  - [ ] 性能调优
  - [ ] 文档更新

**总计**：7-9 个工作日

---

## 十、后续优化（P1 功能）

### 10.1 用户管理界面
- [ ] 用户列表显示多角色
- [ ] 编辑用户时支持分配/撤销角色
- [ ] 角色变更日志

### 10.2 权限管理界面
- [ ] Casbin 策略可视化管理
- [ ] 角色权限编辑（CRUD）
- [ ] 权限测试工具（输入角色+路径，输出是否允许）

### 10.3 审计日志
- [ ] 角色变更记录
- [ ] 权限拒绝日志聚合
- [ ] 异常登录告警

---

## Spec 已锁定 ✅

**核心设计确认**：
- ✅ 端内角色唯一，跨端多角色
- ✅ 不同登录接口（C端、B端）
- ✅ JWT 单角色 Token
- ✅ Casbin RBAC 无继承
- ✅ 实时角色撤销（JWT + Redis 黑名单）
- ✅ 共享模块架构（common/user, common/auth）

**可以进入执行模式** 🚀

请确认是否开始实施，或有其他调整？
