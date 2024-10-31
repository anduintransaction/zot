package imagestore

import storageTypes "zotregistry.dev/zot/pkg/storage/types"

func (is *ImageStore) GetStorageDriver() storageTypes.Driver {
	return is.storeDriver
}
