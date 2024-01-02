package pkg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func performInstallation(version string) error {
	// Add your installation logic here.
	// For example, you can use curl and bash to install:
	cmd := exec.Command("bash", "-c", fmt.Sprintf("curl -L https://install.dojoengine.org | bash -s %s", version))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDojoup() error {
	// Run the dojoup command
	cmd := exec.Command("dojoup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// installCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetString("version")
		fmt.Printf("Installing dependencies with version: %s\n", version)

		// Perform the installation
		if err := performInstallation(version); err != nil {
			fmt.Println("Error during installation:", err)
			os.Exit(1)
		}

		fmt.Println("Installation completed successfully.")

		// Source the bash script
		// if err := sourceBashScript(); err != nil {
		// 	fmt.Println("Error sourcing bash script:", err)
		// 	os.Exit(1)
		// }

		// Run the dojoup command
		if err := runDojoup(); err != nil {
			fmt.Println("Error running dojoup:", err)
			os.Exit(1)
		}
	},
}
