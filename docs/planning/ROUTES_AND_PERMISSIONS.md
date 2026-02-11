# sCare 系统路由资源清单与权限配置

## 一、白名单路由（无需认证）

| 路由 | 方法 | 说明 | 是否公开 |
|------|------|------|---------|
| `/api/v1/health` | GET | 健康检查 | ✅ 公开 |
| `/static/*` | GET | 静态文件访问 | ✅ 公开 |
| `/api/v1/b/auth/login` | POST | B端登录 | ✅ 公开 |
| `/api/v1/b/auth/refresh` | POST | B端刷新Token | ✅ 公开 |
| `/api/v1/c/auth/login` | POST | C端登录 | ✅ 公开 |
| `/api/v1/c/auth/refresh` | POST | C端刷新Token | ✅ 公开 |

## 二、需要认证但无需权限检查的路由

| 路由 | 方法 | 说明 | 需要Token |
|------|------|------|-----------|
| `/api/v1/b/auth/me` | GET | 获取B端用户信息 | ✅ 需要 |
| `/api/v1/c/auth/me` | GET | 获取C端用户信息 | ✅ 需要 |

## 三、受保护资源（需要认证 + 权限检查）

### 3.1 B端资源

#### 服务需求管理
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/b/requests` | GET | 查看需求列表 | staff, station_manager, admin |
| `/api/v1/b/requests/:id` | GET | 查看需求详情 | staff, station_manager, admin |

#### 任务管理
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/b/tasks/pool` | GET | 查看任务池 | staff, station_manager, admin |
| `/api/v1/b/tasks/my` | GET | 查看我的任务 | staff, station_manager, admin |
| `/api/v1/b/tasks/:id/claim` | POST | 认领任务 | staff, station_manager, admin |
| `/api/v1/b/tasks/:id/complete` | POST | 完成任务 | staff, station_manager, admin |

#### 服务站点管理
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/b/stations` | GET | 查看站点列表 | staff, station_manager, admin |
| `/api/v1/b/stations` | POST | 创建站点 | admin |
| `/api/v1/b/stations/:id` | GET | 查看站点详情 | staff, station_manager, admin |
| `/api/v1/b/stations/:id` | PUT | 更新站点 | admin |
| `/api/v1/b/stations/:id` | DELETE | 删除站点 | admin |

#### 服务围栏管理
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/b/zones` | GET | 查看围栏列表 | staff, station_manager, admin |
| `/api/v1/b/zones` | POST | 创建围栏 | station_manager, admin |
| `/api/v1/b/zones/:id` | PUT | 更新围栏 | station_manager, admin |
| `/api/v1/b/zones/:id` | DELETE | 删除围栏 | station_manager, admin |

#### 用户管理
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/b/users` | GET | 查看用户列表 | station_manager, admin |
| `/api/v1/b/users` | POST | 创建用户 | admin |
| `/api/v1/b/users/:id` | GET | 查看用户详情 | station_manager, admin |
| `/api/v1/b/users/:id` | PUT | 更新用户 | admin |

#### 通知管理
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/b/notifications` | GET | 查看通知列表 | staff, station_manager, admin |
| `/api/v1/b/notifications/:id/read` | POST | 标记已读 | staff, station_manager, admin |

#### 文件上传
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/b/upload` | POST | 上传文件 | staff, station_manager, admin |

### 3.2 C端资源

#### 服务需求管理
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/c/requests` | POST | 创建服务需求 | elderly, family |
| `/api/v1/c/requests` | GET | 查看我的需求列表 | elderly, family |
| `/api/v1/c/requests/:id` | GET | 查看需求详情 | elderly, family |
| `/api/v1/c/requests/:id/cancel` | POST | 取消需求 | elderly, family |

#### 通知管理
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/c/notifications` | GET | 查看通知列表 | elderly, family |
| `/api/v1/c/notifications/:id/read` | POST | 标记已读 | elderly, family |

#### 文件上传
| 路由 | 方法 | 说明 | 需要角色 |
|------|------|------|----------|
| `/api/v1/c/upload` | POST | 上传文件 | elderly, family |

## 四、角色权限矩阵

### B端角色

| 角色 | 权限范围 | 说明 |
|------|---------|------|
| **admin** | 全部B端资源 | 超级管理员，拥有所有权限 |
| **station_manager** | 任务管理、需求查看、站点查看、围栏管理、用户查看、通知、上传 | 站长，管理本站点相关事务 |
| **staff** | 任务管理、需求查看、通知、上传 | 工作人员，执行服务任务 |

### C端角色

| 角色 | 权限范围 | 说明 |
|------|---------|------|
| **elderly** | 创建需求、查看需求、取消需求、通知、上传 | 老年人，提交服务需求 |
| **family** | 查看需求、通知 | 家属，查看关联老人的需求 |

## 五、Casbin 策略配置

### 策略格式
```
p, <主体>, <资源>, <操作>
```

### 完整策略见 database/seeds/seed.sql
- B端策略：`role:admin`, `role:station_manager`, `role:staff`
- C端策略：`role:elderly`, `role:family`
