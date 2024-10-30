package extensions

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"zotregistry.dev/zot/pkg/api/config"
	"zotregistry.dev/zot/pkg/api/constants"
	zcommon "zotregistry.dev/zot/pkg/common"
	"zotregistry.dev/zot/pkg/extensions/anduin/blobpresign"
	"zotregistry.dev/zot/pkg/log"
	mTypes "zotregistry.dev/zot/pkg/meta/types"
	zreg "zotregistry.dev/zot/pkg/regexp"
	"zotregistry.dev/zot/pkg/storage"
)

func SetupBlobPresignRoutes(
	conf *config.Config,
	router *mux.Router,
	proxyRouter func(http.HandlerFunc) http.HandlerFunc,
	storeController storage.StoreController,
	metaDB mTypes.MetaDB,
	log log.Logger,
) {
	if !conf.IsBlobPresignEnabled() {
		log.Info().Msg("skip enabling the blob presign route as the config prerequisites are not met")
		return
	}
	log.Info().Msg("setting up blob presign routes")

	allowedMethods := zcommon.AllowedMethods(http.MethodGet, http.MethodPost)
	bp := blobpresign.BlobPresign{
		Config:          conf,
		StoreController: storeController,
		MetaDB:          metaDB,
		Log:             log,
	}

	extRouter := router.PathPrefix(constants.ExtBlobPresignPrefix).Subrouter()
	extRouter.Use(zcommon.CORSHeadersMiddleware(conf.HTTP.AllowOrigin))
	extRouter.Use(zcommon.ACHeadersMiddleware(conf, allowedMethods...))
	extRouter.Use(zcommon.AddExtensionSecurityHeaders())

	extRouter.
		HandleFunc(
			fmt.Sprintf("/{name:%s}/blobs/{digest}", zreg.NameRegexp.String()),
			proxyRouter(bp.GeneratePresignLink),
		).
		Methods(http.MethodGet)

	log.Info().Msg("finished setting up blob presign routes")
}
