package config

import (
	"path/filepath"
	"strings"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ReadFromFiles reads config from file trying find config file in paths listed in `configPaths string`
func ReadFromFiles(config interface{}, configPaths []string) error {
	setDefaultValuesForViperProperties(config)

	for _, configPath := range configPaths {
		configPath, configName, ext := splitPath(configPath)
		viper.AddConfigPath(configPath)
		viper.SetConfigName(configName)
		viper.SetConfigType(ext)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return errors.Errorf("unable to read config file: %+v", err)
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return errors.Errorf("unable to decode into config struct, %v", err)
	}

	return nil
}

func setDefaultValuesForViperProperties(config interface{}) {
	m := structs.Map(config)

	for k, v := range m {
		viper.SetDefault(k, v)
	}
}

// MergeEnv Merges config properties from env variables. for simple types it replaces values.
// For slices if initial value was []string{"pod0", "pod1", "pod2"} and env variable contains only "pod33", merge result will be equal to
// []string{"pod33", "pod1", "pod2"}
func MergeEnv(envPrefix string, config interface{}) error {
	setDefaultValuesForViperProperties(config)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}

	err := viper.Unmarshal(config)
	if err != nil {
		return errors.Errorf("unable to decode into config struct, %v", err)
	}

	return nil
}

func UsedFileName() string {
	return viper.ConfigFileUsed()
}

func splitPath(fullPath string) (dir string, fileBaseName string, fileType string) {
	fullPath = filepath.Clean(fullPath)

	dir, file := filepath.Split(fullPath)
	dir = filepath.Clean(dir)

	fileExt := filepath.Ext(file)
	fileBaseName = strings.TrimSuffix(file, fileExt)
	fileType = strings.TrimPrefix(fileExt, ".")

	return
}
