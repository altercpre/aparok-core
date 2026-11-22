package credentialprovider

import (
	"context"

	"github.com/altercpre/aparok-sync/app"
	"github.com/altercpre/aparok-sync/commonspace/credentialprovider"
	"github.com/altercpre/aparok-sync/commonspace/spacesyncproto"
	"github.com/altercpre/aparok-sync/coordinator/coordinatorclient"
	"github.com/gogo/protobuf/proto"
)

func New() app.Component {
	return &credentialProvider{}
}

type credentialProvider struct {
	client coordinatorclient.CoordinatorClient
}

func (c *credentialProvider) Init(a *app.App) (err error) {
	c.client = a.MustComponent(coordinatorclient.CName).(coordinatorclient.CoordinatorClient)
	return
}

func (c *credentialProvider) Name() (name string) {
	return credentialprovider.CName
}

func (c *credentialProvider) GetCredential(ctx context.Context, spaceHeader *spacesyncproto.RawSpaceHeaderWithId) ([]byte, error) {
	payload := coordinatorclient.SpaceSignPayload{
		SpaceId:     spaceHeader.Id,
		SpaceHeader: spaceHeader.RawHeader,
	}
	receipt, err := c.client.SpaceSign(ctx, payload)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(receipt)
}
