package repository

import (
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/dao/query"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	q *query.Query
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		q: query.Use(db),
	}
}

// Create 创建通知
func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.q.Notification.Create(notification)
}

// ListByUser 获取用户的通知列表（分页）
func (r *NotificationRepository) ListByUser(userID int64, offset, limit int) ([]*model.Notification, int64, error) {
	n := r.q.Notification
	
	// 获取总数
	total, err := n.Where(n.UserID.Eq(userID)).Count()
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	notifications, err := n.Where(n.UserID.Eq(userID)).
		Order(n.ID.Desc()).
		Offset(offset).
		Limit(limit).
		Find()
	
	if err != nil {
		return nil, 0, err
	}
	
	return notifications, total, nil
}

// MarkRead 标记为已读
func (r *NotificationRepository) MarkRead(id int64, userID int64) error {
	n := r.q.Notification
	_, err := n.Where(n.ID.Eq(id), n.UserID.Eq(userID)).Update(n.IsRead, true)
	return err
}
