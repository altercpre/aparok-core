package dependencies

import "github.com/altercpre/aparok-core/space/clientspace"

type BuiltinTemplateService interface {
	RegisterBuiltinTemplates(space clientspace.Space) error
}
