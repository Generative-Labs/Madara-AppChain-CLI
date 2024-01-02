package pkg

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

var cfgFile string

// Config structure to hold configuration parameters
type Config struct {
	Madara string `toml:"madara"`
	Torii  string `toml:"torii"`
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".madcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".madcli")
	}

	viper.AutomaticEnv() // Read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// If config file is not found, create it with default values.
		defaultConfig := Config{
			Madara: "0.1.0",
			Torii:  "0.3.15",
		}

		viper.Set("madara", defaultConfig.Madara)
		viper.Set("torii", defaultConfig.Torii)

		// Write the config file
		if err := writeConfigFile(); err != nil {
			fmt.Println("Error writing config file:", err)
		} else {
			fmt.Println("Created default config file:", viper.ConfigFileUsed())
		}
	}
}

func writeConfigFile() error {
	config := viper.AllSettings()
	tomlData, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	// Create the config file.
	configFilePath := "madcli.toml"
	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// Write the TOML data to the file.
	_, err = configFile.Write(tomlData)
	if err != nil {
		return err
	}

	return nil
}
