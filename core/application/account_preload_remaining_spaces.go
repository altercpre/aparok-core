package application

import (
	"context"

	"github.com/altercpre/aparok-sync/app"

	"github.com/altercpre/aparok-core/space"
)

func (s *Service) AccountPreloadRemainingSpaces(ctx context.Context) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.app == nil {
		return ErrApplicationIsNotRunning
	}
	spaceService := app.MustComponent[space.Service](s.app)
	return spaceService.PreloadRemainingSpaces(ctx)
}
