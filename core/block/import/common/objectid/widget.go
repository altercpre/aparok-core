package objectid

import (
	"context"
	"fmt"
	"time"

	"github.com/altercpre/aparok-sync/commonspace/object/tree/treestorage"

	"github.com/altercpre/aparok-core/core/block/import/common"
	"github.com/altercpre/aparok-core/core/domain/objectorigin"
	"github.com/altercpre/aparok-core/space"
)

type widget struct {
	spaceService space.Service
}

func newWidget(spaceService space.Service) *widget {
	return &widget{spaceService: spaceService}
}

func (w widget) GetIDAndPayload(ctx context.Context, spaceID string, _ *common.Snapshot, _ time.Time, _ bool, _ objectorigin.ObjectOrigin) (string, treestorage.TreeStorageCreatePayload, error) {
	spc, err := w.spaceService.Get(ctx, spaceID)
	if err != nil {
		return "", treestorage.TreeStorageCreatePayload{}, fmt.Errorf("get space : %w", err)
	}
	return spc.DerivedIDs().Widgets, treestorage.TreeStorageCreatePayload{}, nil
}
