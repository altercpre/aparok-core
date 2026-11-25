package linkpreview

import (
	"context"
	"errors"
	neturl "net/url"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

const (
	CName = "linkpreview"
)

var (
	log            = logging.Logger(CName)
	ErrPrivateLink = errors.New("private link")
)

func New() LinkPreview {
	return &linkPreview{}
}

type LinkPreview interface {
	Fetch(ctx context.Context, url string, withResponseBody bool) (linkPreview model.LinkPreview, responseBody []byte, isFile bool, err error)
	app.Component
}

type linkPreview struct{}

func (l *linkPreview) Init(_ *app.App) error {
	return nil
}

func (l *linkPreview) Name() string {
	return CName
}

func (l *linkPreview) Fetch(
	ctx context.Context, fetchUrl string, withResponseBody bool,
) (linkPreview model.LinkPreview, responseBody []byte, isFile bool, err error) {
	title := fetchUrl
	if u, err := neturl.Parse(fetchUrl); err == nil && u.Hostname() != "" {
		title = u.Hostname()
	}
	return model.LinkPreview{
		Url:   fetchUrl,
		Title: title,
	}, nil, false, nil
}
