package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func performUpdateChainSpec(source_path string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Can not get use home dir:", err)
		return err
	}

	source_folder := source_path + "chains/dev"

	source := source_folder + "/genesis-assets/genesis.json"
	destination := source_path + "chains/dev/genesis-assets/old_genesis.json"

	srcFile, err := os.Open(source)
	if err != nil {
		fmt.Println("Can not open source file:", err)
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		fmt.Println("Can not create destination file:", err)
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		fmt.Println("Copy file error:", err)
		return err
	}

	// Move appchain_genesis.json to genesis.json
	source_appchain_json := filepath.Join(homeDir, ".madup/appchain_genesis.json")

	srcJsonFile, err := os.Open(source_appchain_json)
	if err != nil {
		fmt.Println("Can not open source appchain_json file:", err)
		return err
	}
	defer srcFile.Close()

	dst_json := source_folder + "/genesis-assets/genesis.json"
	destJsonFile, err := os.Create(dst_json)
	if err != nil {
		fmt.Println("Can not create destination json file:", err)
		return err
	}
	defer destJsonFile.Close()

	_, err = io.Copy(srcJsonFile, destJsonFile)
	if err != nil {
		fmt.Println("Move new file error:", err)
		return err
	}

	return nil
}
