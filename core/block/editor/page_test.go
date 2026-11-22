package editor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/altercpre/aparok-core/core/block/editor/smartblock"
	"github.com/altercpre/aparok-core/core/block/editor/smartblock/smarttest"
	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/core/block/migration"
	"github.com/altercpre/aparok-core/core/domain"
	"github.com/altercpre/aparok-core/pkg/lib/bundle"
	"github.com/altercpre/aparok-core/pkg/lib/localstore/objectstore"
	"github.com/altercpre/aparok-core/pkg/lib/localstore/objectstore/spaceindex"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

func TestPage_CreationStateMigration(t *testing.T) {
	t.Run("todo layout migration", func(t *testing.T) {
		// given

		initContext := &smartblock.InitContext{
			IsNewObject:    true,
			State:          state.NewDoc("rootId", nil).(*state.State),
			ObjectTypeKeys: []domain.TypeKey{bundle.TypeKeyTask},
		}

		storeFixture := spaceindex.NewStoreFixture(t)
		storeFixture.AddObjects(t, []objectstore.TestObject{
			{
				bundle.RelationKeyId:                domain.String("task"),
				bundle.RelationKeyUniqueKey:         domain.String(bundle.TypeKeyTask.URL()),
				bundle.RelationKeyRecommendedLayout: domain.Int64(int64(model.ObjectType_todo)),
				bundle.RelationKeySpaceId:           domain.String("spaceId"),
			},
		})

		p := &Page{}
		test := smarttest.New("rootId")
		test.SetSpaceId("spaceId")
		p.SmartBlock = test
		p.objectStore = storeFixture

		// when
		migration.RunMigrations(p, initContext)

		// then
		done := initContext.State.Details().GetBool(bundle.RelationKeyDone)
		assert.False(t, done)
	})
}
