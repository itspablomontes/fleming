package audit

import (
	"context"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type Log interface {
	Append(ctx context.Context, entry *Entry) error

	GetByResource(ctx context.Context, resourceID types.ID) ([]Entry, error)

	GetByActor(ctx context.Context, actor types.WalletAddress) ([]Entry, error)

	GetLatest(ctx context.Context) (*Entry, error)

	GetByID(ctx context.Context, id types.ID) (*Entry, error)

	Query(ctx context.Context, filter QueryFilter) ([]Entry, error)
}

type QueryFilter struct {
	Actor types.WalletAddress

	ResourceID types.ID

	ResourceType ResourceType

	Action Action

	StartTime *types.Timestamp

	EndTime *types.Timestamp

	Limit int

	Offset int
}

func NewQueryFilter() QueryFilter {
	return QueryFilter{
		Limit: 100,
	}
}

func (f QueryFilter) WithActor(actor types.WalletAddress) QueryFilter {
	f.Actor = actor
	return f
}

func (f QueryFilter) WithResource(id types.ID) QueryFilter {
	f.ResourceID = id
	return f
}

func (f QueryFilter) WithAction(action Action) QueryFilter {
	f.Action = action
	return f
}

func (f QueryFilter) WithLimit(limit int) QueryFilter {
	f.Limit = limit
	return f
}
