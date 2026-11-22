package dependencies

import (
	"github.com/altercpre/aparok-sync/app"

	"github.com/altercpre/aparok-core/space/clientspace"
)

type SpaceIndexer interface {
	app.Component
	ReindexMarketplaceSpace(space clientspace.Space) error
	ReindexSpace(space clientspace.Space) error
	RemoveIndexes(spaceID string) (err error)
	RemoveAclIndexes(spaceID string) (err error)
}
