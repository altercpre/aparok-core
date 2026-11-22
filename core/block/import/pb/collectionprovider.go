package pb

import (
	"github.com/altercpre/aparok-core/core/block/collection"
	"github.com/altercpre/aparok-core/core/block/import/common"
	"github.com/altercpre/aparok-core/pb"
)

type CollectionProvider interface {
	ProvideCollection(
		snapshots *common.SnapshotContext,
		oldToNewID map[string]string,
		params *pb.RpcObjectImportRequestPbParams,
		isNewSpace bool,
	) ([]*common.Snapshot, error)
}

func GetProvider(importType pb.RpcObjectImportRequestPbParamsType, service *collection.Service) CollectionProvider {
	if importType == pb.RpcObjectImportRequestPbParams_EXPERIENCE {
		return NewGalleryImport(service)
	}
	return NewSpaceImport(service)
}
