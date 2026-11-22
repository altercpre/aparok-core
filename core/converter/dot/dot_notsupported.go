//go:build gomobile || windows || nographviz || ignore || !cgo
// +build gomobile windows nographviz ignore !cgo

package dot

import (
	"github.com/altercpre/aparok-core/core/block/editor/smartblock"
	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/core/converter"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/core/relationutils"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
	"github.com/altercpre/aparok-core/space/spacecore/typeprovider"
)

func NewMultiConverter(format int, _ typeprovider.SmartBlockTypeProvider) converter.MultiConverter {
	return &dot{}
}

const (
	ExportFormatDOT = 0
	ExportFormatSVG = 1
)

type edgeType int

const (
	EdgeTypeRelation edgeType = iota
	EdgeTypeLink
)

type dot struct {
}

func (d *dot) SetKnownDocs(docs map[string]*domain.Details) converter.Converter {
	return d
}

func (d *dot) Add(space smartblock.Space, st *state.State, fetcher relationutils.RelationFormatFetcher) error {
	return nil
}

func (d *dot) Convert(sbType model.SmartBlockType) []byte {
	panic("not supported on windows")
	return nil
}

func (d *dot) Ext() string {
	return ""
}
