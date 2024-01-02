package main

import (
	"fmt"
	"madcli/pkg"

	"github.com/spf13/cobra"

	"os"
)

var rootCmd = &cobra.Command{
	Use:   "madcli",
	Short: "Command-line tool for Madara & Dojo CLI",
}

var cfgFile string

func init() {
	rootCmd.AddCommand(pkg.InitCmd)
	rootCmd.AddCommand(pkg.InstallCmd)
	rootCmd.AddCommand(pkg.StartCmd)
	rootCmd.AddCommand(pkg.NodeStateCmd)
	rootCmd.AddCommand(pkg.GenesisConfigCmd)

	// // Flags for install command
	// installCmd.Flags().String("version", "", "Specify version")

	// Flags for get-node-state command
	pkg.NodeStateCmd.Flags().String("rpc-url", "", "Specify RPC URL")

	// Flags for update-genesis-config command
	pkg.GenesisConfigCmd.Flags().String("custom-genesis-file", "", "Specify custom genesis file")
	pkg.GenesisConfigCmd.Flags().String("default-genesis-file", "", "Specify default genesis file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
