package editor

import (
	"github.com/altercpre/aparok-core/core/block/editor/smartblock"
	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/core/files/fileobject"
)

func migrateFilesToObjects(sb smartblock.SmartBlock, fileObjectService fileobject.Service) func(s *state.State) {
	return func(st *state.State) {
		fileObjectService.MigrateFileIdsInBlocks(st, sb.Space())
	}
}
