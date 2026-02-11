# 老年人档案架构重构实施计划

## 目标
将 C端用户从 `user_roles` 表中独立出来，使用 `elderly_profiles` 表管理老年人档案。

## 核心变更

### 概念变化
- **之前**：elderly/family 作为角色存储在 `user_roles` 表，type='c_end'
- **之后**：
  - 去掉 family 概念
  - elderly 用户通过 `elderly_profiles` 表标识
  - `user_roles` 表只用于 B端权限控制，去掉 type 字段

### 用户身份判定
```
B端用户: 在 user_roles 表中有记录
C端用户: 在 elderly_profiles 表中有记录
跨端用户: 同时在两张表中有记录
```

---

## 步骤1: 创建 elderly_profiles 表

### SQL 脚本
```sql
-- 创建老年人档案表
CREATE TABLE `elderly_profiles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL COMMENT '关联用户ID',

  -- 基本信息
  `id_card` VARCHAR(18) DEFAULT NULL COMMENT '身份证号',
  `gender` VARCHAR(10) DEFAULT NULL COMMENT '性别',
  `birth_date` DATE DEFAULT NULL COMMENT '出生日期',
  `address` TEXT DEFAULT NULL COMMENT '居住地址',

  -- 健康信息
  `health_status` VARCHAR(50) DEFAULT NULL COMMENT '健康状况',
  `disability_level` VARCHAR(20) DEFAULT NULL COMMENT '失能等级',
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
  CONSTRAINT `fk_elderly_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='老年人档案表';
```

### 验证
```sql
-- 确认表结构
DESCRIBE elderly_profiles;

-- 确认索引
SHOW INDEX FROM elderly_profiles;
```

---

## 步骤2: 数据迁移

### 2.1 分析现有数据
```sql
-- 查看 c_end 用户分布
SELECT role, COUNT(*) as count
FROM user_roles
WHERE type = 'c_end'
GROUP BY role;

-- 预期结果：
-- elderly: 3-4条
-- family: 3-4条

-- 查看具体用户
SELECT ur.user_id, u.name, u.phone, ur.role, ur.type
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
WHERE ur.type = 'c_end'
ORDER BY ur.user_id, ur.role;
```

### 2.2 迁移 elderly 用户到 elderly_profiles
```sql
-- 将 role='elderly' 的用户迁移到 elderly_profiles
INSERT INTO elderly_profiles (
    user_id,
    gender,
    birth_date,
    address,
    created_at,
    updated_at
)
SELECT
    ur.user_id,
    u.gender,
    u.birth_date,
    NULL as address,  -- 地址需要后续完善
    NOW(3),
    NOW(3)
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
WHERE ur.type = 'c_end'
  AND ur.role = 'elderly';

-- 验证迁移结果
SELECT
    ep.id,
    ep.user_id,
    u.name,
    u.phone,
    ep.gender,
    ep.birth_date
FROM elderly_profiles ep
JOIN users u ON ep.user_id = u.id;
```

### 2.3 处理 family 用户
```sql
-- family 用户不迁移到 elderly_profiles
-- 直接删除这些记录（步骤3会执行）

-- 先确认这些用户只有 family 角色，没有其他B端角色
SELECT
    ur.user_id,
    u.name,
    GROUP_CONCAT(ur.role) as roles,
    GROUP_CONCAT(ur.type) as types
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
WHERE ur.user_id IN (
    SELECT user_id FROM user_roles WHERE type = 'c_end' AND role = 'family'
)
GROUP BY ur.user_id, u.name;

-- 如果某个 family 用户同时有 B端角色，需要特殊处理
```

---

## 步骤3: 清理 user_roles 表中的 c_end 记录

### 3.1 删除所有 type='c_end' 的记录
```sql
-- 备份要删除的数据（可选）
CREATE TABLE user_roles_c_end_backup AS
SELECT * FROM user_roles WHERE type = 'c_end';

-- 删除 c_end 记录
DELETE FROM user_roles WHERE type = 'c_end';

-- 验证删除结果
SELECT COUNT(*) as remaining_c_end_count
FROM user_roles
WHERE type = 'c_end';
-- 预期：0

-- 确认剩余的都是 B端角色
SELECT role, type, COUNT(*) as count
FROM user_roles
GROUP BY role, type;
-- 预期：只有 admin/station_manager/staff，type='b_end'
```

---

## 步骤4: 修改 user_roles 表结构（删除 type 字段）

### 4.1 确认安全性
```sql
-- 再次确认没有 c_end 记录
SELECT COUNT(*) FROM user_roles WHERE type = 'c_end';
-- 必须为 0

-- 确认所有记录都是 b_end
SELECT COUNT(*) FROM user_roles WHERE type != 'b_end';
-- 必须为 0
```

### 4.2 删除 type 字段
```sql
-- 删除 type 字段
ALTER TABLE user_roles DROP COLUMN type;

-- 验证
DESCRIBE user_roles;
-- 确认 type 字段已不存在
```

### 4.3 清理冗余索引（如果有）
```sql
-- 检查索引
SHOW INDEX FROM user_roles;

-- 如果有基于 type 的索引，删除它们
-- ALTER TABLE user_roles DROP INDEX idx_type;  -- 如果存在
```

---

## 步骤5: 创建 ElderlyProfile 领域模型

### 5.1 创建文件：`internal/domain/elderly_profile.go`
```go
package domain

import "time"

// ElderlyProfile 老年人档案
type ElderlyProfile struct {
    ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    int64     `gorm:"uniqueIndex;not null" json:"user_id"`

    // 基本信息
    IDCard    *string   `gorm:"type:varchar(18)" json:"id_card,omitempty"`
    Gender    *string   `gorm:"type:varchar(10)" json:"gender,omitempty"`
    BirthDate *string   `gorm:"type:date" json:"birth_date,omitempty"`
    Address   *string   `gorm:"type:text" json:"address,omitempty"`

    // 健康信息
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

func (ElderlyProfile) TableName() string {
    return "elderly_profiles"
}
```

### 5.2 创建文件：`internal/repository/elderly_repository.go`
```go
package repository

import (
    "community-elderly-care-platform/internal/domain"
    "gorm.io/gorm"
)

type ElderlyRepository struct {
    db *gorm.DB
}

func NewElderlyRepository(db *gorm.DB) *ElderlyRepository {
    return &ElderlyRepository{db: db}
}

// GetByUserID 根据用户ID获取档案
func (r *ElderlyRepository) GetByUserID(userID int64) (*domain.ElderlyProfile, error) {
    var profile domain.ElderlyProfile
    err := r.db.Where("user_id = ?", userID).First(&profile).Error
    if err != nil {
        return nil, err
    }
    return &profile, nil
}

// GetByUserIDWithUser 获取档案并预加载用户信息
func (r *ElderlyRepository) GetByUserIDWithUser(userID int64) (*domain.ElderlyProfile, error) {
    var profile domain.ElderlyProfile
    err := r.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error
    if err != nil {
        return nil, err
    }
    return &profile, nil
}

// Create 创建档案
func (r *ElderlyRepository) Create(profile *domain.ElderlyProfile) error {
    return r.db.Create(profile).Error
}

// Update 更新档案
func (r *ElderlyRepository) Update(profile *domain.ElderlyProfile) error {
    return r.db.Save(profile).Error
}

// Delete 删除档案
func (r *ElderlyRepository) Delete(userID int64) error {
    return r.db.Where("user_id = ?", userID).Delete(&domain.ElderlyProfile{}).Error
}

// Exists 检查用户是否有档案
func (r *ElderlyRepository) Exists(userID int64) (bool, error) {
    var count int64
    err := r.db.Model(&domain.ElderlyProfile{}).Where("user_id = ?", userID).Count(&count).Error
    return count > 0, err
}
```

---

## 步骤6: 简化 JWT Manager

### 6.1 修改 `pkg/jwt/manager.go`

#### 修改 Claims 结构
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
```

#### 修改方法签名
```go
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

// GenerateRefreshToken 生成刷新令牌（同样修改）
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

### 6.2 更新测试文件 `pkg/jwt/manager_test.go`
```go
func TestGenerateToken(t *testing.T) {
    manager := NewManager("test_secret", 24*time.Hour, 7*24*time.Hour)

    // 测试 B端 Token
    token, err := manager.GenerateToken(42, "b_end", 7, []string{"admin", "staff"})
    if err != nil {
        t.Fatalf("生成token失败: %v", err)
    }

    claims, err := manager.ValidateToken(token)
    if err != nil {
        t.Fatalf("验证token失败: %v", err)
    }

    if claims.UserID != 42 {
        t.Errorf("UserID错误: expected 42, got %d", claims.UserID)
    }

    if claims.Type != "b_end" {
        t.Errorf("Type错误: expected b_end, got %s", claims.Type)
    }

    if len(claims.Roles) != 2 {
        t.Errorf("Roles数量错误: expected 2, got %d", len(claims.Roles))
    }

    // 测试 C端 Token
    cToken, err := manager.GenerateToken(100, "c_end", 0, nil)
    if err != nil {
        t.Fatalf("生成C端token失败: %v", err)
    }

    cClaims, err := manager.ValidateToken(cToken)
    if err != nil {
        t.Fatalf("验证C端token失败: %v", err)
    }

    if cClaims.Type != "c_end" {
        t.Errorf("Type错误: expected c_end, got %s", cClaims.Type)
    }

    if len(cClaims.Roles) != 0 {
        t.Errorf("C端Roles应为空: got %v", cClaims.Roles)
    }
}
```

---

## 步骤7: 重写 C端登录逻辑

### 7.1 修改 `internal/service/auth_service.go`

#### 删除旧的 Login 方法，拆分为两个方法
```go
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

// LoginCEnd C端登录（新增）
func (s *AuthService) LoginCEnd(phone, password string, elderlyRepo *repository.ElderlyRepository) (*Tokens, *domain.User, error) {
    user, err := s.ValidateCredentials(phone, password)
    if err != nil {
        return nil, nil, err
    }

    // 检查是否有老年人档案
    exists, err := elderlyRepo.Exists(user.ID)
    if err != nil || !exists {
        return nil, nil, ErrNoElderlyProfile
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

#### 添加新的错误类型
```go
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserInactive       = errors.New("user inactive")
    ErrNoRoleForEndType   = errors.New("user has no role for this end type")
    ErrNoElderlyProfile   = errors.New("user has no elderly profile")  // 新增
)
```

### 7.2 修改 `internal/handler/b_auth_handler.go`
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

### 7.3 完全重写 `internal/handler/c_auth_handler.go`
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
    elderlyRepo  *repository.ElderlyRepository
}

func NewCAuthHandler(
    authService *service.AuthService,
    userRepo *repository.UserRepository,
    elderlyRepo *repository.ElderlyRepository,
) *CAuthHandler {
    return &CAuthHandler{
        authService: authService,
        userRepo:    userRepo,
        elderlyRepo: elderlyRepo,
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
    tokens, user, err := h.authService.LoginCEnd(req.Phone, req.Password, h.elderlyRepo)
    if err != nil {
        if err == service.ErrInvalidCredentials {
            RespondError(c, http.StatusUnauthorized, "invalid credentials")
            return
        }
        if err == service.ErrUserInactive {
            RespondError(c, http.StatusForbidden, "user inactive")
            return
        }
        if err == service.ErrNoElderlyProfile {
            RespondError(c, http.StatusForbidden, "user not registered as elderly")
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

    // 获取老年人档案
    profile, err := h.elderlyRepo.GetByUserID(userID)
    if err != nil {
        RespondError(c, http.StatusForbidden, "user has no elderly profile")
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

### 7.4 修改 `internal/handler/routes.go`
```go
// 初始化 Repository
userRepo := repository.NewUserRepository(db.DB)
elderlyRepo := repository.NewElderlyRepository(db.DB)  // 新增
stationRepo := repository.NewStationRepository(db.DB)
// ... 其他 repo

// 初始化 Handler
bAuthHandler := NewBAuthHandler(authService, userRepo, permissionService)
cAuthHandler := NewCAuthHandler(authService, userRepo, elderlyRepo)  // 修改参数
```

---

## 步骤8: 更新测试数据 seed.sql

### 8.1 修改 `database/seeds/seed.sql`

```sql
-- =====================================================
-- 9. 初始化用户角色关联表（仅B端）
-- =====================================================

-- 删除旧的注释和 INSERT
-- 去掉所有 type='c_end' 的记录

INSERT INTO `user_roles` (`user_id`, `role`, `station_id`, `is_primary`, `status`) VALUES
-- 管理员(ID=1)
(1, 'admin', NULL, TRUE, 'active'),

-- 站长们(ID=2,3)
(2, 'station_manager', 1, TRUE, 'active'),
(3, 'station_manager', 2, TRUE, 'active'),

-- 工作人员(ID=4,5,6,7)
(4, 'staff', 1, TRUE, 'active'),
(5, 'staff', 1, TRUE, 'active'),
(6, 'staff', 2, TRUE, 'active'),
(7, 'staff', 2, TRUE, 'active');

-- 李奶奶(ID=9)：既是老年人，也是志愿者
INSERT INTO `user_roles` (`user_id`, `role`, `station_id`, `is_primary`, `status`)
VALUES (9, 'staff', 1, FALSE, 'active');

-- =====================================================
-- 10. 初始化老年人档案表
-- =====================================================

INSERT INTO `elderly_profiles` (
    `user_id`,
    `gender`,
    `birth_date`,
    `address`,
    `health_status`,
    `disability_level`,
    `emergency_contact`
) VALUES
-- 张大爷(ID=8)
(8, 'male', '1950-05-15', '北京市朝阳区幸福小区1号楼101', '良好', '自理',
 '{"name":"张小明","phone":"13900000001","relation":"子女"}'),

-- 李奶奶(ID=9)：跨端用户
(9, 'female', '1955-03-20', '北京市朝阳区幸福小区2号楼202', '一般', '轻度失能',
 '{"name":"李华","phone":"13900000002","relation":"子女"}'),

-- 王爷爷(ID=10)
(10, 'male', '1948-11-10', '北京市朝阳区幸福小区3号楼303', '较差', '中度失能',
 '{"name":"王芳","phone":"13900000003","relation":"子女"}');

-- =====================================================
-- 11. 更新 users 表的 primary_role（仅B端用户）
-- =====================================================

UPDATE `users` u
JOIN `user_roles` ur ON u.id = ur.user_id
SET u.primary_role = ur.role
WHERE ur.is_primary = TRUE;

-- =====================================================
-- 12. Casbin 策略（保持不变，去掉 authenticated 的说明）
-- =====================================================

-- 注：C端用户无需 Casbin 权限控制
-- 只有 B端角色需要权限策略
```

### 8.2 重新导入测试数据
```bash
# 1. 备份当前数据库
docker exec scare_mysql mysqldump -uroot -pscare_pass scare_db > backup_before_migration.sql

# 2. 清空相关表
docker exec -i scare_mysql mysql -uroot -pscare_pass scare_db << 'EOF'
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE user_roles;
TRUNCATE TABLE elderly_profiles;
SET FOREIGN_KEY_CHECKS = 1;
EOF

# 3. 重新导入
docker exec -i scare_mysql mysql -uroot -pscare_pass scare_db < database/seeds/seed.sql
```

---

## 步骤9: 测试新架构

### 9.1 编译检查
```bash
cd /Users/yt/Documents/project/sCare/backend
go build -o /dev/null ./...
```

### 9.2 启动后端
```bash
go run cmd/server/main.go
```

### 9.3 测试用例

#### 测试1: C端登录（老年人）
```bash
# 张大爷登录 (user_id=8, 有 elderly_profile)
curl -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000008","password":"Test@123"}'

# 预期响应：
{
  "msg": "ok",
  "data": {
    "token": "...",
    "refresh_token": "...",
    "user_id": 8,
    "type": "c_end",
    "name": "张大爷",
    "phone": "13800000008",
    "status": "active"
  }
}
```

#### 测试2: C端登录失败（非老年人）
```bash
# 王小红登录 (user_id=4, 只有 B端角色，无 elderly_profile)
curl -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000004","password":"Test@123"}'

# 预期响应：
{
  "msg": "user not registered as elderly",
  "data": null
}
```

#### 测试3: B端登录（正常）
```bash
# 王小红登录 (user_id=4, 有 staff 角色)
curl -X POST http://localhost:8080/api/v1/b/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000004","password":"Test@123"}'

# 预期响应：
{
  "msg": "ok",
  "data": {
    "token": "...",
    "refresh_token": "...",
    "user_id": 4,
    "roles": ["staff"],
    "type": "b_end",
    "station_id": 1,
    "name": "王小红",
    "phone": "13800000004",
    "status": "active"
  }
}
```

#### 测试4: 跨端用户（李奶奶）
```bash
# C端登录 (user_id=9, 有 elderly_profile)
curl -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000009","password":"Test@123"}'
# 预期成功，type: "c_end"

# B端登录 (user_id=9, 有 staff 角色)
curl -X POST http://localhost:8080/api/v1/b/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000009","password":"Test@123"}'
# 预期成功，type: "b_end", roles: ["staff"]
```

#### 测试5: C端 /auth/me 接口
```bash
# 先登录获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/c/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000008","password":"Test@123"}' | \
  grep -o '"token":"[^"]*' | cut -d'"' -f4)

# 获取用户信息
curl -X GET http://localhost:8080/api/v1/c/auth/me \
  -H "Authorization: Bearer $TOKEN"

# 预期响应：
{
  "msg": "ok",
  "data": {
    "user": {
      "id": 8,
      "name": "张大爷",
      "phone": "13800000008",
      "email": null,
      "avatar": null,
      "status": "active"
    },
    "profile": {
      "gender": "male",
      "birth_date": "1950-05-15",
      "address": "北京市朝阳区幸福小区1号楼101",
      "health_status": "良好",
      "disability_level": "自理",
      "emergency_contact": {
        "name": "张小明",
        "phone": "13900000001",
        "relation": "子女"
      }
    }
  }
}
```

---

## 回滚方案

如果迁移过程中出现问题，按以下步骤回滚：

### 1. 恢复数据库
```bash
# 恢复到迁移前的状态
docker exec -i scare_mysql mysql -uroot -pscare_pass scare_db < backup_before_migration.sql
```

### 2. 删除新表
```sql
DROP TABLE IF EXISTS elderly_profiles;
```

### 3. 恢复代码
```bash
git checkout internal/handler/c_auth_handler.go
git checkout internal/service/auth_service.go
git checkout pkg/jwt/manager.go
# ... 恢复其他修改的文件
```

---

## 注意事项和风险

### 高风险项
1. **数据迁移**：确保迁移前做好备份
2. **外键约束**：删除 user_roles 记录前确认无其他表依赖
3. **JWT 兼容性**：旧的 Token 会失效，需要用户重新登录

### 中风险项
1. **编译错误**：修改了多个文件，可能遗漏某些引用
2. **测试覆盖**：需要全面测试 B端、C端、跨端用户场景

### 建议
1. **分步执行**：按照步骤顺序执行，每步验证后再继续
2. **多次备份**：关键步骤前都备份数据库
3. **测试环境先行**：如果有测试环境，先在测试环境完整执行一遍

---

## 时间估算

- 步骤1-4（数据库）：30分钟
- 步骤5-6（代码模型）：20分钟
- 步骤7（重写登录）：40分钟
- 步骤8（测试数据）：15分钟
- 步骤9（测试验证）：30分钟

**总计**：约 2-2.5 小时

---

## 检查清单

迁移完成后，确认以下项目：

- [ ] elderly_profiles 表创建成功，索引正确
- [ ] 所有 elderly 用户已迁移到 elderly_profiles
- [ ] user_roles 表中无 type='c_end' 记录
- [ ] user_roles 表的 type 字段已删除
- [ ] ElderlyProfile 模型和 Repository 编译通过
- [ ] JWT Manager 支持 B端（含 roles）和 C端（无 roles）
- [ ] C端登录只检查 elderly_profile，不检查角色
- [ ] B端登录功能正常，返回多角色
- [ ] 跨端用户（李奶奶）两端都能登录
- [ ] C端 /auth/me 返回用户信息 + 档案信息
- [ ] 所有测试用例通过
