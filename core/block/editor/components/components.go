package components

import "github.com/altercpre/aparok-core/core/domain"

type SyncStatusHandler interface {
	HandleSyncStatusUpdate(heads []string, status domain.ObjectSyncStatus, syncError domain.SyncError)
}
