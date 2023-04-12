package main

import (
	"github.com/ythosa/gowrapper/internal/cli"
	"github.com/ythosa/gowrapper/internal/services/generator/gowrap"
)

func main() {
	//builder := options_builder.NewOptionsBuilder()
	////g.SetOutputFolder("tracing_wrappers")
	//builder.SetFilePostfix("tracing")
	//builder.SetTemplate("/Users/rbabin/Projects/gowrapper/internal/templates/opentracing")
	//
	//g := gowrap.NewGenerator(finder.NewInterfacesFinder(), builder)
	//g.SetExcludedDirs([]string{"internal/pkg", "internal/usecase/issue_meta", "internal/repository/issue_meta", "internal/pb", "internal/enum", "internal/model", "auditor/table"})
	//g.SetExcludedFiles([]string{"/auditor/converter.go"})
	//g.SetInitialDirectory("/Users/rbabin/ozon/it-manager-backend/internal")
	//g.Generate()

	cmd := cli.NewRootCommand()
	cmd.RegisterCommands(cli.NewGenCommand(gowrap.NewFactory()))
	cmd.Execute()
}
