package service

import (
	"context"
	"errors"
	"strings"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/notify"
	"community-elderly-care-platform/internal/repository"
)

var ErrNotificationInvalid = errors.New("invalid notification")

const (
	NotificationChannelInApp = "in_app"
	NotificationChannelEmail = "email"
	NotificationTypeSystem   = "system"
	NotificationTypeTask     = "task"
	NotificationTypeAlert    = "alert"
	NotificationTypeRequest  = "request"
)

type NotificationService struct {
	repo       *repository.NotificationRepository
	userRepo   *repository.UserRepository
	mailSender notify.MailSender
}

type NotificationSendOptions struct {
	Type        string
	RelatedID   int64
	RelatedType string
}

type NotificationListFilter struct {
	Type   string
	IsRead *bool
}

func NewNotificationService(repo *repository.NotificationRepository, userRepo *repository.UserRepository, mailSender notify.MailSender) *NotificationService {
	return &NotificationService{
		repo:       repo,
		userRepo:   userRepo,
		mailSender: mailSender,
	}
}

func (s *NotificationService) SendInApp(ctx context.Context, userID int64, title, body string) error {
	return s.SendInAppWithOptions(ctx, userID, title, body, NotificationSendOptions{})
}

func (s *NotificationService) SendInAppWithOptions(ctx context.Context, userID int64, title, body string, opts NotificationSendOptions) error {
	if userID == 0 || title == "" {
		return ErrNotificationInvalid
	}
	notification := &model.Notification{
		UserID:      userID,
		Title:       title,
		Body:        body,
		Type:        normalizeNotificationType(opts.Type),
		RelatedID:   opts.RelatedID,
		RelatedType: opts.RelatedType,
		Channel:     NotificationChannelInApp,
		IsRead:      false,
	}
	return s.repo.Create(notification)
}

func (s *NotificationService) SendEmail(ctx context.Context, userID int64, subject, body string) error {
	return s.SendEmailWithOptions(ctx, userID, subject, body, NotificationSendOptions{})
}

func (s *NotificationService) SendEmailWithOptions(ctx context.Context, userID int64, subject, body string, opts NotificationSendOptions) error {
	if userID == 0 || subject == "" {
		return ErrNotificationInvalid
	}
	if s.mailSender == nil {
		return ErrNotificationInvalid
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user.Email == "" {
		return ErrNotificationInvalid
	}
	if err := s.mailSender.Send(ctx, user.Email, subject, body); err != nil {
		return err
	}
	notification := &model.Notification{
		UserID:      userID,
		Title:       subject,
		Body:        body,
		Type:        normalizeNotificationType(opts.Type),
		RelatedID:   opts.RelatedID,
		RelatedType: opts.RelatedType,
		Channel:     NotificationChannelEmail,
		IsRead:      false,
	}
	return s.repo.Create(notification)
}

func (s *NotificationService) List(userID int64, page, pageSize int, filter NotificationListFilter) ([]*model.Notification, int64, error) {
	if userID == 0 {
		return nil, 0, ErrNotificationInvalid
	}
	offset := (page - 1) * pageSize
	notifications, total, err := s.repo.ListByUser(userID, offset, pageSize, repository.NotificationListFilter{
		Type:   filter.Type,
		IsRead: filter.IsRead,
	})
	if err != nil {
		return nil, 0, err
	}

	for _, notification := range notifications {
		if notification.Type == "" {
			notification.Type = inferNotificationType(notification.Title, notification.Body)
		}
	}

	return notifications, total, nil
}

func (s *NotificationService) MarkRead(userID, id int64) error {
	if userID == 0 || id == 0 {
		return ErrNotificationInvalid
	}
	return s.repo.MarkRead(id, userID)
}

func normalizeNotificationType(notificationType string) string {
	switch notificationType {
	case NotificationTypeTask, NotificationTypeAlert, NotificationTypeRequest:
		return notificationType
	case NotificationTypeSystem, "":
		return NotificationTypeSystem
	default:
		return notificationType
	}
}

func inferNotificationType(title, body string) string {
	content := title + " " + body
	if strings.Contains(content, "任务") || strings.Contains(content, "服务需求") {
		return NotificationTypeTask
	}
	return NotificationTypeSystem
}
