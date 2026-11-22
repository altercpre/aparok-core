package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/altercpre/aparok-core/core/block/simple/test"
	"github.com/altercpre/aparok-core/pb"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

func TestDiff(t *testing.T) {
	testBlock := func() *rowBlock {
		return NewRowBlock(&model.Block{
			Restrictions: &model.BlockRestrictions{},
			Content:      &model.BlockContentOfTableRow{TableRow: &model.BlockContentTableRow{}},
		}).(*rowBlock)
	}
	t.Run("change header", func(t *testing.T) {
		// given
		b1 := testBlock()
		b2 := testBlock()

		// when
		b2.content.IsHeader = true
		diff, err := b1.Diff("", b2)

		// then
		require.NoError(t, err)
		require.Len(t, diff, 1)

		assert.Equal(t, test.MakeEvent(&pb.EventMessageValueOfBlockSetTableRow{
			BlockSetTableRow: &pb.EventBlockSetTableRow{
				Id:       b1.Id,
				IsHeader: &pb.EventBlockSetTableRowIsHeader{Value: true},
			},
		}), diff)
	})
}
