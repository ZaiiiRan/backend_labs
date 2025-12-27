package consumer

import "context"

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, batch []Message) (bool, error)
}
