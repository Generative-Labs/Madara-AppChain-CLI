package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

var cfgFile string

// Config structure to hold configuration parameters
type Config struct {
	Madara string `toml:"madara"`
	Torii  string `toml:"torii"`
}

func parseArgs(args []string, defaultConfig Config) Config {
	config := defaultConfig

	for _, arg := range args {
		parts := strings.Split(arg, "@")
		if len(parts) == 2 {
			switch parts[0] {
			case "torii":
				config.Torii = parts[1]
			case "madara":
				config.Madara = parts[1]
			}
		}
	}

	return config
}

func initConfig(args []string) {
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
		if len(args) > 0 {

		}

		defaultConfig := Config{
			Madara: "0.1.0",
			Torii:  "0.3.15",
		}
		config := parseArgs(args, defaultConfig)

		viper.Set("madara", config.Madara)
		viper.Set("torii", config.Torii)

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
