package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Panic(err)
	}
	// print env from pointer
	fmt.Println(cfg.EnvNumber, ";", cfg.EnvString)
}

// Config,
type Config struct {
	EnvString string `yaml:"env_string" mapstructure:"env_string" json:"env_string"`
	EnvNumber int    `yaml:"env_number" mapstructure:"env_number" json:"env_number"`
}

func getConfigName() string {
	mode := "dev"
	switch os.Getenv("MODE") {
	case "uat":
		mode = "uat"
	case "prod":
		mode = "prod"
	}
	return mode
}

func LoadConfig() (*Config, error) {
	mode := getConfigName()
	fmt.Println("mode:", mode)
	cfg := &Config{}
	path := "/Users/ducnp/src/github/open-dev/usecase/environment"
	vn := viper.New()
	vn.AddConfigPath(path)
	vn.SetConfigName(mode)
	vn.SetConfigType("yaml")
	vn.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	vn.AutomaticEnv()

	err := vn.ReadInConfig()
	if err != nil {
		log.Panic(err)
		return cfg, err
	}

	for _, key := range vn.AllKeys() {
		str := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		log.Println(key, str, vn.Get(key))
		vn.BindEnv(key, str)
	}

	err = vn.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
