package budgets

import (
	"fmt"

	"github.com/spf13/viper"
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
	viper.SetConfigFile("/home/ec2-user/.private/qbase.yaml")

	Logger.Debugw("Config file Used", "file", viper.ConfigFileUsed())

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		Logger.Fatalw("Error reading config file", "err", err)
	}
	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	v := viper.Sub(stage)
	var C qbaseConf
	err := v.Unmarshal(&C)
	if err != nil {
		Logger.Fatalw("unable to decode configuration for quickbase into a struct", "err", err)
	}
	return C
}
