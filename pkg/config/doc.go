// Package config is a wrapper on viper to simplify config struct populating
// current implementation works with single viper instance
// casual flow of using its public API:
// ReadFromFiles() -> MergeEnv()
// if ReadFromFiles() invocation is skipped than MergeEnv() need to be invoked
// to get all config values from env variables
package config
