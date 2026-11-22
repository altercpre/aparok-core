package application

import (
	"github.com/altercpre/aparok-sync/app"

	"github.com/altercpre/aparok-core/space/spacecore/localdiscovery"
)

func (s *Service) EnableLocalNetworkSync() error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.app == nil {
		return ErrApplicationIsNotRunning
	}
	localDiscoveryService := app.MustComponent[localdiscovery.LocalDiscovery](s.app)
	return localDiscoveryService.Start()
}
