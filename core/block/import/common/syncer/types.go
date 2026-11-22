package syncer

import (
	"context"

	"github.com/altercpre/aparok-core/core/block"
	"github.com/altercpre/aparok-core/core/block/editor/smartblock"
	"github.com/altercpre/aparok-core/core/block/simple"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/core/domain/objectorigin"
	"github.com/altercpre/aparok-core/core/session"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

type BlockService interface {
	GetObject(ctx context.Context, objectID string) (sb smartblock.SmartBlock, err error)
	GetObjectByFullID(ctx context.Context, id domain.FullID) (sb smartblock.SmartBlock, err error)
	UploadFile(ctx context.Context, spaceId string, req block.FileUploadRequest) (objectId string, fileType model.BlockContentFileType, details *domain.Details, err error)
	UploadBlockFile(ctx session.Context, req block.UploadRequest, groupID string, isSync bool) (fileObjectId string, err error)
}

type Syncer interface {
	Sync(id domain.FullID, newIdsSet map[string]struct{}, b simple.Block, origin objectorigin.ObjectOrigin) error
}
