package consumer

import "context"

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, batch []MessageInfo) (bool, error)
}
