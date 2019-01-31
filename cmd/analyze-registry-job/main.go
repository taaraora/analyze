package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	command := &cobra.Command{
		Use:          "analyze-registry-job",
		Short:        "analyze-registry-job is job which registers or removes plugin from analyze registry",
		RunE:         runCommand,
		SilenceUsage: true,
	}

	command.PersistentFlags().BoolP(
		"remove",
		"r",
		false,
		"if true job will try to remove plugin from analyze registry")

	if err := command.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}

func runCommand(cmd *cobra.Command, _ []string) error {

	remove, err := cmd.Flags().GetBool("remove")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag remove")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	logger.Debugf("remove: %v", remove)


	return nil
}
