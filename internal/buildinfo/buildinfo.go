package buildinfo

import (
	"github.com/maxpeterkaya/peanut/common"
	"time"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func init() {
	if Version == "dev" {
		Commit = common.SHAHash(time.Now().String())
	}
}
