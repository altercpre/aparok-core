package relationutils

import (
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/pkg/lib/bundle"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

func OptionFromDetails(det *domain.Details) *Option {
	return &Option{
		RelationOption: &model.RelationOption{
			Id:          det.GetString(bundle.RelationKeyId),
			Text:        det.GetString(bundle.RelationKeyName),
			Color:       det.GetString(bundle.RelationKeyRelationOptionColor),
			RelationKey: det.GetString(bundle.RelationKeyRelationKey),
			OrderId:     det.GetString(bundle.RelationKeyOrderId),
		},
	}
}

type Option struct {
	*model.RelationOption
}
