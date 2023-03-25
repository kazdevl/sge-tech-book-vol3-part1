package config

import "os"

func IsEnableCacheRepository() bool {
	return os.Getenv("ENABLE_CACHE_REPOSITORY") != "0"
}
