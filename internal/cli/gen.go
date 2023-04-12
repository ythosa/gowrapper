package cli

import (
	"github.com/spf13/cobra"
	"github.com/ythosa/gowrapper/internal/model"
	"github.com/ythosa/gowrapper/internal/services/generator"
)

type GeneratorFactory interface {
	Construct(options model.Options) generator.Generator
}

type GenCommand struct {
	generatorFactory GeneratorFactory

	options model.Options
}

func NewGenCommand(generatorFactory GeneratorFactory) *GenCommand {
	return &GenCommand{generatorFactory: generatorFactory}
}

func (cmd *GenCommand) descriptor() *cobra.Command {
	cc := &cobra.Command{
		Use:   "gen",
		Short: "Generates wrappers",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.generatorFactory.Construct(cmd.options).Generate()
		},
	}

	cc.Flags().StringVarP(&cmd.options.OutputFolder,
		"out-folder", "o", "", "folder name where to generate wraps")
	cc.Flags().StringVarP(&cmd.options.FilePostfix,
		"postfix", "p", "", "file postfix for wraps files")
	cc.Flags().StringVarP(&cmd.options.TemplatePath,
		"template", "t", "", "full path for template")
	cc.Flags().StringSliceVar(&cmd.options.ExcludedDirs,
		"excluded-dirs", nil, "parts of paths of excluded dirs")
	cc.Flags().StringSliceVar(&cmd.options.ExcludedFiles,
		"excluded-files", nil, "parts of paths of excluded files")
	cc.Flags().StringVarP(&cmd.options.InitialDirectory,
		"dir", "d", "", "path of initial directory")

	return cc
}
