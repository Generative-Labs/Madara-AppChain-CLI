package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func performAppChainEVMInstallation(version string) error {
	InstallAppChainEVMPkg(version)

	return nil
}

func performMadaraInstallation(version string) error {
	InstallMadaraPkg(version)

	return nil
}

func performDojoInstallation(version string) error {
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

func parsePackageInfo(arg string) (string, string) {
	parts := strings.Split(arg, "@")
	if len(parts) == 1 {
		return parts[0], "latest"
	}

	if len(parts) != 2 {
		fmt.Println("Invalid argument format")
		return "", ""
	}

	packageName := parts[0]
	version := parts[1]

	return packageName, version
}

// installCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		// version, _ := cmd.Flags().GetString("version")
		// fmt.Printf("Installing dependencies with version: %s\n", version)

		for _, arg := range args {
			_, version := parsePackageInfo(arg)

			if strings.HasPrefix(arg, "madara") {
				// Perform the installation
				if err := performMadaraInstallation(version); err != nil {
					fmt.Println("Error during madara installation:", err)
					os.Exit(1)
				}
			}

			if strings.HasPrefix(arg, "appchain-evm") {
				// Perform the installation
				if err := performAppChainEVMInstallation(version); err != nil {
					fmt.Println("Error during appchain-evm installation:", err)
					os.Exit(1)
				}
			}

			if strings.HasPrefix(arg, "dojo") {
				// Perform the installation
				if err := performDojoInstallation(version); err != nil {
					fmt.Println("Error during dojo installation:", err)
					os.Exit(1)
				}

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
			}
		}

		fmt.Println("Installation completed successfully.")
	},
}

// UpdateSpecCmd represents the update spec command
var UpdateSpecCmd = &cobra.Command{
	Use:   "update-spec",
	Short: "Update spec dependencies",
	Run: func(cmd *cobra.Command, args []string) {

		for _, arg := range args {

			if strings.HasPrefix(arg, "madara") {
				// Perform the installation
				if err := performUpdateChainSpec(); err != nil {
					fmt.Println("Error during update chain spec :", err)
					os.Exit(1)
				}
			}
		}
	},
}
