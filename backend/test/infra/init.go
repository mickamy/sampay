package infra

import (
	"os"
)

var (
	useTestContainers = false
	reuseContainer    = false
)

func init() {
	val, found := os.LookupEnv("TEST_CONTAINERS")
	if found {
		useTestContainers = val == "1"
	}
	val, found = os.LookupEnv("REUSE_CONTAINERS")
	if found {
		reuseContainer = val == "1"
	}
	if !useTestContainers && reuseContainer {
		panic("TEST_CONTAINERS=0 and REUSE_CONTAINERS=1 is not supported")
	}
}
