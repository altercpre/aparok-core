package objectid

import (
	"context"
	"strings"
	"time"

	"github.com/altercpre/aparok-sync/commonspace/object/tree/treestorage"

	"github.com/altercpre/aparok-core/core/block/import/common"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/core/domain/objectorigin"
	"github.com/altercpre/aparok-core/pkg/lib/bundle"
)

type participant struct{}

func newParticipant() *participant {
	return &participant{}
}

func (w *participant) GetIDAndPayload(ctx context.Context, spaceID string, sn *common.Snapshot, _ time.Time, _ bool, _ objectorigin.ObjectOrigin) (string, treestorage.TreeStorageCreatePayload, error) {
	participantId := sn.Snapshot.Data.Details.GetString(bundle.RelationKeyId)
	splitId := strings.Split(participantId, "_")
	identity := splitId[len(splitId)-1]
	newParticipantID := domain.NewParticipantId(spaceID, identity)
	return newParticipantID, treestorage.TreeStorageCreatePayload{}, nil
}
