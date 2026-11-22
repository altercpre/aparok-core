package testutil

import (
	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/pb"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
	"github.com/altercpre/aparok-core/tests/blockbuilder"
)

func BuildStateFromAST(root *blockbuilder.Block) *state.State {
	st, err := state.NewDocFromSnapshot("", &pb.ChangeSnapshot{
		Data: &model.SmartBlockSnapshotBase{
			Blocks: root.Build(),
		},
	})
	if err != nil {
		panic(err)
	}
	_, _, err = state.ApplyState("", st, true)
	if err != nil {
		panic(err)
	}
	return st.NewState()
}
