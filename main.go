package main

import (
	"storage/app/infra/config"
	"storage/cmd"
)

//
var (
	VERSION string
	COMMIT  string
	BUILD   string
)

func main() {
	config.C.Info.Version = VERSION
	config.C.Info.Commit = COMMIT
	config.C.Info.Build = BUILD

	cmd.Execute()
}
