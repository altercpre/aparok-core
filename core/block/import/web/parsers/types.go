package parsers

import (
	"github.com/altercpre/aparok-core/core/block/import/common"
)

type RegisterParser func() Parser

var Parsers []RegisterParser

func RegisterFunc(p RegisterParser) {
	Parsers = append(Parsers, p)
}

type Parser interface {
	ParseUrl(url string) (*common.StateSnapshot, error)
	MatchUrl(url string) bool
}
