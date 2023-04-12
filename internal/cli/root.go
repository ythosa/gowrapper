package cli

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type descriptable interface {
	descriptor() *cobra.Command
}

type RootCommand struct {
	*cobra.Command
}

func NewRootCommand() RootCommand {
	cmd := RootCommand{Command: &cobra.Command{}}
	cmd.init()

	return cmd
}

func (cmd RootCommand) RegisterCommands(cmds ...descriptable) RootCommand {
	cmd.AddCommand(lo.Map(cmds, func(c descriptable, _ int) *cobra.Command { return c.descriptor() })...)

	return cmd
}

func (cmd RootCommand) Execute() {
	cmd.init()

	if err := cmd.Command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (cmd RootCommand) init() {
	cmd.Use = "gowrapper"
}
