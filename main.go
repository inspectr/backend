package main

import (
	"github.com/inspectr/backend/cmd"
	_ "github.com/inspectr/backend/plugins"
	_ "github.com/inspectr/backend/plugins/api"
	_ "github.com/inspectr/backend/plugins/heartbeat"
)

func main() {
	cmd.Execute()
}
