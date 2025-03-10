package imagestore

import (
	"net/http"

	godigest "github.com/opencontainers/go-digest"
)

func (is *ImageStore) GenPresignLink(r *http.Request, repo string, digest godigest.Digest) (string, error) {
	if err := digest.Validate(); err != nil {
		return "", err
	}

	presignedLink := ""
	err := is.lock.WithRepoReadLock(repo, func() error {
		binfo, err := is.originalBlobInfo(repo, digest)
		if err != nil {
			return err
		}

		presignedLink, err = is.storeDriver.RedirectURL(r, binfo.Path())
		return err
	})
	return presignedLink, err
}
