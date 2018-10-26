package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/supergiant/robot/cmd/robot/run"
)

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "robot",
		Short: "robot kelly service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		SilenceUsage: true,
	}

	root.PersistentFlags().StringArrayP(
		"config",
		"c",
		[]string{"./robot.yaml", "/etc/robot/robot.yaml", "$HOME/robot.yaml"},
		"config file path")

	root.AddCommand(run.NewCommand())

	return root
}

func main() {
	root := newRootCommand()
	if err := root.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}
