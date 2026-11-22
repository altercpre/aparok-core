package converter

import (
	"github.com/altercpre/aparok-core/core/block/editor/smartblock"
	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/core/relationutils"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

type Converter interface {
	Convert(sbType model.SmartBlockType) (result []byte)
	SetKnownDocs(docs map[string]*domain.Details) Converter
	Ext() string
}

type MultiConverter interface {
	Converter
	Add(space smartblock.Space, state *state.State, fetcher relationutils.RelationFormatFetcher) error
}
