//go:build !anydebug

package debug

import "github.com/altercpre/aparok-sync/app"

func (d *debug) initHandlers(a *app.App) {
	// no-op
}
