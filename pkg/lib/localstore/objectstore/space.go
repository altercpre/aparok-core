package objectstore

import (
	"fmt"

	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/pkg/lib/bundle"
	"github.com/altercpre/aparok-core/pkg/lib/database"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

type SpaceNameGetter interface {
	GetSpaceName(spaceId string) string
}

func (d *dsObjectStore) GetSpaceName(spaceId string) string {
	records, err := d.SpaceIndex(d.techSpaceId).Query(database.Query{
		Filters: []database.FilterRequest{
			{
				RelationKey: bundle.RelationKeyTargetSpaceId,
				Condition:   model.BlockContentDataviewFilter_Equal,
				Value:       domain.String(spaceId),
			},
			{
				RelationKey: bundle.RelationKeyResolvedLayout,
				Condition:   model.BlockContentDataviewFilter_Equal,
				Value:       domain.Int64(int64(model.ObjectType_spaceView)),
			},
		},
	})
	if err != nil {
		log.Errorf("failed to get details: %s", err)
	}
	var spaceName string
	if len(records) > 0 {
		spaceName = records[0].Details.GetString(bundle.RelationKeyName)
	}
	return spaceName
}

func (d *dsObjectStore) GetSpaceViewDetails(spaceId string) (*domain.Details, error) {
	records, err := d.SpaceIndex(d.techSpaceId).Query(database.Query{
		Filters: []database.FilterRequest{
			{
				RelationKey: bundle.RelationKeyTargetSpaceId,
				Condition:   model.BlockContentDataviewFilter_Equal,
				Value:       domain.String(spaceId),
			},
			{
				RelationKey: bundle.RelationKeyResolvedLayout,
				Condition:   model.BlockContentDataviewFilter_Equal,
				Value:       domain.Int64(int64(model.ObjectType_spaceView)),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(records) > 0 {
		return records[0].Details, nil
	}
	return nil, fmt.Errorf("spaceView not found")
}
