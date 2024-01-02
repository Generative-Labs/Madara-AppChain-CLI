package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{Use: "madcli"}

// initCmd represents the init command
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize madcli",
    Run: func(cmd *cobra.Command, args []string) {
    	fmt.Println("args:", args)
        fmt.Println("Initializing madcli...")
    },
}

// installCmd represents the install command
var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install dependencies",
    Run: func(cmd *cobra.Command, args []string) {
        version, _ := cmd.Flags().GetString("version")
        fmt.Printf("Installing dependencies with version: %s\n", version)
    },
}

// startCmd represents the start command
var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start madara and torii services etc...",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Starting madara and torii services...")
    },
}

// getNodeStateCmd represents the get-node-state command
var getNodeStateCmd = &cobra.Command{
    Use:   "get-node-state",
    Short: "Query node state",
    Run: func(cmd *cobra.Command, args []string) {
        rpcURL, _ := cmd.Flags().GetString("rpc-url")
        fmt.Printf("Querying node state with RPC URL: %s\n", rpcURL)
    },
}

// updateGenesisConfigCmd represents the update-genesis-config command
var updateGenesisConfigCmd = &cobra.Command{
    Use:   "update-genesis-config",
    Short: "Update genesis config",
    Run: func(cmd *cobra.Command, args []string) {
        customGenesisFile, _ := cmd.Flags().GetString("custom-genesis-file")
        defaultGenesisFile, _ := cmd.Flags().GetString("default-genesis-file")
        fmt.Printf("Updating genesis config with custom file: %s, default file: %s\n", customGenesisFile, defaultGenesisFile)
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
    rootCmd.AddCommand(installCmd)
    rootCmd.AddCommand(startCmd)
    rootCmd.AddCommand(getNodeStateCmd)
    rootCmd.AddCommand(updateGenesisConfigCmd)

    // Flags for install command
    installCmd.Flags().String("version", "", "Specify version")

    // Flags for get-node-state command
    getNodeStateCmd.Flags().String("rpc-url", "", "Specify RPC URL")

    // Flags for update-genesis-config command
    updateGenesisConfigCmd.Flags().String("custom-genesis-file", "", "Specify custom genesis file")
    updateGenesisConfigCmd.Flags().String("default-genesis-file", "", "Specify default genesis file")
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
