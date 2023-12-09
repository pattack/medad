package loaders

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"

	"github.com/pattack/medad"
)

type localLoader struct {
}

func LocalLoader(opts ...tricks.Option[localLoader]) medad.Loader {
	l := localLoader{}
	tricks.ApplyOptions(&l, opts...)

	return l
}

func (c localLoader) Load(ctx context.Context, sources ...string) bricks.Bag[medad.Article] {
	// Name: filepath.Base(filepath.Dir(src))
	// LangCode: extensionless(filepath.Base(src))

	panic(bricks.ErrNotImplemented)
}

func extensionless(basename string) string {
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}
