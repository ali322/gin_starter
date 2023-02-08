package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var App = new(AppConf)

type AppConf struct {
	Port           string `yaml:"port"`
	Locale         string `yaml:"locale"`
	LogDir         string `yaml:"logDir"`
	JWTSecret      string `yaml:"jwtSecret"`
	GroupAdminRole string `yaml:"groupAdminRole"`
	DefaultRole    string `yaml:"defaultRole"`
	Dsn            string `yaml:"dsn"`
}

func Read() {
	workDir, _ := os.Getwd()
	viper.SetConfigFile(filepath.Join(workDir, "config.yml"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Sub("app").Unmarshal(App); err != nil {
		log.Fatal(err)
	}
}
