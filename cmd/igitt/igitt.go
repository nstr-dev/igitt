package main

import (
	"github.com/nstr-dev/igitt/internal/utilities/initialize"
)

func main() {

	var (
		version   = "dev"
		commit    = "none"
		buildDate = "unknown"
	)

	initialize.InitializeIgitt(version, commit, buildDate)
}
