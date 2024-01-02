package pkg

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getNodeStateCmd represents the get-node-state command
var NodeStateCmd = &cobra.Command{
	Use:   "get-node-state",
	Short: "Query node state",
	Run: func(cmd *cobra.Command, args []string) {
		rpcURL, _ := cmd.Flags().GetString("rpc-url")
		fmt.Printf("Querying node state with RPC URL: %s\n", rpcURL)
	},
}
