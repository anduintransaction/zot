package blobpresign

import (
	"net/http"

	"github.com/gorilla/mux"
	"zotregistry.dev/zot/pkg/api/config"
	"zotregistry.dev/zot/pkg/log"
	mTypes "zotregistry.dev/zot/pkg/meta/types"
	"zotregistry.dev/zot/pkg/storage"
)

type BlobPresign struct {
	Config          *config.Config
	StoreController storage.StoreController
	MetaDB          mTypes.MetaDB
	Log             log.Logger
}

func (bp *BlobPresign) GeneratePresignLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name, ok := vars["name"]
	if !ok || name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	digestStr, ok := vars["digest"]
	if !ok || digestStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bp.Log.Info().Msgf("generating presign link for module `%s`, digest `%s`", name, digestStr)

	w.WriteHeader(http.StatusOK)
}
