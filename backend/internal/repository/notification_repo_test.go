package repository

import (
	"testing"

	"community-elderly-care-platform/internal/dao/model"
)

func TestNotificationRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNotificationRepository(db)

	notification := &model.Notification{
		UserID: 1,
		Title:  "测试通知",
		Body:   "这是一条测试通知",
		Type:   "system",
		IsRead: false,
	}

	err := repo.Create(notification)
	if err != nil {
		t.Fatalf("failed to create notification: %v", err)
	}

	if notification.ID == 0 {
		t.Error("expected notification ID to be set after creation")
	}
}

func TestNotificationRepository_ListByUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNotificationRepository(db)

	// 创建用户1的通知
	for i := 0; i < 5; i++ {
		notification := &model.Notification{
			UserID: 1,
			Title:  "通知" + string(rune('A'+i)),
			Body:   "内容",
			Type:   "system",
			IsRead: false,
		}
		repo.Create(notification)
	}

	// 创建用户2的通知
	notification := &model.Notification{
		UserID: 2,
		Title:  "用户2的通知",
		Body:   "内容",
		Type:   "system",
		IsRead: false,
	}
	repo.Create(notification)

	// 测试分页
	notifications, total, err := repo.ListByUser(1, 0, 3)
	if err != nil {
		t.Fatalf("failed to list notifications: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(notifications) != 3 {
		t.Errorf("expected 3 notifications in page, got %d", len(notifications))
	}
}

func TestNotificationRepository_MarkRead(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNotificationRepository(db)

	// 创建测试通知
	notification := &model.Notification{
		UserID: 1,
		Title:  "未读通知",
		Body:   "内容",
		Type:   "system",
		IsRead: false,
	}
	repo.Create(notification)

	// 标记为已读
	err := repo.MarkRead(notification.ID, 1)
	if err != nil {
		t.Fatalf("failed to mark notification as read: %v", err)
	}

	// 验证更新
	var updated model.Notification
	db.First(&updated, notification.ID)
	if !updated.IsRead {
		t.Error("expected notification to be marked as read")
	}
}

func TestNotificationRepository_MarkRead_WrongUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewNotificationRepository(db)

	// 创建测试通知
	notification := &model.Notification{
		UserID: 1,
		Title:  "用户1的通知",
		Body:   "内容",
		Type:   "system",
		IsRead: false,
	}
	repo.Create(notification)

	// 尝试用错误的用户ID标记为已读
	err := repo.MarkRead(notification.ID, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 验证通知仍然未读（因为用户ID不匹配）
	var updated model.Notification
	db.First(&updated, notification.ID)
	if updated.IsRead {
		t.Error("expected notification to remain unread for wrong user")
	}
}
