// Package filterstring preserves Heart's former import path while the parser
// is owned and implemented by github.com/altercpre/aparok-block.
package filterstring

import external "github.com/altercpre/aparok-block/codec/anyblockjson/filterstring"

const EBNF = external.EBNF

type (
	Options = external.Options
	Error   = external.Error
)

var (
	Examples  = external.Examples
	Parse     = external.Parse
	ParseDate = external.ParseDate
	IsBareKey = external.IsBareKey
)
