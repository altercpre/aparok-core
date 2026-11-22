package spaceinfo

import (
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

type ParticipantAclInfo struct {
	Id          string
	SpaceId     string
	Identity    string
	Permissions model.ParticipantPermissions
	Status      model.ParticipantStatus
}

type OneToOneParticipantData struct {
	Identity           string
	RequestMetadataKey []byte // a.k.a RequestMetadata, symKey
}
