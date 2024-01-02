package pkg

import (
	"fmt"

	"github.com/spf13/cobra"
)

// updateGenesisConfigCmd represents the update-genesis-config command
var GenesisConfigCmd = &cobra.Command{
	Use:   "update-genesis-config",
	Short: "Update genesis config",
	Run: func(cmd *cobra.Command, args []string) {
		customGenesisFile, _ := cmd.Flags().GetString("custom-genesis-file")
		defaultGenesisFile, _ := cmd.Flags().GetString("default-genesis-file")
		fmt.Printf("Updating genesis config with custom file: %s, default file: %s\n", customGenesisFile, defaultGenesisFile)
	},
}
