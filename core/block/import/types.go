package importer

import (
	"context"

	"github.com/altercpre/aparok-sync/app"

	"github.com/altercpre/aparok-core/core/block"
	"github.com/altercpre/aparok-core/core/block/import/common"
	creator "github.com/altercpre/aparok-core/core/block/import/common/objectcreator"
	"github.com/altercpre/aparok-core/core/block/import/common/objectid"
	_ "github.com/altercpre/aparok-core/core/block/import/markdown"
	_ "github.com/altercpre/aparok-core/core/block/import/pb"
	_ "github.com/altercpre/aparok-core/core/block/import/web"
	"github.com/altercpre/aparok-core/core/block/process"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/core/domain/objectorigin"
	"github.com/altercpre/aparok-core/core/event"
	"github.com/altercpre/aparok-core/core/files/filesync"
	"github.com/altercpre/aparok-core/core/notifications"
	"github.com/altercpre/aparok-core/pb"
	"github.com/altercpre/aparok-core/pkg/lib/localstore/objectstore"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
	"github.com/altercpre/aparok-core/space"
)

type ImportRequest struct {
	*pb.RpcObjectImportRequest
	Origin           objectorigin.ObjectOrigin
	Progress         process.Progress
	SendNotification bool
	IsSync           bool
}

type ImportResponse struct {
	RootCollectionId string
	RootWidgetLayout model.BlockContentWidgetLayout
	ProcessId        string
	ObjectsCount     int64
	Err              error
}

type Dependencies struct {
	converters          map[string]common.Converter
	blockService        *block.Service
	objectCreator       creator.Service
	idProvider          objectid.IdAndKeyProvider
	fileSync            filesync.FileSync
	notificationService notifications.Notifications
	eventSender         event.Sender
	objectStore         objectstore.ObjectStore
	spaceService        space.Service
}

// Importer encapsulate logic with import
type Importer interface {
	app.Component
	Import(ctx context.Context, importRequest *ImportRequest) *ImportResponse

	ListImports(req *pb.RpcObjectImportListRequest) ([]*pb.RpcObjectImportListImportResponse, error)
	ImportWeb(ctx context.Context, req *ImportRequest) (string, *domain.Details, error)
	// nolint: lll
	ValidateNotionToken(ctx context.Context, req *pb.RpcObjectImportNotionValidateTokenRequest) (pb.RpcObjectImportNotionValidateTokenResponseErrorCode, error)
}
