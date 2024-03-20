package internal

import (
	"github.com/apex/log"
	"github.com/spf13/viper"
)

var ViperReader *viper.Viper

func init() {
	// Initialize New Viper to avoid collisions
	ViperReader = viper.New()
	setupConfig()
}

func setupConfig() {
	// set config name and config type -> config.json
	ViperReader.SetConfigName("config")
	ViperReader.SetConfigType("json")
	// set multiple path for config to accomodate docker configuration and local for future configuration and testing configuration
	ViperReader.AddConfigPath("../../")
	ViperReader.AddConfigPath("../")
	ViperReader.AddConfigPath(".")
	ViperReader.AddConfigPath("../../app")
	ViperReader.AddConfigPath("../app")
	ViperReader.AddConfigPath("/app")
	// read config
	if err := ViperReader.ReadInConfig(); err != nil {
		log.Fatalf("[viper][init] cannot read config file: %v\n", err)
	}
	log.Debugf("[viper][init] config file used: %s", ViperReader.ConfigFileUsed())
	log.Debugf("[viper][init] config value: %v", ViperReader.AllSettings())
}
