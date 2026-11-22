package csv

import (
	"github.com/altercpre/aparok-core/core/block/import/common"
	"github.com/altercpre/aparok-core/core/block/process"
	"github.com/altercpre/aparok-core/pb"
)

type Strategy interface {
	CreateObjects(path string, csvTable [][]string, params *pb.RpcObjectImportRequestCsvParams, progress process.Progress) (string, []*common.Snapshot, error)
}
