package dependencies

import (
	"context"

	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/space/clientspace"
)

type BundledObjectsInstaller interface {
	InstallBundledObjects(ctx context.Context, spc clientspace.Space, ids []string) ([]string, []*domain.Details, error)
	BundledObjectsIdsToInstall(ctx context.Context, spc clientspace.Space, sourceObjectIds []string) (ids domain.BundledObjectIds, err error)
}
