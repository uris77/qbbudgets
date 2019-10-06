package budgets

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type qbaseConf struct {
	QbUsername    string
	QbPassword    string
	QbAppToken    string
	TeamworkToken string
	DbUser        string
	DbPassword    string
	DbName        string
	DbHost        string
}

func ReadConf(stage string) qbaseConf {
	viper.AddConfigPath("$HOME/.private")
	viper.AddConfigPath("./")
	viper.SetConfigName("qbase")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	v := viper.Sub(stage)
	var C qbaseConf
	err := v.Unmarshal(&C)
	if err != nil {
		log.Fatalf("unable to decode configuration for quickbase into a struct, %v", err)
	}
	return C
}
