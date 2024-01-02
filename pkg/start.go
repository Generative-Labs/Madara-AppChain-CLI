package pkg

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start madara and torii services etc...",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting madara and torii services...")
	},
}
