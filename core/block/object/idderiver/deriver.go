package idderiver

import (
	"context"

	"github.com/altercpre/aparok-sync/app"

	"github.com/altercpre/aparok-core/core/domain"
)

type Deriver interface {
	app.Component
	DeriveObjectId(ctx context.Context, spaceId string, key domain.UniqueKey) (id string, err error)
}
