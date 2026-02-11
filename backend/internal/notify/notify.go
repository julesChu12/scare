package notify

import "context"

type MailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

type InAppSender interface {
	Send(ctx context.Context, userID int64, title, body string) error
}
