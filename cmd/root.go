package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/supergiant/robot/robot"
)

const (
	ConfigFileName = "robot"
)

var serverCmd = &cobra.Command{
	Use:     "robot",
	Short:   "Use qbox robot for ",
	Aliases: []string{"server", "srv", "s"},
}

func initCfg() {
	viper.SetConfigFile(ConfigFileName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/robot")
	viper.AddConfigPath("$HOME")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("using config file %s", viper.ConfigFileUsed())
	cfg := new(robot.Config)
	if err := viper.Unmarshal(cfg); err != nil {
		logrus.Fatal("failed to read configuration")
	}
}

func Execute() {
	cobra.OnInitialize(initCfg)
	if err := serverCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
