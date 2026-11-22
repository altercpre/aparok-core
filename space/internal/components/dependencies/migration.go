package dependencies

import (
	"github.com/altercpre/aparok-core/pkg/lib/threads"
)

type MigrationService interface {
	RunMigrationsWhenIdle(spaceId string, derivedIDs threads.DerivedSmartblockIds)
}
