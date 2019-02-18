package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/supergiant/analyze/cmd/analyzed/app"
)

func main() {

	command := &cobra.Command{
		Use:          "analyzed",
		Short:        "analyze service checks K8s cluster by means of installed plugins and gives recommendations",
		RunE:         app.RunCommand,
		SilenceUsage: true,
	}

	command.PersistentFlags().StringArrayP(
		"config",
		"c",
		[]string{"./config.yaml", "/etc/analyzed/config.yaml", "$HOME/.analyzed/config.yaml"},
		"config file path")

	if err := command.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}
