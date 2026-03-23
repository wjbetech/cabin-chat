package store

import (
	"context"

	"github.com/wjbetech/cabin-chat/pkg/chat"
)

type UserStore interface {
	CreateUser(ctx context.Context, user chat.User) error
	GetUserByID(ctx context.Context, id string) (chat.User, error)
	GetUserByUsername(ctx context.Context, username string) (chat.User, error)
}
