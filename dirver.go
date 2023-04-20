package notifier

import "context"

type Title = string
type Message = string

type Driver interface {
	Push(ctx context.Context, title Title, message Message) error
}
