package anytype

import (
	"context"
	"os"
	"regexp"
	"time"

	"github.com/altercpre/aparok-sync/app"
	"github.com/altercpre/aparok-sync/app/debugstat"
	"github.com/altercpre/aparok-sync/commonfile/fileservice"
	"github.com/altercpre/aparok-sync/commonspace"
	"github.com/altercpre/aparok-sync/commonspace/acl/aclclient"
	anysyncpubsub "github.com/altercpre/aparok-sync/commonspace/pubsub"
	anysyncinboxclient "github.com/altercpre/aparok-sync/coordinator/inboxclient"

	"github.com/altercpre/aparok-sync/coordinator/nodeconfsource"
	"github.com/altercpre/aparok-sync/coordinator/subscribeclient"
	"github.com/altercpre/aparok-sync/metric"
	"github.com/altercpre/aparok-sync/net/peerservice"
	"github.com/altercpre/aparok-sync/net/pool"
	"github.com/altercpre/aparok-sync/net/quicdemotion"
	"github.com/altercpre/aparok-sync/net/rpc/debugserver"
	"github.com/altercpre/aparok-sync/net/rpc/server"
	"github.com/altercpre/aparok-sync/net/secureservice"
	"github.com/altercpre/aparok-sync/net/streampool"
	"github.com/altercpre/aparok-sync/net/transport/quic"
	"github.com/altercpre/aparok-sync/net/transport/yamux"
	"github.com/altercpre/aparok-sync/node/nodeclient"
	"github.com/altercpre/aparok-sync/nodeconf"
	"github.com/altercpre/aparok-sync/nodeconf/nodeconfstore"
	"github.com/altercpre/aparok-sync/util/crypto"
	"github.com/altercpre/aparok-sync/util/syncqueues"
	"github.com/altercpre/anytype-publish-server/publishclient"
	"go.uber.org/zap"

	"github.com/altercpre/aparok-sync/nameservice/nameserviceclient"
	"github.com/altercpre/aparok-sync/paymentservice/paymentserviceclient"
	"github.com/altercpre/aparok-sync/paymentservice/paymentserviceclient2"

	"github.com/altercpre/aparok-core/core/acl"
	"github.com/altercpre/aparok-core/core/aparok/account"
	"github.com/altercpre/aparok-core/core/aparok/config"
	"github.com/altercpre/aparok-core/core/api"
	"github.com/altercpre/aparok-core/core/block"
	"github.com/altercpre/aparok-core/core/block/backlinks"
	"github.com/altercpre/aparok-core/core/block/bookmark"
	decorator "github.com/altercpre/aparok-core/core/block/bookmark/bookmarkimporter"
	"github.com/altercpre/aparok-core/core/block/chats"
	"github.com/altercpre/aparok-core/core/block/chats/chatrepository"
	"github.com/altercpre/aparok-core/core/block/chats/chatsubscription"
	"github.com/altercpre/aparok-core/core/block/collection"
	"github.com/altercpre/aparok-core/core/block/dataviewservice"
	"github.com/altercpre/aparok-core/core/block/deletionaudit"
	"github.com/altercpre/aparok-core/core/block/detailservice"
	"github.com/altercpre/aparok-core/core/block/editor"
	"github.com/altercpre/aparok-core/core/block/editor/converter"
	"github.com/altercpre/aparok-core/core/block/export"
	importer "github.com/altercpre/aparok-core/core/block/import"
	importv2adapter "github.com/altercpre/aparok-core/core/block/importv2/adapter"
	"github.com/altercpre/aparok-core/core/block/object/idderiver/idderiverimpl"
	"github.com/altercpre/aparok-core/core/block/object/idresolver"
	"github.com/altercpre/aparok-core/core/block/object/objectcreator"
	"github.com/altercpre/aparok-core/core/block/object/objectgraph"
	"github.com/altercpre/aparok-core/core/block/object/treemanager"
	"github.com/altercpre/aparok-core/core/block/objectgc"
	"github.com/altercpre/aparok-core/core/block/personalfavorites"
	"github.com/altercpre/aparok-core/core/block/process"
	"github.com/altercpre/aparok-core/core/block/source/sourceimpl"
	"github.com/altercpre/aparok-core/core/block/template/templateimpl"
	"github.com/altercpre/aparok-core/core/configfetcher"
	"github.com/altercpre/aparok-core/core/debug"
	"github.com/altercpre/aparok-core/core/debug/debugreporter"
	"github.com/altercpre/aparok-core/core/debug/profiler"
	"github.com/altercpre/aparok-core/core/device"
	"github.com/altercpre/aparok-core/core/durability"
	"github.com/altercpre/aparok-core/core/files"
	"github.com/altercpre/aparok-core/core/files/fileacl"
	"github.com/altercpre/aparok-core/core/files/filedownloader"
	"github.com/altercpre/aparok-core/core/files/fileobject"
	"github.com/altercpre/aparok-core/core/files/fileoffloader"
	"github.com/altercpre/aparok-core/core/files/filespaceusage"
	"github.com/altercpre/aparok-core/core/files/filestorage"
	"github.com/altercpre/aparok-core/core/files/filestorage/rpcstore"
	"github.com/altercpre/aparok-core/core/files/filesync"
	"github.com/altercpre/aparok-core/core/files/fileuploader"
	"github.com/altercpre/aparok-core/core/files/reconciler"
	"github.com/altercpre/aparok-core/core/gallery"
	"github.com/altercpre/aparok-core/core/history"
	"github.com/altercpre/aparok-core/core/identity"
	"github.com/altercpre/aparok-core/core/inbox/inboxclient"
	"github.com/altercpre/aparok-core/core/inbox/inboxservice"
	"github.com/altercpre/aparok-core/core/indexer"
	"github.com/altercpre/aparok-core/core/invitecleanup"
	"github.com/altercpre/aparok-core/core/inviteservice"
	"github.com/altercpre/aparok-core/core/invitestore"
	"github.com/altercpre/aparok-core/core/kanban"
	"github.com/altercpre/aparok-core/core/migration"
	"github.com/altercpre/aparok-core/core/nameservice"
	"github.com/altercpre/aparok-core/core/notifications"
	"github.com/altercpre/aparok-core/core/order"
	"github.com/altercpre/aparok-core/core/payments"
	paymentscache "github.com/altercpre/aparok-core/core/payments/cache"
	"github.com/altercpre/aparok-core/core/payments/emailcollector"
	"github.com/altercpre/aparok-core/core/peerstatus"
	"github.com/altercpre/aparok-core/core/publish"
	"github.com/altercpre/aparok-core/core/pubsub"
	"github.com/altercpre/aparok-core/core/pushnotification"
	"github.com/altercpre/aparok-core/core/pushnotification/pushclient"
	"github.com/altercpre/aparok-core/core/relationutils/formatfetcher"
	"github.com/altercpre/aparok-core/core/session"
	"github.com/altercpre/aparok-core/core/subscription"
	"github.com/altercpre/aparok-core/core/subscription/crossspacesub"
	"github.com/altercpre/aparok-core/core/syncstatus"
	"github.com/altercpre/aparok-core/core/syncstatus/detailsupdater"
	"github.com/altercpre/aparok-core/core/syncstatus/nodestatus"
	"github.com/altercpre/aparok-core/core/syncstatus/spacesyncstatus"
	"github.com/altercpre/aparok-core/core/syncstatus/syncsubscriptions"
	"github.com/altercpre/aparok-core/core/wallet"
	"github.com/altercpre/aparok-core/metrics"
	"github.com/altercpre/aparok-core/net/peerobservermux"
	"github.com/altercpre/aparok-core/net/transportpenalty"
	"github.com/altercpre/aparok-core/pkg/lib/core"
	"github.com/altercpre/aparok-core/pkg/lib/datastore/anystoreprovider"
	"github.com/altercpre/aparok-core/pkg/lib/datastore/clientds"
	"github.com/altercpre/aparok-core/pkg/lib/gateway"
	"github.com/altercpre/aparok-core/pkg/lib/localstore/ftsearch"
	"github.com/altercpre/aparok-core/pkg/lib/localstore/objectstore"
	"github.com/altercpre/aparok-core/pkg/lib/logging"
	"github.com/altercpre/aparok-core/space"
	"github.com/altercpre/aparok-core/space/coordinatorclient"
	"github.com/altercpre/aparok-core/space/deletioncontroller"
	"github.com/altercpre/aparok-core/space/spacecore"
	"github.com/altercpre/aparok-core/space/spacecore/clientserver"
	"github.com/altercpre/aparok-core/space/spacecore/credentialprovider"
	"github.com/altercpre/aparok-core/space/spacecore/localdiscovery"
	"github.com/altercpre/aparok-core/space/spacecore/oldstorage"
	"github.com/altercpre/aparok-core/space/spacecore/peermanager"
	"github.com/altercpre/aparok-core/space/spacecore/peerstore"
	"github.com/altercpre/aparok-core/space/spacecore/storage"
	"github.com/altercpre/aparok-core/space/spacecore/storage/migrator"
	"github.com/altercpre/aparok-core/space/spacecore/storage/migratorfinisher"
	"github.com/altercpre/aparok-core/space/spacecore/typeprovider"
	"github.com/altercpre/aparok-core/space/spacefactory"
	"github.com/altercpre/aparok-core/space/virtualspaceservice"
	"github.com/altercpre/aparok-core/util/builtinobjects"
	"github.com/altercpre/aparok-core/util/builtintemplate"
	"github.com/altercpre/aparok-core/util/linkpreview"
	"github.com/altercpre/aparok-core/util/unsplash"
	"github.com/altercpre/aparok-core/util/vcs"
)

var (
	log          = logging.LoggerNotSugared("anytype-app")
	WarningAfter = time.Second * 1
)

func BootstrapConfig(newAccount bool, spaceStreamAutoJoinUrl string) *config.Config {
	return config.New(
		config.WithDebugAddr(os.Getenv("ANYTYPE_DEBUG_ADDR")),
		config.WithNewAccount(newAccount),
		config.WithAutoJoinStream(spaceStreamAutoJoinUrl),
	)
}

func BootstrapWallet(rootPath string, derivationResult crypto.DerivationResult, lang string) wallet.Wallet {
	return wallet.NewWithAccountRepo(rootPath, derivationResult, lang)
}

func StartNewApp(ctx context.Context, clientWithVersion string, components ...app.Component) (a *app.App, err error) {
	a = new(app.App)
	complexAppVersion := appVersion(a, clientWithVersion)
	a.SetVersionName(complexAppVersion)
	Bootstrap(a, components...)
	metrics.Service.SetAppVersion(a.VersionName())
	metrics.Service.Run()
	log.Info("starting app", zap.String("mwVersion", MiddlewareVersion()), zap.String("appVersion", complexAppVersion))
	startTime := time.Now()
	if err = a.Start(ctx); err != nil {
		metrics.Service.Close()
		a = nil
		return
	}
	// Wire the profiler component as the process-wide Reporter so the
	// metrics long-method interceptor and any future critical-signal sites
	// can produce snapshots + events without reaching into internal APIs.
	metrics.SetReporter(app.MustComponent[debugreporter.Reporter](a))
	totalSpent := time.Since(startTime)
	l := log.With(zap.Int64("total", totalSpent.Milliseconds()))
	stat := a.StartStat()

	if v, ok := ctx.Value(metrics.CtxKeyRPC).(string); ok {
		l = l.With(zap.String("rpc", v))
	}

	for comp, spent := range stat.SpentMsPerComp {
		if spent == 0 {
			continue
		}
		l = l.With(zap.Int64(comp, spent))
	}

	l.With(zap.Int64("totalRun", stat.SpentMsTotal))
	a.IterateComponents(func(comp app.Component) {
		if c, ok := comp.(ComponentLogFieldsGetter); ok {
			for _, field := range c.GetLogFields() {
				field.Key = comp.Name() + "_" + field.Key
				l = l.With(field)
			}
		}
	})

	if totalSpent > WarningAfter {
		l.Warn("app started")
	} else {
		l.Debug("app started")
	}
	return
}

func appVersion(a *app.App, clientWithVersion string) string {
	clientWithVersion = regexp.MustCompile(`(@|\/)+`).ReplaceAllString(clientWithVersion, "_")
	middleVersion := MiddlewareVersion()
	anySyncVersion := a.AnySyncVersion()
	return clientWithVersion + "/middle:" + middleVersion + "/any-sync:" + anySyncVersion
}

func BootstrapMigration(a *app.App, components ...app.Component) {
	for _, c := range components {
		a.Register(c)
	}
	a.Register(migratorfinisher.New()).
		Register(clientds.New()).
		Register(oldstorage.New()).
		Register(storage.New()).
		Register(process.New()).
		Register(migrator.New())
}

func Bootstrap(a *app.App, components ...app.Component) {
	for _, c := range components {
		a.Register(c)
	}

	// the pubsub engine takes its client-side deps (crypto, membership, peers)
	// from the heart-side component, resolved lazily after Init
	pubsubService := pubsub.New()

	a.
		// profiler is registered early so its Init wires the event sender
		// before any storage component runs its own Init — storage corruption
		// detection below relies on the Reporter being live.
		Register(profiler.New()).
		// Data storages
		Register(debugstat.New()).
		Register(anystoreprovider.New()).
		Register(ftsearch.TantivyNew()).
		Register(objectstore.New()).
		Register(backlinks.New()).
		// Services
		Register(collection.New()).
		Register(subscription.New()).
		Register(crossspacesub.New()).
		Register(formatfetcher.New()).
		Register(nodeconfsource.New()).
		Register(nodeconfstore.New()).
		Register(nodeconf.New()).
		Register(syncqueues.New()).
		// the single peerobserver slot, as a fan-out: the recovery status
		// tracker registers with it at Init, diagnostics can later
		Register(peerstore.New()).
		Register(peerobservermux.New()).
		Register(storage.New()).
		Register(secureservice.New()).
		Register(metric.New()).
		Register(server.New()).
		Register(debugserver.New()).
		Register(pool.New())

	if os.Getenv(transportpenalty.DisableEnv) != "0" {
		// quicdemotion is a standalone app component (any-sync net/quicdemotion):
		// registering it is what enables auto-demoting peers to yamux-first
		// when their QUIC connections keep dying under DPI-style degradation
		// (GO-7467). peerservice picks it up itself in its own Init, so it
		// must be registered before peerservice's Init runs.
		a.Register(quicdemotion.New())
	}

	a.
		Register(peerservice.New()).
		Register(yamux.New()).
		Register(quic.New()).
		Register(clientserver.New()).
		Register(coordinatorclient.New()).
		Register(nodeclient.New()).
		Register(credentialprovider.New()).
		Register(spacecore.NewStreamOpener()).
		Register(streampool.New()).
		Register(commonspace.New()).
		Register(aclclient.NewAclJoiningClient()).
		Register(virtualspaceservice.New()).
		Register(spacecore.New()).
		Register(idresolver.New(200*time.Millisecond, 2*time.Second)).
		Register(device.New()).
		Register(transportpenalty.New()).
		Register(localdiscovery.New()).
		Register(peermanager.New()).
		Register(anysyncpubsub.New(pubsubService.EngineDeps())).
		// Cancel accepted publishes before the engine waits for its workers.
		Register(pubsubService).
		Register(typeprovider.New()).
		Register(fileuploader.New()).
		Register(rpcstore.New()).
		Register(fileservice.New()).
		Register(filestorage.New()).
		Register(files.New()).
		Register(filespaceusage.New()).
		Register(fileoffloader.New()).
		Register(filedownloader.New()).
		Register(fileacl.New()).
		Register(chatrepository.New()).
		Register(chatsubscription.New()).
		Register(chats.New()).
		Register(personalfavorites.New()).
		Register(sourceimpl.New()).
		Register(spacefactory.New()).
		Register(space.New()).
		Register(idderiverimpl.New()).
		Register(deletioncontroller.New()).
		Register(invitestore.New()).
		Register(invitecleanup.New()).
		Register(filesync.New()).
		Register(reconciler.New()).
		Register(fileobject.New()).
		Register(inviteservice.New()).
		Register(publish.New()).
		Register(publishclient.New()).
		Register(acl.New()).
		Register(builtintemplate.New()).
		Register(converter.NewLayoutConverter()).
		Register(configfetcher.New()).
		Register(process.New()).
		Register(core.NewTempDirService()).
		Register(treemanager.New()).
		Register(block.New()).
		Register(detailservice.New()).
		Register(objectgc.New()).
		Register(deletionaudit.New()).
		Register(dataviewservice.New()).
		Register(migration.New()).
		Register(indexer.New()).
		Register(detailsupdater.New()).
		Register(session.NewHookRunner()).
		Register(spacesyncstatus.NewSpaceSyncStatus()).
		Register(nodestatus.NewNodeStatus()).
		Register(syncstatus.New()).
		Register(history.New()).
		Register(gateway.New()).
		Register(export.New()).
		Register(linkpreview.NewWithCache()).
		Register(unsplash.New()).
		Register(debug.New()).
		Register(syncsubscriptions.New()).
		Register(builtinobjects.New()).
		Register(gallery.New()).
		Register(bookmark.New()).
		Register(importer.New()).
		Register(importv2adapter.New()).
		Register(decorator.New()).
		Register(objectcreator.NewCreator()).
		Register(kanban.New()).
		Register(device.NewDevices()).
		Register(editor.NewObjectFactory()).
		Register(objectgraph.NewBuilder()).
		Register(account.New()).
		Register(identity.New(5*time.Minute, 10*time.Second)).
		Register(templateimpl.New()).
		Register(notifications.New(time.Second * 10)).
		Register(paymentserviceclient.New()).
		Register(paymentserviceclient2.New()).
		Register(nameservice.New()).
		Register(nameserviceclient.New()).
		Register(payments.New()).
		Register(paymentscache.New()).
		Register(emailcollector.New()).
		Register(peerstatus.New()).
		Register(order.New()).
		Register(api.New()).
		Register(pushclient.New()).
		Register(pushnotification.New()).
		Register(subscribeclient.New()).
		Register(anysyncinboxclient.New()).
		Register(inboxclient.New()).
		Register(inboxservice.New()).
		Register(durability.New()) // leave it the last one
}

func MiddlewareVersion() string {
	return vcs.GetVCSInfo().Version()
}

type ComponentLogFieldsGetter interface {
	// GetLogFields returns additional useful fields for logs to debug long app start/stop duration or something else in the future
	// You don't need to provide the component name in the field's Key, because it will be added automatically
	GetLogFields() []zap.Field
}
