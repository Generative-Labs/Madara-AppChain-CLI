package pkg

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize madcli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args:", args)
		initConfig(args)
		fmt.Println("Initializing madcli...")
	},
}
