package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" default:"development"`
	PostgresConfig `yaml:"database"`
}

// postgres config
type PostgresConfig struct {
	Host     string `yaml:"host" default:"localhost"`
	Port     int    `yaml:"port" default:"5432"`
	Username string `yaml:"username" default:"user"`
	Password string `yaml:"password" default:"password"`
	Name     string `yaml:"name" default:"app_db"`
}

func (cfg PostgresConfig) DSN() string {
	return "host=" + cfg.Host +
		" port=" + string(rune(cfg.Port)) +
		" user=" + cfg.Username +
		" password=" + cfg.Password +
		" dbname=" + cfg.Name +
		" sslmode=disable"
}

// -------------Get Config Path from Flag or Env --------------
var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "Path to the config file")
}

func fetchConfigPath() string {
	var res string

	if !flag.Parsed() {
		flag.Parse()
	}

	res = configPath

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		panic("config path is not provided")
	}

	return res
}
func LoadConfig() Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	return LoadConfigFromPath(path)
}

func LoadConfigFromPath(path string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
