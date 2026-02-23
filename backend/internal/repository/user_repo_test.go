package repository

import (
	"testing"

	"community-elderly-care-platform/internal/dao/model"
)

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Phone:        "13800000001",
		PasswordHash: "hashed_password",
		Name:         "测试用户",
		Status:       "active",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("expected user ID to be set after creation")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user := &model.User{
		Phone:        "13800000002",
		PasswordHash: "hashed_password",
		Name:         "测试用户2",
		Status:       "active",
	}
	repo.Create(user)

	// 查询用户
	result, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if result.Phone != "13800000002" {
		t.Errorf("expected phone 13800000002, got %s", result.Phone)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestUserRepository_GetByPhone(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user := &model.User{
		Phone:        "13800000003",
		PasswordHash: "hashed_password",
		Name:         "测试用户3",
		Status:       "active",
	}
	repo.Create(user)

	// 根据手机号查询
	result, err := repo.GetByPhone("13800000003")
	if err != nil {
		t.Fatalf("failed to get user by phone: %v", err)
	}

	if result.Name != "测试用户3" {
		t.Errorf("expected name 测试用户3, got %s", result.Name)
	}
}

func TestUserRepository_GetByPhone_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.GetByPhone("19999999999")
	if err == nil {
		t.Error("expected error for non-existent phone, got nil")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user := &model.User{
		Phone:        "13800000004",
		PasswordHash: "hashed_password",
		Name:         "原名称",
		Status:       "active",
	}
	repo.Create(user)

	// 更新用户
	user.Name = "新名称"
	err := repo.Update(user)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	// 验证更新
	updated, _ := repo.GetByID(user.ID)
	if updated.Name != "新名称" {
		t.Errorf("expected name 新名称, got %s", updated.Name)
	}
}

func TestUserRepository_ListWithFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user1 := &model.User{
		Phone:        "13800000020",
		PasswordHash: "hashed_password",
		Name:         "张三",
		Status:       "active",
		StationID:    1,
	}
	repo.Create(user1)

	user2 := &model.User{
		Phone:        "13800000021",
		PasswordHash: "hashed_password",
		Name:         "李四",
		Status:       "inactive",
		StationID:    2,
	}
	repo.Create(user2)

	// 按状态筛选
	users, total, err := repo.ListWithFilter(0, 10, UserFilter{Status: "active"})
	if err != nil {
		t.Fatalf("failed to list users with filter: %v", err)
	}

	if total != 1 {
		t.Errorf("expected 1 active user, got %d", total)
	}

	// 按站点筛选
	users, total, err = repo.ListWithFilter(0, 10, UserFilter{StationID: 1})
	if err != nil {
		t.Fatalf("failed to list users with station filter: %v", err)
	}

	if total != 1 {
		t.Errorf("expected 1 user in station 1, got %d", total)
	}

	// 按关键词搜索
	users, total, err = repo.ListWithFilter(0, 10, UserFilter{Keyword: "张"})
	if err != nil {
		t.Fatalf("failed to list users with keyword filter: %v", err)
	}

	if total != 1 {
		t.Errorf("expected 1 user matching keyword, got %d", total)
	}

	if len(users) > 0 && users[0].Name != "张三" {
		t.Errorf("expected user 张三, got %s", users[0].Name)
	}
}
