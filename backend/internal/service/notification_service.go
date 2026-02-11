package service

import (
	"context"
	"errors"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/notify"
	"community-elderly-care-platform/internal/repository"
)

var ErrNotificationInvalid = errors.New("invalid notification")

const (
	NotificationChannelInApp = "in_app"
	NotificationChannelEmail = "email"
)

type NotificationService struct {
	repo       *repository.NotificationRepository
	userRepo   *repository.UserRepository
	mailSender notify.MailSender
}

func NewNotificationService(repo *repository.NotificationRepository, userRepo *repository.UserRepository, mailSender notify.MailSender) *NotificationService {
	return &NotificationService{
		repo:       repo,
		userRepo:   userRepo,
		mailSender: mailSender,
	}
}

func (s *NotificationService) SendInApp(ctx context.Context, userID int64, title, body string) error {
	if userID == 0 || title == "" {
		return ErrNotificationInvalid
	}
	notification := &model.Notification{
		UserID:  userID,
		Title:   title,
		Body:    body,
		Channel: NotificationChannelInApp,
		IsRead:  false,
	}
	return s.repo.Create(notification)
}

func (s *NotificationService) SendEmail(ctx context.Context, userID int64, subject, body string) error {
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
		UserID:  userID,
		Title:   subject,
		Body:    body,
		Channel: NotificationChannelEmail,
		IsRead:  false,
	}
	return s.repo.Create(notification)
}

func (s *NotificationService) List(userID int64, page, pageSize int) ([]*model.Notification, int64, error) {
	if userID == 0 {
		return nil, 0, ErrNotificationInvalid
	}
	offset := (page - 1) * pageSize
	return s.repo.ListByUser(userID, offset, pageSize)
}

func (s *NotificationService) MarkRead(userID, id int64) error {
	if userID == 0 || id == 0 {
		return ErrNotificationInvalid
	}
	return s.repo.MarkRead(id, userID)
}
