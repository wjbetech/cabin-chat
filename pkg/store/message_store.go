package store

import (
	"context"

	"github.com/wjbetech/cabin-chat/pkg/chat"
)

type MessageStore interface {
	CreateMessage(ctx context.Context, message chat.Message) error
	GetMessagesByRoomID(ctx context.Context, roomId string) ([]chat.Message, error)
}
