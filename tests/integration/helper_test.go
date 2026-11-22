package integration

import (
	"golang.org/x/exp/constraints"

	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/pkg/lib/database"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

func filterEqualsToString(key domain.RelationKey, value string) database.FilterRequest {
	return database.FilterRequest{
		RelationKey: key,
		Condition:   model.BlockContentDataviewFilter_Equal,
		Value:       domain.String(value),
	}
}

func filterNotEmpty(key domain.RelationKey) database.FilterRequest {
	return database.FilterRequest{
		RelationKey: key,
		Condition:   model.BlockContentDataviewFilter_NotEmpty,
	}
}

func filterEqualsToInteger[T constraints.Integer](key domain.RelationKey, value T) database.FilterRequest {
	return database.FilterRequest{
		RelationKey: key,
		Condition:   model.BlockContentDataviewFilter_Equal,
		Value:       domain.Int64(value),
	}
}
