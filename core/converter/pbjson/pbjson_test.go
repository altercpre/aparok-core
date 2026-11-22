package pbjson

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/altercpre/aparok-core/core/block/editor/state"
	"github.com/altercpre/aparok-core/core/block/editor/template"
	"github.com/altercpre/aparok-core/pkg/lib/pb/model"
)

func TestPbc_Convert(t *testing.T) {
	s := state.NewDoc("root", nil).(*state.State)
	template.InitTemplate(s, template.WithTitle)
	c := NewConverter(s)
	result := c.Convert(model.SmartBlockType_Page)
	assert.NotEmpty(t, result)
}
