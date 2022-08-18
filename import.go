package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var importerBasePath, _ = filepath.Abs("pdf-importer")
var importerCallerPath, _ = filepath.Abs(importerBasePath + "\\" + "pdf-import-runner.exe")
var importerExePath, _ = filepath.Abs(importerBasePath + "\\" + "gswin64c.exe")
var importerDllPath, _ = filepath.Abs(importerBasePath + "\\" + "gsdll64.dll")
var importerLicensePath, _ = filepath.Abs(importerBasePath + "\\" + "THIRD-PARTY-LICENSE.md")

const importerBaseUrl = "https://github.com/Christian1984/pdf-importer/releases/download/v1.1.1/"
const importerCallerUrl = importerBaseUrl + "pdf-import-runner.exe"
const importerExeUrl = importerBaseUrl + "gswin64c.exe"
const importerDllUrl = importerBaseUrl + "gsdll64.dll"
const importerLicenseUrl = importerBaseUrl + "THIRD-PARTY-LICENSE.md"

const sourceFolder = "charts\\!import"
const outFolder = "charts\\imported"

type PdfFileInfo struct {
	SourcePath string
	TargetPath string
	FileName   string
}

func importPdfChart(sourcePath string, targetBasePath string, filename string) error {
	fmt.Println("Enter importPdfChart...")

	documentTargetPath, _ := filepath.Abs(targetBasePath + "\\" + strings.TrimSuffix(filename, ".pdf"))

	in, _ := filepath.Abs(sourcePath + "\\" + filename)
	out, _ := filepath.Abs(documentTargetPath + "\\" + strings.TrimSuffix(filename, ".pdf") + "-%03d.png")

	fmt.Println("Starting PDF Import of " + in)
	fmt.Println("Creating out path at " + documentTargetPath)
	mkDirErr := os.MkdirAll(documentTargetPath, os.ModePerm)

	if mkDirErr != nil {
		fmt.Println("Could not create target folder, reason: " + mkDirErr.Error())
		return mkDirErr
	}

	cmdParams := []string{
		"--out",
		out,
		"--in",
		in,
		"--gspath",
		".\\pdf-importer",
		"--verbose",
	}

	cmd := exec.Command(importerCallerPath, cmdParams...)
	fmt.Println("Import command is: " + cmd.String())

	s, importErr := cmd.Output()
	result := string(s)
	fmt.Println("Import output:\n" + result)

	if importErr != nil {
		fmt.Println("Could not import PDF file, reason: " + importErr.Error())
		return importErr
	} else {
		fmt.Println("Import successful!")
	}

	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) {
		fmt.Println(filePath + " not found, reason: " + err.Error())
		return false
	} else if err != nil {
		fmt.Println("Something when wrong when searching for " + filePath + ", reason: " + err.Error())
		return false
	}

	fmt.Println(filePath + " found!")
	return true
}

func downloadFile(filepath string, url string) error {
	fmt.Println("Downloading from " + url + " to " + filepath)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func createAndOpenFolder(path string) {
	absPath, _ := filepath.Abs(path)
	os.MkdirAll(absPath, os.ModePerm)

	fmt.Println("Trying to open folder [" + path + "]...")
	err := OpenExplorer(path)

	if err != nil {
		fmt.Println("Could not open folder [" + path + "]")
		ShowError("Folder could not be opened! Reason: " + err.Error())
	}
}

func CreatePdfFileList() ([]PdfFileInfo, error) {
	list := []PdfFileInfo{}

	absPath, _ := filepath.Abs(sourceFolder)
	sourceFolderErr := os.MkdirAll(absPath, os.ModePerm)

	if sourceFolderErr != nil {
		fmt.Println("Import failed! Could not find or create import source folder, reason: " + sourceFolderErr.Error())
		return list, sourceFolderErr
	}

	walkErr := filepath.Walk(sourceFolder, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Println("%v", info)
		} else {
			// check if pdf
			if !strings.HasSuffix(strings.ToLower(info.Name()), ".pdf") {
				fmt.Println("Skipping file " + info.Name() + ", reason: Extension [.pdf] missing!")
				return nil
			}

			fmt.Println("%v", info)
			sourcePath := strings.TrimSuffix(filePath, "\\"+info.Name()) // e.g. charts\!import\folder 1
			relPath := strings.TrimPrefix(sourcePath, sourceFolder)      // e.g. \folder 1
			relPath = strings.TrimPrefix(relPath, "\\")                  // e.g. folder 1
			targetPath := outFolder + "\\" + relPath                     // e.g. charts\imported\folder 1

			fmt.Println("filePath: " + filePath)
			fmt.Println("sourcePath: " + sourcePath)
			fmt.Println("sourceRoot: " + sourceFolder)
			fmt.Println("relPath: " + relPath)
			fmt.Println("targetPath: " + targetPath)

			list = append(list, PdfFileInfo{SourcePath: sourcePath, TargetPath: targetPath, FileName: info.Name()})
		}

		return nil
	})

	if walkErr != nil {
		fmt.Println("Could not import folder " + sourceFolder + ", reason: " + walkErr.Error())
		return list, walkErr
	}

	return list, nil
}

func ImportPdfFolder(updateStatusCallback func(string)) error {
	fmt.Println("Enter ImportPdfFolder...")

	pdfList, listErr := CreatePdfFileList()

	if listErr != nil {
		fmt.Println("Could not create list of PDF files, reason: " + listErr.Error())
		return listErr
	}

	if len(pdfList) == 0 {
		fmt.Println("No PDFs found in import folder!")
		return errors.New("No PDFs found in import folder!")
	}

	for _, fileInfo := range pdfList {
		if updateStatusCallback != nil {
			updateStatusCallback("Importing PDF: [" + fileInfo.FileName + "]")
		}

		importErr := importPdfChart(fileInfo.SourcePath, fileInfo.TargetPath, fileInfo.FileName)

		if importErr != nil {
			fmt.Println("Could not import file " + fileInfo.FileName + ", reason: " + importErr.Error())
			return importErr
		}
	}

	fmt.Println("Import process finished!")
	return nil
}

func HasImporter() bool {
	if !fileExists(importerCallerPath) || !fileExists(importerExePath) || !fileExists(importerDllPath) {
		fmt.Println("Local importer binaries not found or incomplete!")
		return false
	}

	fmt.Println("Local importer binaries found!")
	return true
}

func DownloadImporter() error {
	fmt.Println("Downloading Importer...")

	mkDirErr := os.MkdirAll(importerBasePath, os.ModePerm)

	if mkDirErr != nil {
		fmt.Println("Could not create importer base directory, reason: " + mkDirErr.Error())
		return mkDirErr
	}

	if err := downloadFile(importerCallerPath, importerCallerUrl); err != nil {
		fmt.Println("Could not download importer caller, reason: " + err.Error())
		return err
	}

	if err := downloadFile(importerExePath, importerExeUrl); err != nil {
		fmt.Println("Could not download importer executable, reason: " + err.Error())
		return err
	}

	if err := downloadFile(importerDllPath, importerDllUrl); err != nil {
		fmt.Println("Could not download importer dll, reason: " + err.Error())
		return err
	}

	if err := downloadFile(importerLicensePath, importerLicenseUrl); err != nil {
		fmt.Println("Could not download importer license, reason: " + err.Error())
		return err
	}

	fmt.Println("Importer successfully downloaded...")

	return nil
}

func OpenPdfSourceFolder() {
	fmt.Println("Trying to open PDF import source folder [" + sourceFolder + "]...")
	createAndOpenFolder(sourceFolder)
}

func OpenPdfOutFolder() {
	fmt.Println("Trying to open PDF import out folder [" + outFolder + "]...")
	createAndOpenFolder(outFolder)
}

func ClearPdfImportFolder() error {
	absSourcePath, _ := filepath.Abs(sourceFolder)
	err := os.RemoveAll(absSourcePath)

	if err != nil {
		fmt.Println("PDF import source folder [" + sourceFolder + "] could not be cleared, reason: " + err.Error())
		return err
	}

	os.MkdirAll(absSourcePath, os.ModePerm)

	return nil
}
