package main

import (
	"github.com/ereztaiar/agenter/pkg/agents"
	"github.com/ereztaiar/agenter/pkg/commands"
)

var (
	rootConfig *agent.RootConfig
)

func main() {
	commands.Start(rootConfig)
}
