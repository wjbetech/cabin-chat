package store

import (
	"context"

	"github.com/wjbetech/cabin-chat/pkg/chat"
)

type ReactionStore interface {
	AddReaction(ctx context.Context, reaction chat.Reaction) error
	GetReactionsByMessageID(ctx context.Context, messageID string) ([]chat.Reaction, error)
	RemoveReaction(ctx context.Context, reactionID string) error
}
