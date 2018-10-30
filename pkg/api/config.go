package api

type Config struct {
	ServerHost string `mapstructure:"server_host"`
	ServerPort int    `mapstructure:"server_port"`
}
