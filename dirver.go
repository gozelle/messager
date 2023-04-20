package notifier

import "context"

type Driver interface {
	PushMessage(ctx context.Context, subject, content string) error
}
