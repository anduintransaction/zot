package imagestore

import (
	"net/http"
	"time"

	godigest "github.com/opencontainers/go-digest"
)

func (is *ImageStore) GenPresignLink(r *http.Request, repo string, digest godigest.Digest) (string, error) {
	var lockLatency time.Time

	if err := digest.Validate(); err != nil {
		return "", err
	}

	is.RLock(&lockLatency)
	defer is.RUnlock(&lockLatency)

	binfo, err := is.originalBlobInfo(repo, digest)
	if err != nil {
		return "", err
	}

	return is.storeDriver.RedirectURL(r, binfo.Path())
}
