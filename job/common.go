package job

import (
	"github.com/margostino/climateline-processor/config"
	"os"
)

var urls []*config.UrlConfig

// GetBaseJobUrl Rather than a global and one-time assigment, this method is convenient for overriding on testing
func GetBaseJobUrl() string {
	return os.Getenv("JOB_BASE_URL")
}
