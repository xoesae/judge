package filesystem

import (
	"fmt"
	"github.com/xoesae/judge/runner/internal/config"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"syscall"
)

func MakeUploadDir() {
	if err := os.MkdirAll(config.GetConfig().UploadDir, 0755); err != nil {
		panic(err)
	}
}

func SaveScriptFile(file multipart.File, filename string, destinationPath string) (string, error) {
	if err := os.MkdirAll(destinationPath, 0755); err != nil {
		return "", err
	}

	fmt.Println("Created directory: " + destinationPath)

	fullPath := filepath.Join(destinationPath, filename)

	fmt.Println("Writing file: " + fullPath)

	// create the file to the given path
	savedFile, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	fmt.Println("File saved")

	// copy the contents to the file
	_, err = io.Copy(savedFile, file)
	savedFile.Close()
	if err != nil {
		return "", err
	}

	fmt.Println("File copied")

	return fullPath, nil

	//outputFile, err := os.Create(tmpPath)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = io.Copy(outputFile, file)
	//outputFile.Close()
	//if err != nil {
	//	return err
	//}
	//defer os.Remove(tmpPath)
	//
	//fmt.Println(tmpPath, destinationPath)
	//
	//// copy file to rootfs
	//srcBytes, _ := os.ReadFile(tmpPath)
	//os.WriteFile(destinationPath, srcBytes, 0755)
	//
	//return nil
}

func Mount() {
	procPath := filepath.Join(config.GetConfig().RootFs, "proc")
	err := os.MkdirAll(procPath, 0555)
	if err != nil {
		panic(err)
	}

	err = syscall.Mount("proc", procPath, "proc", 0, "")
	if err != nil {
		panic(err)
	}
}
