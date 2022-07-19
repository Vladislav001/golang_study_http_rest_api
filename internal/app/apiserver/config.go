package apiserver

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLever string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		LogLever: "debug",
	}
}
