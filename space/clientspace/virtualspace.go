package clientspace

import (
	"context"
	"time"

	"github.com/altercpre/aparok-sync/accountservice"
	"github.com/altercpre/aparok-sync/commonspace"
	"github.com/altercpre/aparok-sync/commonspace/acl/aclclient"
	"github.com/altercpre/aparok-sync/commonspace/headsync"
	"github.com/altercpre/aparok-sync/commonspace/object/acl/syncacl"
	"github.com/altercpre/aparok-sync/commonspace/object/keyvalue/kvinterfaces"
	"github.com/altercpre/aparok-sync/commonspace/object/treesyncer"
	"github.com/altercpre/aparok-sync/commonspace/objecttreebuilder"
	"github.com/altercpre/aparok-sync/commonspace/spacestorage"
	"github.com/altercpre/aparok-sync/commonspace/spacesyncproto"
	"github.com/altercpre/aparok-sync/commonspace/sync/objectsync/objectmessages"
	"github.com/altercpre/aparok-sync/commonspace/syncstatus"
	"github.com/altercpre/aparok-sync/net/peer"
	"storj.io/drpc"

	"github.com/altercpre/aparok-core/core/block/object/objectcache"
	"github.com/altercpre/aparok-core/core/domain"
)

type VirtualSpaceDeps struct {
	ObjectFactory   objectcache.ObjectFactory
	AccountService  accountservice.Service
	PersonalSpaceId string
	Indexer         spaceIndexer
	Installer       bundledObjectsInstaller
	TypePrefix      string
	RelationPrefix  string
}

type VirtualSpace struct {
	*space
	TypePrefix, RelationPrefix string
}

func NewVirtualSpace(spaceId string, deps VirtualSpaceDeps) *VirtualSpace {
	vs := &VirtualSpace{
		space: &space{
			indexer:                deps.Indexer,
			installer:              deps.Installer,
			common:                 newVirtualCommonSpace(spaceId),
			loadMandatoryObjectsCh: make(chan struct{}),
			personalSpaceId:        deps.PersonalSpaceId,
		},
		TypePrefix:     deps.TypePrefix,
		RelationPrefix: deps.RelationPrefix,
	}
	vs.space.Cache = objectcache.New(deps.AccountService, deps.ObjectFactory, deps.PersonalSpaceId, vs)
	return vs
}

// IsReadOnly is always true: a virtual space has no ACL to ask, and the embedded space's check
// would dereference it
func (vs *VirtualSpace) IsReadOnly() bool {
	return true
}

func (vs *VirtualSpace) GetRelationIdByKey(ctx context.Context, key domain.RelationKey) (id string, err error) {
	return vs.RelationPrefix + key.String(), nil
}

func (vs *VirtualSpace) GetTypeIdByKey(ctx context.Context, key domain.TypeKey) (id string, err error) {
	return vs.TypePrefix + key.String(), nil
}

func newVirtualCommonSpace(spaceId string) commonspace.Space {
	return &virtualCommonSpace{spaceId: spaceId}
}

type virtualCommonSpace struct {
	spaceId string
}

func (c *virtualCommonSpace) HandleMessage(ctx context.Context, msg *objectmessages.HeadUpdate) (err error) {
	return nil
}

func (c *virtualCommonSpace) HandleStreamSyncRequest(ctx context.Context, req *spacesyncproto.ObjectSyncMessage, stream drpc.Stream) (err error) {
	return nil
}

func (c *virtualCommonSpace) HandleStream(stream spacesyncproto.DRPCSpaceSync_ObjectSyncStreamStream) error {
	return nil
}

func (c *virtualCommonSpace) AclClient() aclclient.AclSpaceClient {
	return nil
}

func (c *virtualCommonSpace) IsPersonal() bool {
	return false
}

func (c *virtualCommonSpace) Id() string {
	return c.spaceId
}

func (c *virtualCommonSpace) Init(ctx context.Context) error {
	return nil
}

func (c *virtualCommonSpace) Acl() syncacl.SyncAcl {
	return nil
}

func (c *virtualCommonSpace) StoredIds() []string {
	return nil
}

func (c *virtualCommonSpace) SyncHeads(ctx context.Context) error {
	return nil
}

func (c *virtualCommonSpace) DebugAllHeads() []headsync.TreeHeads {
	return nil
}

func (c *virtualCommonSpace) Description(ctx context.Context) (desc commonspace.SpaceDescription, err error) {
	return
}

func (c *virtualCommonSpace) TreeBuilder() objecttreebuilder.TreeBuilder {
	return nil
}

func (c *virtualCommonSpace) TreeSyncer() treesyncer.TreeSyncer {
	return nil
}

func (c *virtualCommonSpace) SyncStatus() syncstatus.StatusUpdater {
	return nil
}

func (c *virtualCommonSpace) Storage() spacestorage.SpaceStorage {
	return nil
}

func (c *virtualCommonSpace) DeleteTree(ctx context.Context, id string) (err error) {
	return nil
}

func (c *virtualCommonSpace) GetNodePeers(ctx context.Context) (peer []peer.Peer, err error) {
	return
}

func (c *virtualCommonSpace) HandleRangeRequest(ctx context.Context, req *spacesyncproto.HeadSyncRequest) (resp *spacesyncproto.HeadSyncResponse, err error) {
	return
}

func (c *virtualCommonSpace) TryClose(objectTTL time.Duration) (close bool, err error) {
	return
}

func (c *virtualCommonSpace) Close() error {
	return nil
}

func (c *virtualCommonSpace) IsReadOnly() bool {
	return false
}

func (c *virtualCommonSpace) KeyValue() kvinterfaces.KeyValueService {
	return nil
}
