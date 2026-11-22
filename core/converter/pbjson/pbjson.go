package pbjson

import (
	"github.com/gogo/protobuf/jsonpb"

	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/core/converter"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/pb"
	"github.com/altercpre/aparok-core/pkg/lib/logging"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

var log = logging.Logger("json-converter")

func NewConverter(s state.Doc) converter.Converter {
	return &pbj{s: s}
}

type pbj struct {
	s state.Doc
}

func (p *pbj) Convert(sbType model.SmartBlockType) []byte {
	st := p.s.NewState()
	snapshot := &pb.ChangeSnapshot{
		Data: &model.SmartBlockSnapshotBase{
			Blocks:      st.BlocksToSave(),
			Details:     st.CombinedDetails().ToProto(),
			ObjectTypes: domain.MarshalTypeKeys(st.ObjectTypeKeys()),
			Collections: st.Store(),
			Key:         p.s.UniqueKeyInternal(),
			FileInfo:    st.GetFileInfo().ToModel(),
		},
	}
	mo := &pb.SnapshotWithType{
		SbType:   sbType,
		Snapshot: snapshot,
	}
	m := jsonpb.Marshaler{Indent: " "}
	result, err := m.MarshalToString(mo)
	if err != nil {
		log.Errorf("failed to convert object to json: %s", err)
	}
	return []byte(result)
}

func (p *pbj) Ext() string {
	return ".pb.json"
}

func (p *pbj) SetKnownDocs(map[string]*domain.Details) converter.Converter {
	return p
}

