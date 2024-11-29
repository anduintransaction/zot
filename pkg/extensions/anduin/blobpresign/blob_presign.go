package blobpresign

import (
	"net/http"

	"github.com/gorilla/mux"
	godigest "github.com/opencontainers/go-digest"
	"zotregistry.dev/zot/pkg/api/config"
	zcommon "zotregistry.dev/zot/pkg/common"
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
	digest := godigest.Digest(digestStr)

	bp.Log.Info().Msgf("generating presign link for module `%s`, digest `%s`", name, digestStr)

	is := bp.StoreController.GetImageStore(name)
	presignLink, err := is.GenPresignLink(r, name, digest)
	if err != nil {
		bp.Log.Err(err).Msgf("unable to generate presign link for digest: %s", digestStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zcommon.WriteJSON(w, http.StatusOK, map[string]string{
		"presign": presignLink,
	})
}
