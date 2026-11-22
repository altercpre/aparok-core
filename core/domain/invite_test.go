package domain

import (
	"testing"

	"github.com/altercpre/aparok-sync/commonspace/object/acl/list"
	"github.com/stretchr/testify/assert"

	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

func TestConvertParticipantPermissions_Admin(t *testing.T) {
	assert.Equal(t, list.AclPermissionsAdmin, ConvertParticipantPermissions(model.ParticipantPermissions_Admin))
}

func TestConvertAclPermissions_Admin(t *testing.T) {
	assert.Equal(t, model.ParticipantPermissions_Admin, ConvertAclPermissions(list.AclPermissionsAdmin))
}
