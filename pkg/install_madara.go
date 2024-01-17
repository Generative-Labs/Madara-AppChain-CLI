package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func GetPkgName() string {
	osName := runtime.GOOS
	machineName := runtime.GOARCH

	// Map OS names to the desired format
	osMap := map[string]string{
		"darwin": "apple-darwin",
		"linux":  "unknown-linux-gnu",
	}

	// Map machine names to the desired format
	machineMap := map[string]string{
		"x86_64": "x86_64",
		"arm64":  "aarch64",
		"amd64":  "x86_64",
	}

	// Validate OS
	os, osExists := osMap[osName]
	if !osExists {
		fmt.Println("Unknown OS")
		return ""
	}

	// fmt.Println(">>machineName: ", machineName)

	// Validate machine
	machine, machineExists := machineMap[machineName]
	if !machineExists {
		fmt.Println("Unknown machine")
		return ""
	}

	// Construct the package name
	PkgName := fmt.Sprintf("%s-%s-madara", machine, os)
	return PkgName
}

func GetMadupDir() string {
	baseDir := os.Getenv("XDG_CONFIG_HOME")
	if baseDir == "" {
		baseDir = os.Getenv("HOME")
	}

	// Set MADUP_DIR with default value if not set

	madupDir := filepath.Join(baseDir, ".madup")

	return madupDir
}

func InstallMadaraPkg(version string) {
	PkgName := GetPkgName()
	if PkgName == "" {
		return
	}

	MADUP_DIR := GetMadupDir()

	// Build the download URL
	url := fmt.Sprintf("https://github.com/Generative-Labs/madara-and-dojo/releases/latest/download/%s.zip", PkgName)

	fmt.Println("Download madara file...")
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

	// Unzip the file
	err = unzip(filepath.Join(MADUP_DIR, PkgName+".zip"), MADUP_DIR)
	if err != nil {
		fmt.Println("Failed to unzip file:", err)
		return
	}

	// Move the madara directory to bin
	srcDir := filepath.Join(MADUP_DIR, PkgName, "madara")
	destDir := filepath.Join(MADUP_DIR, "bin", "madara")

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
	err = os.Chmod(filepath.Join(MADUP_DIR, "bin", "madara"), 0755)
	if err != nil {
		fmt.Println("Failed to set permissions:", err)
		return
	}

	fmt.Println("Installation completed")
}

// Download file function
func downloadFile(url, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}

// Unzip a file
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
