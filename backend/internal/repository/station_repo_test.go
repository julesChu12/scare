package repository

import (
	"testing"

	"community-elderly-care-platform/internal/dao/model"
)

func TestStationRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	station := &model.ServiceStation{
		Name:      "测试站点",
		Address:   "测试地址",
		Phone:     "13800000001",
		Latitude:  30.0,
		Longitude: 120.0,
		Status:    "active",
	}

	err := repo.Create(station)
	if err != nil {
		t.Fatalf("failed to create station: %v", err)
	}

	if station.ID == 0 {
		t.Error("expected station ID to be set after creation")
	}
}

func TestStationRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	// 创建测试站点
	station := &model.ServiceStation{
		Name:    "测试站点",
		Address: "测试地址",
		Status:  "active",
	}
	repo.Create(station)

	// 查询站点
	result, err := repo.GetByID(station.ID)
	if err != nil {
		t.Fatalf("failed to get station: %v", err)
	}

	if result.Name != "测试站点" {
		t.Errorf("expected name 测试站点, got %s", result.Name)
	}
}

func TestStationRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("expected error for non-existent station, got nil")
	}
}

func TestStationRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	// 创建多个测试站点
	for i := 0; i < 5; i++ {
		station := &model.ServiceStation{
			Name:    "站点" + string(rune('A'+i)),
			Address: "地址" + string(rune('A'+i)),
			Status:  "active",
		}
		repo.Create(station)
	}

	// 测试分页
	stations, total, err := repo.List(0, 3, StationListFilter{})
	if err != nil {
		t.Fatalf("failed to list stations: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(stations) != 3 {
		t.Errorf("expected 3 stations in page, got %d", len(stations))
	}
}

func TestStationRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	// 创建测试站点
	station := &model.ServiceStation{
		Name:    "原名称",
		Address: "原地址",
		Status:  "active",
	}
	repo.Create(station)

	// 更新站点
	station.Name = "新名称"
	err := repo.Update(station)
	if err != nil {
		t.Fatalf("failed to update station: %v", err)
	}

	// 验证更新
	updated, _ := repo.GetByID(station.ID)
	if updated.Name != "新名称" {
		t.Errorf("expected name 新名称, got %s", updated.Name)
	}
}

func TestStationRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	// 创建测试站点
	station := &model.ServiceStation{
		Name:    "待删除站点",
		Address: "测试地址",
		Status:  "active",
	}
	repo.Create(station)

	// 删除站点
	err := repo.Delete(station.ID)
	if err != nil {
		t.Fatalf("failed to delete station: %v", err)
	}

	// 验证删除（软删除）
	_, err = repo.GetByID(station.ID)
	if err == nil {
		t.Error("expected error for deleted station, got nil")
	}
}

func TestStationRepository_ListActive(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	// 创建活跃站点
	for i := 0; i < 3; i++ {
		station := &model.ServiceStation{
			Name:    "活跃站点" + string(rune('A'+i)),
			Address: "地址",
			Status:  "active",
		}
		repo.Create(station)
	}

	// 创建非活跃站点
	station := &model.ServiceStation{
		Name:    "非活跃站点",
		Address: "地址",
		Status:  "inactive",
	}
	repo.Create(station)

	// 查询活跃站点
	stations, err := repo.ListActive()
	if err != nil {
		t.Fatalf("failed to list active stations: %v", err)
	}

	if len(stations) != 3 {
		t.Errorf("expected 3 active stations, got %d", len(stations))
	}
}

func TestStationRepository_GetByName(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	// 创建测试站点
	station := &model.ServiceStation{
		Name:    "唯一名称站点",
		Address: "测试地址",
		Status:  "active",
	}
	repo.Create(station)

	// 根据名称查询
	result, err := repo.GetByName("唯一名称站点")
	if err != nil {
		t.Fatalf("failed to get station by name: %v", err)
	}

	if result.ID != station.ID {
		t.Errorf("expected station ID %d, got %d", station.ID, result.ID)
	}
}

func TestStationRepository_GetByName_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStationRepository(db)

	_, err := repo.GetByName("不存在的站点")
	if err == nil {
		t.Error("expected error for non-existent station name, got nil")
	}
}
