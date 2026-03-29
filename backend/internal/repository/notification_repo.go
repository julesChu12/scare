package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	q *query.Query
}

type NotificationListFilter struct {
	Type   string
	IsRead *bool
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		q: query.Use(db),
	}
}

// Create 创建通知
// Omit sent_at/read_at：这两个字段在创建时应为 NULL，
// 但 Gen 模型使用 time.Time（值类型），零值会被序列化为 '0000-00-00' 导致 MySQL 报错
func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.q.Notification.UnderlyingDB().
		Omit("sent_at", "read_at").
		Create(notification).Error
}

// ListByUser 获取用户的通知列表（分页）
func (r *NotificationRepository) ListByUser(userID int64, offset, limit int, filter NotificationListFilter) ([]*model.Notification, int64, error) {
	db := r.q.Notification.UnderlyingDB().Model(&model.Notification{}).Where("user_id = ?", userID)

	if filter.Type != "" {
		db = applyNotificationTypeFilter(db, filter.Type)
	}

	if filter.IsRead != nil {
		db = db.Where("is_read = ?", *filter.IsRead)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var notifications []*model.Notification
	if err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func applyNotificationTypeFilter(db *gorm.DB, notificationType string) *gorm.DB {
	switch notificationType {
	case "task":
		return db.Where(
			"type = ? OR ((type IS NULL OR type = '') AND (title LIKE ? OR title LIKE ?))",
			"task", "%任务%", "%服务需求%",
		)
	case "system":
		return db.Where(
			"type = ? OR ((type IS NULL OR type = '') AND title NOT LIKE ? AND title NOT LIKE ?)",
			"system", "%任务%", "%服务需求%",
		)
	default:
		return db.Where("type = ?", notificationType)
	}
}

// MarkRead 标记为已读
func (r *NotificationRepository) MarkRead(id int64, userID int64) error {
	n := r.q.Notification
	_, err := n.Where(n.ID.Eq(id), n.UserID.Eq(userID)).Update(n.IsRead, true)
	return err
}
