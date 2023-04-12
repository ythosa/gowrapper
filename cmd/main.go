package main

import (
	"github.com/ythosa/gowrapper/internal/cli"
	"github.com/ythosa/gowrapper/internal/services/generator/gowrap"
)

func main() {
	cli.NewRootCommand().RegisterCommands(cli.NewGenCommand(gowrap.NewFactory())).Execute()
}
