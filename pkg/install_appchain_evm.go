package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

func InstallAppChainEVMPkg(version string) {
	PkgName := GetPkgName("appchain-evm")
	if PkgName == "" {
		return
	}

	MADUP_DIR := GetMadupDir()

	// Build the download URL
	url := fmt.Sprintf("https://github.com/Generative-Labs/AppChainEVMPackage/releases/latest/download/%s.zip", PkgName)

	fmt.Println("Download appchain-evm file...")
	// Download the file
	err := downloadFile(url, PkgName+".zip")
	if err != nil {
		fmt.Println("Failed to download file:", err)
		return
	}

	// Create the target directory
	err = os.MkdirAll(MADUP_DIR, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		return
	}

	// Move the file to the target directory
	err = os.Rename(PkgName+".zip", filepath.Join(MADUP_DIR, PkgName+".zip"))
	if err != nil {
		fmt.Println("Failed to move file:", err)
		return
	}

	filePath := filepath.Join(MADUP_DIR, PkgName+".zip")
	fmt.Println(">>> filePath: ", filePath)

	err = os.Chown(filePath, os.Getuid(), os.Getgid())
	if err != nil {
		fmt.Println("Failed to set zip Chown:", err)
		return
	}

	// Set zip permissions
	err = os.Chmod(filePath, 0655)
	if err != nil {
		fmt.Println("Failed to set zip permissions:", err)
		return
	}

	// Unzip the file
	err = unzipv2(filePath, MADUP_DIR)
	if err != nil {
		fmt.Println("Failed to unzip file:", err)
		return
	}

	// Move the appchain evm directory to bin
	srcDir := filepath.Join(MADUP_DIR, PkgName, "AppChainEVM")
	destDir := filepath.Join(MADUP_DIR, "bin", "AppChainEVM")

	if _, err := os.Stat(destDir); err == nil {
		err := os.RemoveAll(destDir)
		if err != nil {
			fmt.Println("Failed to RemoveAll file:", err)

			return
		}
	}

	// Rename the directory
	err = os.Rename(srcDir, destDir)
	if err != nil {
		fmt.Println("Failed to move directory:", err)
		return
	}

	// Set executable permissions
	err = os.Chmod(filepath.Join(MADUP_DIR, "bin", "AppChainEVM"), 0755)
	if err != nil {
		fmt.Println("Failed to set permissions:", err)
		return
	}

	fmt.Println("Installation completed")
}
