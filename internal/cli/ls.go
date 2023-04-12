package cli

import (
	"github.com/spf13/cobra"
	"github.com/ythosa/gowrapper/internal/model"
	"github.com/ythosa/gowrapper/internal/services/generator/gowrap"
)

type LsCommand struct {
	finder  gowrap.Finder
	options model.SearchOptions
}

func NewLsCommand(finder gowrap.Finder) *LsCommand {
	return &LsCommand{finder: finder}
}

func (cmd *LsCommand) descriptor() *cobra.Command {
	cc := &cobra.Command{
		Use:   "ls",
		Short: "Shows found interfaces for files by search options",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.finder.Find(cmd.options.InitialDirectory, cmd.options.ExcludedDirs, cmd.options.ExcludedFiles)
		},
	}

	cc.Flags().StringSliceVar(&cmd.options.ExcludedDirs,
		"excluded-dirs", nil, "parts of paths of excluded dirs")
	cc.Flags().StringSliceVar(&cmd.options.ExcludedFiles,
		"excluded-files", nil, "parts of paths of excluded files")
	cc.Flags().StringVarP(&cmd.options.InitialDirectory,
		"dir", "d", "", "path of initial directory")

	return cc
}
