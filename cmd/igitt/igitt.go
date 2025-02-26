package main

import (
	"github.com/nstr-dev/igitt/internal/utilities/initialize"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func main() {
	initialize.InitializeIgitt(Version, Commit, BuildDate)
}
