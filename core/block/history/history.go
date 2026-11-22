package history

import (
	"fmt"

	"github.com/altercpre/aparok-core/core/block/editor/smartblock"
	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/core/block/migration"
)

func ResetToVersion(sb smartblock.SmartBlock, st *state.State) error {
	if err := sb.ResetToVersion(st); err != nil {
		return fmt.Errorf("resetting smartblock to version: %w", err)
	}
	migration.RunMigrations(sb, &smartblock.InitContext{
		State:   sb.NewState(),
		SpaceID: sb.SpaceID(),
	})
	return nil
}
