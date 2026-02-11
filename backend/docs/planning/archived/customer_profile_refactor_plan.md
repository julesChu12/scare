# 客户档案架构重构实施计划（简化版）

## 目标
将 C端用户从 `user_roles` 表中独立出来，使用 `customer_profiles` 表管理服务对象档案。

## 核心变更

### 概念变化
- **之前**：elderly/family 作为角色存储在 `user_roles` 表，type='c_end'
- **之后**：
  - 去掉 family 概念
  - 服务对象（老幼病残孕等）通过 `customer_profiles` 表标识
  - `user_roles` 表只用于 B端权限控制，删除 type 字段

### 术语统一
- **customer**（客户/服务对象）：需要照护的人群（老年人、失能人士、孕妇、儿童等）
- **B端用户**：管理员、站长、工作人员（有权限控制）
- **C端用户**：客户本人或其家属（通过 customer_profile 识别）

### 用户身份判定
```
B端用户: 在 user_roles 表中有记录
C端用户: 在 customer_profiles 表中有记录
跨端用户: 同时在两张表中有记录
```

---

## 实施步骤

### 步骤1: 清空相关表并修改结构

#### 1.1 清空数据库
```sql
SET FOREIGN_KEY_CHECKS = 0;

-- 清空所有相关表
TRUNCATE TABLE user_roles;
TRUNCATE TABLE users;
TRUNCATE TABLE casbin_rule;
-- 清空其他业务表...

SET FOREIGN_KEY_CHECKS = 1;
```

#### 1.2 修改 user_roles 表（删除 type 字段）
```sql
-- 删除 type 字段（只用于B端，不再需要区分）
ALTER TABLE user_roles DROP COLUMN IF EXISTS type;

-- 确认表结构
DESCRIBE user_roles;
```

---

### 步骤2: 创建 customer_profiles 表

```sql
-- 客户档案表（服务对象：老幼病残孕等需要照护的人群）
CREATE TABLE `customer_profiles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '关联用户ID',

  -- 基本信息
  `id_card` VARCHAR(18) DEFAULT NULL COMMENT '身份证号',
  `gender` VARCHAR(10) DEFAULT NULL COMMENT '性别',
  `birth_date` DATE DEFAULT NULL COMMENT '出生日期',
  `address` TEXT DEFAULT NULL COMMENT '居住地址',

  -- 客户类型和健康信息
  `customer_type` VARCHAR(20) DEFAULT NULL COMMENT '客户类型：elderly/disabled/pregnant/child/other',
  `health_status` VARCHAR(50) DEFAULT NULL COMMENT '健康状况',
  `disability_level` VARCHAR(20) DEFAULT NULL COMMENT '失能等级：自理/轻度/中度/重度',
  `medical_history` TEXT DEFAULT NULL COMMENT '病史',
  `special_needs` TEXT DEFAULT NULL COMMENT '特殊需求',

  -- 紧急联系人（JSON格式）
  `emergency_contact` JSON DEFAULT NULL COMMENT '紧急联系人 {"name":"","phone":"","relation":""}',

  -- 时间戳
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  KEY `idx_id_card` (`id_card`),
  KEY `idx_customer_type` (`customer_type`),
  CONSTRAINT `fk_customer_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户档案表（服务对象）';
```

**字段说明**：
- `customer_type`：客户类型
  - `elderly` - 老年人
  - `disabled` - 失能人士
  - `pregnant` - 孕妇
  - `child` - 儿童
  - `other` - 其他需要照护的人群

---

### 步骤3: 创建领域模型

#### 3.1 创建 `internal/domain/customer_profile.go`
```go
package domain

import "time"

// CustomerProfile 客户档案（服务对象）
type CustomerProfile struct {
    ID     int64 `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID int64 `gorm:"uniqueIndex;not null" json:"user_id"`

    // 基本信息
    IDCard    *string `gorm:"type:varchar(18)" json:"id_card,omitempty"`
    Gender    *string `gorm:"type:varchar(10)" json:"gender,omitempty"`
    BirthDate *string `gorm:"type:date" json:"birth_date,omitempty"`
    Address   *string `gorm:"type:text" json:"address,omitempty"`

    // 客户类型和健康信息
    CustomerType    *string `gorm:"type:varchar(20)" json:"customer_type,omitempty"` // elderly/disabled/pregnant/child/other
    HealthStatus    *string `gorm:"type:varchar(50)" json:"health_status,omitempty"`
    DisabilityLevel *string `gorm:"type:varchar(20)" json:"disability_level,omitempty"`
    MedicalHistory  *string `gorm:"type:text" json:"medical_history,omitempty"`
    SpecialNeeds    *string `gorm:"type:text" json:"special_needs,omitempty"`

    // 紧急联系人
    EmergencyContact *EmergencyContact `gorm:"type:json" json:"emergency_contact,omitempty"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    // 关联
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// EmergencyContact 紧急联系人结构
type EmergencyContact struct {
    Name     string `json:"name"`
    Phone    string `json:"phone"`
    Relation string `json:"relation"`
}

func (CustomerProfile) TableName() string {
    return "customer_profiles"
}

// CustomerType 常量
const (
    CustomerTypeElderly  = "elderly"  // 老年人
    CustomerTypeDisabled = "disabled" // 失能人士
    CustomerTypePregnant = "pregnant" // 孕妇
    CustomerTypeChild    = "child"    // 儿童
    CustomerTypeOther    = "other"    // 其他
)
```

#### 3.2 创建 `internal/repository/customer_repository.go`
```go
package repository

import (
    "community-elderly-care-platform/internal/domain"
    "gorm.io/gorm"
)

type CustomerRepository struct {
    db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
    return &CustomerRepository{db: db}
}

// GetByUserID 根据用户ID获取档案
func (r *CustomerRepository) GetByUserID(userID int64) (*domain.CustomerProfile, error) {
    var profile domain.CustomerProfile
    err := r.db.Where("user_id = ?", userID).First(&profile).Error
    if err != nil {
        return nil, err
    }
    return &profile, nil
}

// GetByUserIDWithUser 获取档案并预加载用户信息
func (r *CustomerRepository) GetByUserIDWithUser(userID int64) (*domain.CustomerProfile, error) {
    var profile domain.CustomerProfile
    err := r.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error
    if err != nil {
        return nil, err
    }
    return &profile, nil
}

// Create 创建档案
func (r *CustomerRepository) Create(profile *domain.CustomerProfile) error {
    return r.db.Create(profile).Error
}

// Update 更新档案
func (r *CustomerRepository) Update(profile *domain.CustomerProfile) error {
    return r.db.Save(profile).Error
}

// Delete 删除档案
func (r *CustomerRepository) Delete(userID int64) error {
    return r.db.Where("user_id = ?", userID).Delete(&domain.CustomerProfile{}).Error
}

// Exists 检查用户是否有档案
func (r *CustomerRepository) Exists(userID int64) (bool, error) {
    var count int64
    err := r.db.Model(&domain.CustomerProfile{}).Where("user_id = ?", userID).Count(&count).Error
    return count > 0, err
}

// ListByType 根据客户类型查询
func (r *CustomerRepository) ListByType(customerType string) ([]*domain.CustomerProfile, error) {
    var profiles []*domain.CustomerProfile
    err := r.db.Where("customer_type = ?", customerType).Find(&profiles).Error
    return profiles, err
}
```

---

### 步骤4: 修改 JWT Manager

#### 4.1 修改 `pkg/jwt/manager.go`

```go
// Claims JWT 载荷
type Claims struct {
    UserID    int64    `json:"uid"`
    Type      string   `json:"type"`  // "b_end" 或 "c_end"
    StationID int64    `json:"station_id,omitempty"`

    // 仅 B端使用
    Roles     []string `json:"roles,omitempty"`

    jwtlib.RegisteredClaims
}

// GenerateToken 生成访问令牌
// endType: "b_end" 或 "c_end"
// roles: 仅 B端需要，C端传 nil 或空数组
func (m *Manager) GenerateToken(userID int64, endType string, stationID int64, roles []string) (string, error) {
    jti := uuid.New().String()

    claims := Claims{
        UserID:    userID,
        Type:      endType,
        StationID: stationID,
        Roles:     roles,  // C端时为空
        RegisteredClaims: jwtlib.RegisteredClaims{
            ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.expiresIn)),
            IssuedAt:  jwtlib.NewNumericDate(time.Now()),
            ID:        jti,
        },
    }

    token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
    return token.SignedString([]byte(m.secret))
}

// GenerateRefreshToken 生成刷新令牌
func (m *Manager) GenerateRefreshToken(userID int64, endType string, stationID int64, roles []string) (string, error) {
    jti := uuid.New().String()

    claims := Claims{
        UserID:    userID,
        Type:      endType,
        StationID: stationID,
        Roles:     roles,  // C端时为空
        RegisteredClaims: jwtlib.RegisteredClaims{
            ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.refreshExpiresIn)),
            IssuedAt:  jwtlib.NewNumericDate(time.Now()),
            ID:        jti,
        },
    }

    token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
    return token.SignedString([]byte(m.secret))
}
```

---

### 步骤5: 重写认证服务

#### 5.1 修改 `internal/service/auth_service.go`

```go
// 新增错误类型
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserInactive       = errors.New("user inactive")
    ErrNoRoleForEndType   = errors.New("user has no role for this end type")
    ErrNoCustomerProfile  = errors.New("user has no customer profile")  // 新增
)

// ValidateCredentials 验证用户名密码（B端和C端共用）
func (s *AuthService) ValidateCredentials(phone, password string) (*domain.User, error) {
    user, err := s.userRepo.GetByPhone(phone)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    if !s.verifyPassword(user.PasswordHash, password) {
        return nil, ErrInvalidCredentials
    }

    if user.Status != "active" {
        return nil, ErrUserInactive
    }

    return user, nil
}

// LoginBEnd B端登录
func (s *AuthService) LoginBEnd(phone, password string) (*Tokens, *domain.User, []string, error) {
    user, err := s.ValidateCredentials(phone, password)
    if err != nil {
        return nil, nil, nil, err
    }

    // 获取B端角色
    userRoles, err := s.userRepo.GetActiveRoles(user.ID)
    if err != nil || len(userRoles) == 0 {
        return nil, nil, nil, ErrNoRoleForEndType
    }

    stationID := int64(0)
    if user.StationID != nil {
        stationID = *user.StationID
    }

    // 生成Token（包含角色）
    accessToken, err := s.jwtManager.GenerateToken(user.ID, "b_end", stationID, userRoles)
    if err != nil {
        return nil, nil, nil, err
    }

    refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, "b_end", stationID, userRoles)
    if err != nil {
        return nil, nil, nil, err
    }

    return &Tokens{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, user, userRoles, nil
}

// LoginCEnd C端登录
func (s *AuthService) LoginCEnd(phone, password string, customerRepo *repository.CustomerRepository) (*Tokens, *domain.User, error) {
    user, err := s.ValidateCredentials(phone, password)
    if err != nil {
        return nil, nil, err
    }

    // 检查是否有客户档案
    exists, err := customerRepo.Exists(user.ID)
    if err != nil || !exists {
        return nil, nil, ErrNoCustomerProfile
    }

    // 生成Token（不包含角色）
    accessToken, err := s.jwtManager.GenerateToken(user.ID, "c_end", 0, nil)
    if err != nil {
        return nil, nil, err
    }

    refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, "c_end", 0, nil)
    if err != nil {
        return nil, nil, err
    }

    return &Tokens{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, user, nil
}
```

---

### 步骤6: 重写 C端认证处理器

#### 6.1 完全重写 `internal/handler/c_auth_handler.go`

```go
package handler

import (
    "net/http"

    "community-elderly-care-platform/internal/repository"
    "community-elderly-care-platform/internal/service"

    "github.com/gin-gonic/gin"
)

// CAuthHandler C端认证处理器
type CAuthHandler struct {
    authService  *service.AuthService
    userRepo     *repository.UserRepository
    customerRepo *repository.CustomerRepository
}

func NewCAuthHandler(
    authService *service.AuthService,
    userRepo *repository.UserRepository,
    customerRepo *repository.CustomerRepository,
) *CAuthHandler {
    return &CAuthHandler{
        authService:  authService,
        userRepo:     userRepo,
        customerRepo: customerRepo,
    }
}

type cLoginRequest struct {
    Phone    string `json:"phone" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type cRefreshRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login C端登录接口
// @route POST /api/v1/c/auth/login
func (h *CAuthHandler) Login(c *gin.Context) {
    var req cLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        RespondError(c, http.StatusBadRequest, "invalid payload")
        return
    }

    // 使用新的 LoginCEnd 方法
    tokens, user, err := h.authService.LoginCEnd(req.Phone, req.Password, h.customerRepo)
    if err != nil {
        if err == service.ErrInvalidCredentials {
            RespondError(c, http.StatusUnauthorized, "invalid credentials")
            return
        }
        if err == service.ErrUserInactive {
            RespondError(c, http.StatusForbidden, "user inactive")
            return
        }
        if err == service.ErrNoCustomerProfile {
            RespondError(c, http.StatusForbidden, "user not registered as customer")
            return
        }
        RespondError(c, http.StatusInternalServerError, "login failed")
        return
    }

    data := gin.H{
        "token":         tokens.AccessToken,
        "refresh_token": tokens.RefreshToken,
        "user_id":       user.ID,
        "type":          "c_end",
        "name":          user.Name,
        "phone":         user.Phone,
        "status":        user.Status,
    }
    Respond(c, http.StatusOK, "ok", data)
}

// Refresh C端刷新Token
func (h *CAuthHandler) Refresh(c *gin.Context) {
    var req cRefreshRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        RespondError(c, http.StatusBadRequest, "invalid payload")
        return
    }

    tokens, err := h.authService.Refresh(req.RefreshToken)
    if err != nil {
        RespondError(c, http.StatusUnauthorized, "invalid token")
        return
    }

    Respond(c, http.StatusOK, "ok", gin.H{
        "token":         tokens.AccessToken,
        "refresh_token": tokens.RefreshToken,
    })
}

// Me 获取当前登录用户信息（C端）
func (h *CAuthHandler) Me(c *gin.Context) {
    userID, ok := GetUserID(c)
    if !ok {
        RespondError(c, http.StatusUnauthorized, "missing user")
        return
    }

    // 获取用户基本信息
    user, err := h.userRepo.GetByID(userID)
    if err != nil {
        RespondError(c, http.StatusNotFound, "user not found")
        return
    }

    // 获取客户档案
    profile, err := h.customerRepo.GetByUserID(userID)
    if err != nil {
        RespondError(c, http.StatusForbidden, "user has no customer profile")
        return
    }

    Respond(c, http.StatusOK, "ok", gin.H{
        "user": gin.H{
            "id":     user.ID,
            "name":   user.Name,
            "phone":  user.Phone,
            "email":  user.Email,
            "avatar": user.Avatar,
            "status": user.Status,
        },
        "profile": gin.H{
            "customer_type":     profile.CustomerType,
            "gender":            profile.Gender,
            "birth_date":        profile.BirthDate,
            "address":           profile.Address,
            "health_status":     profile.HealthStatus,
            "disability_level":  profile.DisabilityLevel,
            "emergency_contact": profile.EmergencyContact,
        },
    })
}
```

---

### 步骤7: 修改 B端认证处理器

#### 7.1 修改 `internal/handler/b_auth_handler.go`

```go
// Login B端登录接口
func (h *BAuthHandler) Login(c *gin.Context) {
    var req bLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        RespondError(c, http.StatusBadRequest, "invalid payload")
        return
    }

    // 使用新的 LoginBEnd 方法
    tokens, user, roles, err := h.authService.LoginBEnd(req.Phone, req.Password)
    if err != nil {
        if err == service.ErrInvalidCredentials {
            RespondError(c, http.StatusUnauthorized, "invalid credentials")
            return
        }
        if err == service.ErrUserInactive {
            RespondError(c, http.StatusForbidden, "user inactive")
            return
        }
        if err == service.ErrNoRoleForEndType {
            RespondError(c, http.StatusForbidden, "user has no role for B-end")
            return
        }
        RespondError(c, http.StatusInternalServerError, "login failed")
        return
    }

    data := gin.H{
        "token":         tokens.AccessToken,
        "refresh_token": tokens.RefreshToken,
        "user_id":       user.ID,
        "roles":         roles,
        "type":          "b_end",
        "station_id":    user.StationID,
        "name":          user.Name,
        "phone":         user.Phone,
        "status":        user.Status,
    }
    Respond(c, http.StatusOK, "ok", data)
}
```

---

### 步骤8: 修改路由注册

#### 8.1 修改 `internal/handler/routes.go`

```go
func RegisterRoutes(router *gin.Engine, db *database.DB, rdb *redis.Client, cfg *config.Config) error {
    // ... 前面代码保持不变

    // Repository 初始化
    userRepo := repository.NewUserRepository(db.DB)
    customerRepo := repository.NewCustomerRepository(db.DB)  // 新增
    stationRepo := repository.NewStationRepository(db.DB)
    // ... 其他 repo

    // Service 初始化
    jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.JWT.RefreshExpiresIn)
    authService := service.NewAuthService(userRepo, jwtManager)
    // ... 其他 service

    // Handler 初始化
    bAuthHandler := NewBAuthHandler(authService, userRepo, permissionService)
    cAuthHandler := NewCAuthHandler(authService, userRepo, customerRepo)  // 修改参数

    // ... 路由注册保持不变
}
```

---

### 步骤9: 创建新的测试数据

#### 9.1 修改 `database/seeds/seed.sql`

```sql
-- =====================================================
-- 1. 清空所有表（重新开始）
-- =====================================================
SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE user_roles;
TRUNCATE TABLE customer_profiles;
TRUNCATE TABLE users;
TRUNCATE TABLE casbin_rule;
TRUNCATE TABLE stations;
TRUNCATE TABLE zones;
TRUNCATE TABLE requests;
TRUNCATE TABLE tasks;
TRUNCATE TABLE notifications;

SET FOREIGN_KEY_CHECKS = 1;

-- =====================================================
-- 2. 初始化服务站点
-- =====================================================
INSERT INTO `stations` (`id`, `name`, `address`, `phone`, `manager_id`, `status`) VALUES
(1, '朝阳区幸福街道服务站', '北京市朝阳区幸福大街100号', '010-12345678', 2, 'active'),
(2, '朝阳区康乐街道服务站', '北京市朝阳区康乐大街200号', '010-87654321', 3, 'active');

-- =====================================================
-- 3. 初始化用户
-- =====================================================
INSERT INTO `users` (`id`, `phone`, `password_hash`, `name`, `email`, `gender`, `birth_date`, `station_id`, `status`) VALUES
-- B端用户
(1, '13800000001', '$2a$10$HASH_FOR_Test@123', '系统管理员', 'admin@scare.com', 'male', '1980-01-01', NULL, 'active'),
(2, '13800000002', '$2a$10$HASH_FOR_Test@123', '李站长', 'zhangzhang@scare.com', 'male', '1975-05-15', 1, 'active'),
(3, '13800000003', '$2a$10$HASH_FOR_Test@123', '王站长', 'wangzhang@scare.com', 'female', '1978-08-20', 2, 'active'),
(4, '13800000004', '$2a$10$HASH_FOR_Test@123', '王小红', 'xiaohong@scare.com', 'female', '1990-03-10', 1, 'active'),
(5, '13800000005', '$2a$10$HASH_FOR_Test@123', '刘师傅', 'liushifu@scare.com', 'male', '1985-07-22', 1, 'active'),
(6, '13800000006', '$2a$10$HASH_FOR_Test@123', '陈护士', 'chenhushi@scare.com', 'female', '1992-11-05', 2, 'active'),
(7, '13800000007', '$2a$10$HASH_FOR_Test@123', '赵大哥', 'zhaodage@scare.com', 'male', '1988-09-18', 2, 'active'),

-- C端用户（服务对象）
(8, '13800000008', '$2a$10$HASH_FOR_Test@123', '张大爷', NULL, 'male', '1950-05-15', NULL, 'active'),
(9, '13800000009', '$2a$10$HASH_FOR_Test@123', '李奶奶', NULL, 'female', '1955-03-20', 1, 'active'),  -- 跨端用户
(10, '13800000010', '$2a$10$HASH_FOR_Test@123', '王爷爷', NULL, 'male', '1948-11-10', NULL, 'active'),
(11, '13800000011', '$2a$10$HASH_FOR_Test@123', '孙女士', NULL, 'female', '1990-06-25', NULL, 'active'),  -- 孕妇
(12, '13800000012', '$2a$10$HASH_FOR_Test@123', '赵先生', NULL, 'male', '1965-02-14', NULL, 'active');  -- 失能人士

-- =====================================================
-- 4. 初始化用户角色（仅B端）
-- =====================================================
INSERT INTO `user_roles` (`user_id`, `role`, `station_id`, `is_primary`, `status`) VALUES
-- 管理员
(1, 'admin', NULL, TRUE, 'active'),

-- 站长
(2, 'station_manager', 1, TRUE, 'active'),
(3, 'station_manager', 2, TRUE, 'active'),

-- 工作人员
(4, 'staff', 1, TRUE, 'active'),
(5, 'staff', 1, TRUE, 'active'),
(6, 'staff', 2, TRUE, 'active'),
(7, 'staff', 2, TRUE, 'active'),

-- 李奶奶：跨端用户（既是服务对象，也是志愿者）
(9, 'staff', 1, FALSE, 'active');

-- =====================================================
-- 5. 初始化客户档案（C端服务对象）
-- =====================================================
INSERT INTO `customer_profiles` (
    `user_id`,
    `customer_type`,
    `gender`,
    `birth_date`,
    `address`,
    `health_status`,
    `disability_level`,
    `emergency_contact`
) VALUES
-- 张大爷：老年人
(8, 'elderly', 'male', '1950-05-15', '北京市朝阳区幸福小区1号楼101', '良好', '自理',
 '{"name":"张小明","phone":"13900000001","relation":"子女"}'),

-- 李奶奶：老年人（跨端用户）
(9, 'elderly', 'female', '1955-03-20', '北京市朝阳区幸福小区2号楼202', '一般', '轻度失能',
 '{"name":"李华","phone":"13900000002","relation":"子女"}'),

-- 王爷爷：老年人
(10, 'elderly', 'male', '1948-11-10', '北京市朝阳区幸福小区3号楼303', '较差', '中度失能',
 '{"name":"王芳","phone":"13900000003","relation":"子女"}'),

-- 孙女士：孕妇
(11, 'pregnant', 'female', '1990-06-25', '北京市朝阳区康乐小区5号楼501', '良好', NULL,
 '{"name":"孙先生","phone":"13900000004","relation":"配偶"}'),

-- 赵先生：失能人士
(12, 'disabled', 'male', '1965-02-14', '北京市朝阳区康乐小区6号楼602', '较差', '重度失能',
 '{"name":"赵女士","phone":"13900000005","relation":"配偶"}');

-- =====================================================
-- 6. 更新 users 表的 primary_role（仅B端用户）
-- =====================================================
UPDATE `users` u
JOIN `user_roles` ur ON u.id = ur.user_id
SET u.primary_role = ur.role
WHERE ur.is_primary = TRUE;

-- =====================================================
-- 7. Casbin 权限策略（仅B端）
-- =====================================================

-- admin: 全部权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'role:admin', '/api/v1/*', '*');

-- station_manager: 站点管理权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'role:station_manager', '/api/v1/b/requests', 'GET'),
('p', 'role:station_manager', '/api/v1/b/requests/*', 'GET'),
('p', 'role:station_manager', '/api/v1/b/tasks/pool', 'GET'),
('p', 'role:station_manager', '/api/v1/b/tasks/my', 'GET'),
('p', 'role:station_manager', '/api/v1/b/stations', 'GET'),
('p', 'role:station_manager', '/api/v1/b/stations/*', 'GET'),
('p', 'role:station_manager', '/api/v1/b/stations/*', 'PUT'),
('p', 'role:station_manager', '/api/v1/b/zones', 'GET'),
('p', 'role:station_manager', '/api/v1/b/zones', 'POST'),
('p', 'role:station_manager', '/api/v1/b/zones/*', 'PUT'),
('p', 'role:station_manager', '/api/v1/b/zones/*', 'DELETE'),
('p', 'role:station_manager', '/api/v1/b/users', 'GET'),
('p', 'role:station_manager', '/api/v1/b/users', 'POST'),
('p', 'role:station_manager', '/api/v1/b/users/*', 'GET'),
('p', 'role:station_manager', '/api/v1/b/users/*', 'PUT');

-- staff: 基础工作权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'role:staff', '/api/v1/b/tasks/pool', 'GET'),
('p', 'role:staff', '/api/v1/b/tasks/my', 'GET'),
('p', 'role:staff', '/api/v1/b/tasks/*/claim', 'POST'),
('p', 'role:staff', '/api/v1/b/tasks/*/complete', 'POST'),
('p', 'role:staff', '/api/v1/b/requests', 'GET'),
('p', 'role:staff', '/api/v1/b/requests/*', 'GET'),
('p', 'role:staff', '/api/v1/b/stations', 'GET'),
('p', 'role:staff', '/api/v1/b/stations/*', 'GET'),
('p', 'role:staff', '/api/v1/b/zones', 'GET');

-- authenticated: 公共权限（B端和C端都可访问）
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
-- B端公共
('p', 'role:authenticated', '/api/v1/b/auth/me', 'GET'),
('p', 'role:authenticated', '/api/v1/b/auth/logout', 'POST'),
('p', 'role:authenticated', '/api/v1/b/upload', 'POST'),
('p', 'role:authenticated', '/api/v1/b/notifications', 'GET'),
('p', 'role:authenticated', '/api/v1/b/notifications/*/read', 'POST'),

-- C端公共
('p', 'role:authenticated', '/api/v1/c/auth/me', 'GET'),
('p', 'role:authenticated', '/api/v1/c/auth/logout', 'POST'),
('p', 'role:authenticated', '/api/v1/c/upload', 'POST'),
('p', 'role:authenticated', '/api/v1/c/notifications', 'GET'),
('p', 'role:authenticated', '/api/v1/c/notifications/*/read', 'POST');
```

---

### 步骤10: 测试新架构

#### 10.1 执行数据库变更
```bash
cd /Users/yt/Documents/project/sCare/backend

# 1. 修改表结构（删除 type 字段）
docker exec -i scare_mysql mysql -uroot -pscare_pass scare_db << 'EOF'
ALTER TABLE user_roles DROP COLUMN IF EXISTS type;
EOF

# 2. 创建 customer_profiles 表
docker exec -i scare_mysql mysql -uroot -pscare_pass scare_db < database/migrations/create_customer_profiles.sql

# 3. 重新导入测试数据
docker exec -i scare_mysql mysql -uroot -pscare_pass scare_db < database/seeds/seed.sql
```

#### 10.2 编译检查
```bash
go build -o /dev/null ./...
```

#### 10.3 运行测试
```bash
# 运行 JWT 测试
go test ./pkg/jwt -v

# 启动后端
go run cmd/server/main.go
```

#### 10.4 功能测试

**测试1: C端登录（老年人）**
```bash
curl -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000008","password":"Test@123"}'

# 预期：成功，type="c_end"
```

**测试2: C端登录（孕妇）**
```bash
curl -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000011","password":"Test@123"}'

# 预期：成功，type="c_end"
```

**测试3: C端登录失败（B端用户）**
```bash
curl -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000004","password":"Test@123"}'

# 预期：失败，"user not registered as customer"
```

**测试4: B端登录**
```bash
curl -X POST http://localhost:8080/api/v1/b/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000004","password":"Test@123"}'

# 预期：成功，type="b_end", roles=["staff"]
```

**测试5: 跨端用户（李奶奶）**
```bash
# C端登录
curl -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000009","password":"Test@123"}'
# 预期：成功，type="c_end"

# B端登录
curl -X POST http://localhost:8080/api/v1/b/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000009","password":"Test@123"}'
# 预期：成功，type="b_end", roles=["staff"]
```

**测试6: C端 /auth/me**
```bash
# 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000011","password":"Test@123"}' | \
  grep -o '"token":"[^"]*' | cut -d'"' -f4)

# 获取用户信息
curl -X GET http://localhost:8080/api/v1/c/auth/me \
  -H "Authorization: Bearer $TOKEN"

# 预期：返回用户信息 + customer_profile（包含 customer_type="pregnant"）
```

---

## 实施检查清单

- [ ] user_roles 表的 type 字段已删除
- [ ] customer_profiles 表创建成功
- [ ] CustomerProfile 模型和 Repository 创建完成
- [ ] JWT Manager 修改完成（B端有 roles，C端无 roles）
- [ ] AuthService 拆分为 LoginBEnd 和 LoginCEnd
- [ ] CAuthHandler 使用 customerRepo 检查档案
- [ ] BAuthHandler 使用新的 LoginBEnd 方法
- [ ] routes.go 注入 customerRepo
- [ ] seed.sql 更新完成，包含多种客户类型
- [ ] 所有测试用例通过

---

## 优势总结

### 1. 更通用的命名
- **customer_profiles** 代替 elderly_profiles
- 支持"老幼病残孕"等各类服务对象
- 通过 `customer_type` 字段区分

### 2. 简化的实施流程
- 无需数据迁移，直接清空重建
- 减少了复杂的数据处理步骤
- 降低出错风险

### 3. 可扩展性
- 新增客户类型只需在 `customer_type` 枚举中添加
- 不同类型客户可以有不同的健康评估标准
- 为未来业务扩展预留空间

### 4. 清晰的架构
- users：所有人的基础账号
- user_roles：B端权限控制
- customer_profiles：C端服务对象档案
