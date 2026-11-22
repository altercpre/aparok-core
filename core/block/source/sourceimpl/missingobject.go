package sourceimpl

import (
	"context"

	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/core/block/source"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/pb"
	"github.com/altercpre/aparok-core/pkg/lib/bundle"
	"github.com/altercpre/aparok-core/pkg/lib/core/smartblock"
	"github.com/altercpre/aparok-core/pkg/lib/localstore/addr"
)

type missingObject struct {
}

func NewMissingObject() (s source.Source) {
	return &missingObject{}
}

func (m *missingObject) ListIds() ([]string, error) {
	return []string{addr.MissingObject}, nil
}

func (m *missingObject) ReadOnly() bool {
	return true
}

// nolint:revive
func (m *missingObject) Id() string {
	return addr.MissingObject
}

func (m *missingObject) SpaceID() string {
	return addr.AnytypeMarketplaceWorkspace
}

func (m *missingObject) Type() smartblock.SmartBlockType {
	return smartblock.SmartBlockTypeMissingObject
}

func (m *missingObject) getDetails() (p *domain.Details) {
	det := domain.NewDetails()
	det.SetString(bundle.RelationKeyId, addr.MissingObject)
	det.SetBool(bundle.RelationKeyIsDeleted, true)
	det.SetBool(bundle.RelationKeyIsReadonly, true)
	det.SetBool(bundle.RelationKeyIsHidden, true)
	return det
}

func (m *missingObject) ReadDoc(ctx context.Context, receiver source.ChangeReceiver, empty bool) (doc state.Doc, err error) {
	s := state.NewDoc(addr.MissingObject, nil).(*state.State)

	d := m.getDetails()

	s.SetDetails(d)

	return s, nil
}

func (m *missingObject) Close() (err error) {
	return
}

func (m *missingObject) Heads() []string {
	return []string{"todo"}
}

func (m *missingObject) GetFileKeysSnapshot() []*pb.ChangeFileKeys {
	return nil
}

func (m *missingObject) PushChange(params source.PushChangeParams) (id string, err error) {
	return
}

func (m *missingObject) GetCreationInfo() (creatorObjectId string, createdDate int64, err error) {
	return addr.AnytypeProfileId, 0, nil
}
