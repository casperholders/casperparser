/*
	Used to provide useful helpers
*/
package utils

import (
	"github.com/spf13/viper"
	"log"
)

// InitViper Helper Init Viper for the unit tests
func InitViper() error {
	viper.SetConfigName(".env")
	viper.AddConfigPath("../..")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".casperParser")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	return nil
}
