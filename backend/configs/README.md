# Casbin 权限配置说明

## 文件说明

### 1. casbin_model.conf
RBAC权限模型定义文件，定义了：
- **request_definition**: 请求格式(subject, object, action)
- **policy_definition**: 策略格式
- **role_definition**: 角色继承关系
- **policy_effect**: 策略效果(允许/拒绝)
- **matchers**: 匹配规则

### 2. policy.csv
权限策略数据文件，定义了：
- 每个角色可以访问哪些API
- API的HTTP方法权限
- 资源访问控制

## 角色权限矩阵

| 角色 | 中文名 | 权限范围 |
|------|--------|---------|
| `elderly` | 老年人 | 提交需求、查看自己的需求、接收通知 |
| `family` | 家属 | 代老年人提交需求、查看需求、接收通知 |
| `staff` | 工作人员 | 查看任务池、认领任务、完成任务、查看站点需求 |
| `station_manager` | 站点负责人 | staff权限 + 任务转派 + 站点统计 |
| `admin` | 系统管理员 | 所有权限 + 用户管理 + 站点管理 + 围栏管理 |

## 权限控制实现

### 1. 中间件集成

```go
// backend/internal/middleware/casbin.go
func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从JWT Token获取用户角色
        role := c.GetString("user_role")

        // 获取请求资源和方法
        obj := c.Request.URL.Path
        act := c.Request.Method

        // 执行权限检查
        ok, err := enforcer.Enforce(role, obj, act)
        if err != nil {
            c.JSON(500, gin.H{"error": "权限检查失败"})
            c.Abort()
            return
        }

        if !ok {
            c.JSON(403, gin.H{"error": "权限不足"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 2. 路由注册

```go
// backend/internal/handler/routes.go
func RegisterRoutes(r *gin.Engine, enforcer *casbin.Enforcer) {
    api := r.Group("/api")

    // 公开路由(不需要认证)
    api.POST("/auth/login", authHandler.Login)

    // 需要认证的路由
    auth := api.Group("")
    auth.Use(middleware.JWTAuth())
    auth.Use(middleware.CasbinMiddleware(enforcer))
    {
        // 需求管理
        auth.POST("/requests", requestHandler.Create)
        auth.GET("/requests/:id", requestHandler.Get)
        auth.GET("/requests", requestHandler.List)

        // 任务管理
        auth.GET("/tasks", taskHandler.List)
        auth.POST("/tasks/:id/claim", taskHandler.Claim)
        auth.POST("/tasks/:id/complete", taskHandler.Complete)

        // ... 其他路由
    }
}
```

### 3. 动态路径处理

对于`:id`这样的动态路径参数，Casbin需要特殊处理：

**方案A**: 使用通配符匹配
```csv
p, staff, /api/tasks/*, GET
```

**方案B**: 在中间件中规范化路径
```go
// 将 /api/tasks/123 转换为 /api/tasks/:id
func normalizePath(path string) string {
    // 使用正则或路由匹配转换
}
```

**MVP推荐**: 方案A(简单直接)

## 权限策略更新

### 添加新API权限

1. 编辑 `policy.csv`
2. 添加新的权限行
3. 重启服务(或热加载)

**示例**: 添加新的评价接口
```csv
p, elderly, /api/requests/:id/rate, POST
p, family, /api/requests/:id/rate, POST
```

### 添加新角色

1. 在 `policy.csv` 添加角色权限
2. 在 `domain/roles.go` 添加角色常量
3. 更新权限矩阵文档

## 测试验证

### 1. 单元测试

```go
func TestCasbinPermissions(t *testing.T) {
    e, _ := casbin.NewEnforcer("casbin_model.conf", "policy.csv")

    // 测试工作人员可以认领任务
    ok, _ := e.Enforce("staff", "/api/tasks/:id/claim", "POST")
    assert.True(t, ok)

    // 测试老年人不能认领任务
    ok, _ = e.Enforce("elderly", "/api/tasks/:id/claim", "POST")
    assert.False(t, ok)
}
```

### 2. 集成测试

```bash
# 使用staff角色Token访问任务认领接口
curl -H "Authorization: Bearer <staff_token>" \
     -X POST http://localhost:8080/api/tasks/1/claim

# 期望: 200 OK

# 使用elderly角色Token访问任务认领接口
curl -H "Authorization: Bearer <elderly_token>" \
     -X POST http://localhost:8080/api/tasks/1/claim

# 期望: 403 Forbidden
```

## 常见问题

### Q1: 如何处理资源级权限(如只能查看自己的需求)?

**A**: Casbin策略控制路由级权限，资源级权限在业务逻辑中实现：

```go
func (h *RequestHandler) Get(c *gin.Context) {
    requestID := c.Param("id")
    userID := c.GetInt64("user_id")
    userRole := c.GetString("user_role")

    request := // 查询需求

    // 资源级权限检查
    if userRole == "elderly" || userRole == "family" {
        if request.UserID != userID {
            c.JSON(403, gin.H{"error": "只能查看自己的需求"})
            return
        }
    }

    // 返回数据
}
```

### Q2: 权限策略太多，如何简化?

**A**: 启用角色继承关系：

```csv
# 在 policy.csv 添加
g, station_manager, staff
g, admin, station_manager

# 这样只需定义 staff 权限，station_manager 和 admin 自动继承
```

### Q3: 如何动态添加权限(不重启服务)?

**A**: 使用Casbin的API动态添加：

```go
// 添加新策略
enforcer.AddPolicy("staff", "/api/new-api", "GET")

// 保存到文件
enforcer.SavePolicy()
```

## 参考资料

- Casbin官方文档: https://casbin.org/
- RBAC模型说明: https://casbin.org/docs/rbac
- API接口文档: `/docs/API_PLANNING.md`
