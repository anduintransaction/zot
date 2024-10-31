package s3

import "github.com/distribution/distribution/v3/registry/storage/driver"

func (driver *Driver) GetDistStore() driver.StorageDriver {
	return driver.store
}
