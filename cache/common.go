package cache

import (
	"os"
)

// GetBaseCacheUrl Rather than a global and one-time assigment, this method is convenient for overriding on testing
func GetBaseCacheUrl() string {
	return os.Getenv("CACHE_BASE_URL")
}
