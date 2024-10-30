package config

func (c *Config) IsBlobPresignEnabled() bool {
	return c.Extensions != nil && c.Extensions.BlobPresign != nil && *c.Extensions.BlobPresign.Enable
}
