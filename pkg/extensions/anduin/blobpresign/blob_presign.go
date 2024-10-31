package blobpresign

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	godigest "github.com/opencontainers/go-digest"
	"zotregistry.dev/zot/pkg/api/config"
	zcommon "zotregistry.dev/zot/pkg/common"
	"zotregistry.dev/zot/pkg/log"
	mTypes "zotregistry.dev/zot/pkg/meta/types"
	"zotregistry.dev/zot/pkg/storage"
	"zotregistry.dev/zot/pkg/storage/imagestore"
	storageS3 "zotregistry.dev/zot/pkg/storage/s3"
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

	var lockLatency time.Time
	is.RLock(&lockLatency)
	defer is.RUnlock(&lockLatency)

	internalIS, ok := is.(*imagestore.ImageStore)
	if !ok {
		bp.Log.Error().Msgf("unexpected image store class. Expecing `*imagestore.ImageStore`, got `%s`", reflect.TypeOf(is).String())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	driver := internalIS.GetStorageDriver()
	internalDriver, ok := driver.(*storageS3.Driver)
	if !ok {
		bp.Log.Error().Msgf("unexpected driver class. Expecing `*s3.Driver`, got `%s`", reflect.TypeOf(driver).String())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	blobPath := is.BlobPath(name, digest)
	distDriver := internalDriver.GetDistStore()
	presignLink, err := distDriver.RedirectURL(r, blobPath)
	if err != nil {
		bp.Log.Err(err).Msgf("unable to generate presign link for blob path: %s", blobPath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zcommon.WriteJSON(w, http.StatusOK, map[string]string{
		"presign": presignLink,
	})
}
