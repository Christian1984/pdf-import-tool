package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PdfFileInfo struct {
	SourcePath string
	TargetPath string
	FileName   string
}

func importPdfChart(sourcePath string, targetBasePath string, filename string) error {
	fmt.Println("Enter importPdfChart...")

	documentTargetPath, _ := filepath.Abs(targetBasePath + "\\" + strings.TrimSuffix(filename, ".pdf"))

	absIn, _ := filepath.Abs(sourcePath + "\\" + filename)
	absOut, _ := filepath.Abs(documentTargetPath + "\\" + strings.TrimSuffix(filename, ".pdf") + "-%03d.png")

	fmt.Println("Starting PDF Import of " + absIn)
	fmt.Println("Creating out path at " + documentTargetPath)
	mkDirErr := os.MkdirAll(documentTargetPath, os.ModePerm)

	if mkDirErr != nil {
		fmt.Println("Could not create target folder, reason: " + mkDirErr.Error())
		return mkDirErr
	}

	cmdParams := []string{
		//"-q",
		//"-dQUIET",
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-dNOPROMPT",
		"-dMaxBitmap=500000000",
		"-dAlignToPixels=0",
		"-dGridFitTT=2",
		"-sDEVICE=png16m",
		"-dTextAlphaBits=4",
		"-dGraphicsAlphaBits=4",
		"-r600",

		"-o",
		absOut,
		absIn,
	}

	cmd := exec.Command(ImporterExePath, cmdParams...)
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
	fmt.Println("Looking for file " + filePath)
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

func createAndOpenFolder(path string) {
	absPath, _ := filepath.Abs(path)
	os.MkdirAll(absPath, os.ModePerm)

	fmt.Println("Trying to open folder [" + absPath + "]...")
	cmd := exec.Command("explorer", absPath)
	cmd.Run()

	// if err != nil {
	// fmt.Println("Could not open folder [" + path + "]")
	// ShowError("Folder could not be opened! Reason: " + err.Error())
	// }
}

func CreatePdfFileList() ([]PdfFileInfo, error) {
	list := []PdfFileInfo{}

	absPath := AbsInputRoot
	sourceFolderErr := os.MkdirAll(absPath, os.ModePerm)

	if sourceFolderErr != nil {
		fmt.Println("Import failed! Could not find or create import source folder, reason: " + sourceFolderErr.Error())
		return list, sourceFolderErr
	}

	walkErr := filepath.Walk(AbsInputRoot, func(filePath string, info os.FileInfo, err error) error {
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
			relPath := strings.TrimPrefix(sourcePath, AbsInputRoot)      // e.g. \folder 1
			relPath = strings.TrimPrefix(relPath, "\\")                  // e.g. folder 1
			targetPath := AbsOutputRoot + "\\" + relPath                 // e.g. charts\imported\folder 1

			fmt.Println("filePath: " + filePath)
			fmt.Println("sourcePath: " + sourcePath)
			fmt.Println("sourceRoot: " + AbsInputRoot)
			fmt.Println("relPath: " + relPath)
			fmt.Println("targetPath: " + targetPath)

			list = append(list, PdfFileInfo{SourcePath: sourcePath, TargetPath: targetPath, FileName: info.Name()})
		}

		return nil
	})

	if walkErr != nil {
		fmt.Println("Could not import folder " + AbsInputRoot + ", reason: " + walkErr.Error())
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
	if !fileExists(ImporterExePath) || !fileExists(ImporterDllPath) {
		fmt.Println("Local importer binaries not found or incomplete!")
		return false
	}

	fmt.Println("Local importer binaries found!")
	return true
}

func OpenPdfSourceFolder() {
	fmt.Println("Trying to open PDF import source folder [" + AbsInputRoot + "]...")
	createAndOpenFolder(AbsInputRoot)
}

func OpenPdfOutFolder() {
	fmt.Println("Trying to open PDF import out folder [" + AbsOutputRoot + "]...")
	createAndOpenFolder(AbsOutputRoot)
}

func ClearPdfImportFolder() error {
	absSourcePath := AbsInputRoot
	err := os.RemoveAll(absSourcePath)

	if err != nil {
		fmt.Println("PDF import source folder [" + AbsInputRoot + "] could not be cleared, reason: " + err.Error())
		return err
	}

	os.MkdirAll(absSourcePath, os.ModePerm)

	return nil
}
