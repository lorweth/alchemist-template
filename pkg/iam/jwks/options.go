package jwks

import (
	"github.com/virsavik/alchemist-template/pkg/logger"
)

type Option func(provider *CacheProvider)

func WithLogger(log logger.Logger) Option {
	return func(provider *CacheProvider) {
		provider.logger = log
	}
}
