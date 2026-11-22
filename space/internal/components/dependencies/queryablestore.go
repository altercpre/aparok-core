package dependencies

import "github.com/altercpre/aparok-core/pkg/lib/database"

type QueryableStore interface {
	Query(q database.Query) (records []database.Record, err error)
}
