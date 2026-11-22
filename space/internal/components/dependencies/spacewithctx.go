package dependencies

import (
	"context"

	"github.com/altercpre/aparok-core/core/block/editor/smartblock"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/pkg/lib/threads"
)

type SpaceWithCtx interface {
	DoCtx(ctx context.Context, objectId string, apply func(sb smartblock.SmartBlock) error) error
	Id() string
	DerivedIDs() threads.DerivedSmartblockIds
	DeriveObjectID(ctx context.Context, uniqueKey domain.UniqueKey) (id string, err error)
}
