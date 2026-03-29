package repository

import (
	"testing"

	"community-elderly-care-platform/internal/dao/model"
)

func TestZoneRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewZoneRepository(db)

	zone := &model.ServiceZone{
		StationID: 1,
		Name:      "测试围栏",
		Points:    `[[30.0,120.0],[30.1,120.0],[30.1,120.1],[30.0,120.1]]`,
		Priority:  10,
		Status:    "active",
	}

	err := repo.Create(zone)
	if err != nil {
		t.Fatalf("failed to create zone: %v", err)
	}

	if zone.ID == 0 {
		t.Error("expected zone ID to be set after creation")
	}
}

func TestZoneRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewZoneRepository(db)

	// 创建测试围栏
	zone := &model.ServiceZone{
		StationID: 1,
		Name:      "测试围栏",
		Points:    `[[30.0,120.0],[30.1,120.0],[30.1,120.1],[30.0,120.1]]`,
		Status:    "active",
	}
	repo.Create(zone)

	// 查询围栏
	result, err := repo.GetByID(zone.ID)
	if err != nil {
		t.Fatalf("failed to get zone: %v", err)
	}

	if result.Name != "测试围栏" {
		t.Errorf("expected name 测试围栏, got %s", result.Name)
	}
}

func TestZoneRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewZoneRepository(db)

	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("expected error for non-existent zone, got nil")
	}
}

func TestZoneRepository_ListActive(t *testing.T) {
	db := setupTestDB(t)
	repo := NewZoneRepository(db)

	// 创建活跃围栏
	for i := 0; i < 3; i++ {
		zone := &model.ServiceZone{
			StationID: 1,
			Name:      "活跃围栏" + string(rune('A'+i)),
			Points:    `[[30.0,120.0],[30.1,120.0],[30.1,120.1],[30.0,120.1]]`,
			Priority:  int64(i * 10),
			Status:    "active",
		}
		repo.Create(zone)
	}

	// 创建非活跃围栏
	zone := &model.ServiceZone{
		StationID: 1,
		Name:      "非活跃围栏",
		Points:    `[[30.0,120.0],[30.1,120.0],[30.1,120.1],[30.0,120.1]]`,
		Status:    "inactive",
	}
	repo.Create(zone)

	// 查询活跃围栏
	zones, err := repo.ListActive()
	if err != nil {
		t.Fatalf("failed to list active zones: %v", err)
	}

	if len(zones) != 3 {
		t.Errorf("expected 3 active zones, got %d", len(zones))
	}
}

func TestZoneRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewZoneRepository(db)

	// 创建多个测试围栏
	for i := 0; i < 5; i++ {
		zone := &model.ServiceZone{
			StationID: 1,
			Name:      "围栏" + string(rune('A'+i)),
			Points:    `[[30.0,120.0],[30.1,120.0],[30.1,120.1],[30.0,120.1]]`,
			Status:    "active",
		}
		repo.Create(zone)
	}

	// 测试分页
	zones, total, err := repo.List(0, 3, ZoneListFilter{})
	if err != nil {
		t.Fatalf("failed to list zones: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(zones) != 3 {
		t.Errorf("expected 3 zones in page, got %d", len(zones))
	}
}

func TestZoneRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewZoneRepository(db)

	// 创建测试围栏
	zone := &model.ServiceZone{
		StationID: 1,
		Name:      "原名称",
		Points:    `[[30.0,120.0],[30.1,120.0],[30.1,120.1],[30.0,120.1]]`,
		Priority:  10,
		Status:    "active",
	}
	repo.Create(zone)

	// 更新围栏
	zone.Name = "新名称"
	zone.Priority = 20
	err := repo.Update(zone)
	if err != nil {
		t.Fatalf("failed to update zone: %v", err)
	}

	// 验证更新
	updated, _ := repo.GetByID(zone.ID)
	if updated.Name != "新名称" {
		t.Errorf("expected name 新名称, got %s", updated.Name)
	}
	if updated.Priority != 20 {
		t.Errorf("expected priority 20, got %d", updated.Priority)
	}
}

func TestZoneRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewZoneRepository(db)

	// 创建测试围栏
	zone := &model.ServiceZone{
		StationID: 1,
		Name:      "待删除围栏",
		Points:    `[[30.0,120.0],[30.1,120.0],[30.1,120.1],[30.0,120.1]]`,
		Status:    "active",
	}
	repo.Create(zone)

	// 删除围栏
	err := repo.Delete(zone.ID)
	if err != nil {
		t.Fatalf("failed to delete zone: %v", err)
	}

	// 验证删除（软删除）
	_, err = repo.GetByID(zone.ID)
	if err == nil {
		t.Error("expected error for deleted zone, got nil")
	}
}
