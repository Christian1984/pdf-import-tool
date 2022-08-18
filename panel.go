package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var importRunningBinding = binding.NewBool()
var statusBinding = binding.NewString()
var fileListBinding = binding.NewStringList()

func updateStatus(status string) {
	statusBinding.Set(status)
}

func runImport() {
	fmt.Println("Starting PDF import...")
	updateStatus("Preparing PDF batch import...")

	err := ImportPdfFolder(updateStatus)

	importRunningBinding.Set(false)

	if err != nil {
		fmt.Println("Something went wrong, reason: " + err.Error())
		updateStatus("PDF batch import failed!")

		ShowError("PDF Import failed! Please refer to the Console Panel and/or logs for details!")
	} else {
		fmt.Println("Import finished!")
		updateStatus("PDF batch import finished!")

		ShowInformation("The PDF Import finished!")
	}
}

func clearImportFolderPromptCallback(proceed bool) {
	if proceed {
		updateStatus("Clearing PDF import folder...")
		err := ClearPdfImportFolder()

		if err != nil {
			fmt.Println("Could not clear PDF import folder, reason: " + err.Error())
			updateStatus("Clearing PDF import folder failed!")

			ShowError("PDF import folder could not be cleared.")
		}

		go refreshImportDir()

		updateStatus("PDF import folder cleared!")
	}
}

func refreshImportDir() {
	fileListBinding.Set([]string{})
	list, err := CreatePdfFileList()

	if err != nil {
		fmt.Println("Could not refresh PDF import folder, reason: " + err.Error())
		updateStatus("Refreshing PDF import folder failed!")

		ShowError("PDF import folder could not be refreshed.")
	}

	sList := []string{}
	for _, info := range list {
		sList = append(sList, info.FileName)
	}

	fileListBinding.Set(sList)
}

func PdfImportPanel() *fyne.Container {
	fmt.Println("Initializing PDF Import Panel...")

	// top
	refreshFileListBtn := widget.NewButtonWithIcon("Refresh File List", theme.ViewRefreshIcon(), func() {
		go func() {
			updateStatus("Refreshing file list...")
			refreshImportDir()
			updateStatus("File list refreshed!")
		}()
	})

	clearImportDirBtn := widget.NewButtonWithIcon("Clear Import Directory", theme.ContentClearIcon(), func() {
		ShowClearImportFolderPrompt(clearImportFolderPromptCallback)
	})

	openImportDirBtn := widget.NewButtonWithIcon("Open Import Directory", theme.FolderOpenIcon(), func() {
		OpenPdfSourceFolder()
	})

	top := container.NewHBox(openImportDirBtn, refreshFileListBtn, clearImportDirBtn)

	// bottom
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Stop()

	statusLabel := widget.NewLabelWithData(statusBinding)
	statusLabel.Alignment = fyne.TextAlignCenter
	statusBinding.Set("Idle...")

	startImportBtn := widget.NewButtonWithIcon("Start Import", theme.MediaPlayIcon(), func() {
		go func() {
			importRunningBinding.Set(true)

			refreshImportDir()

			if HasImporter() {
				runImport()
			} else {
				fmt.Println("Importer module not found!")
				ShowError("GhostScript exectuable not found. Please reinstall and try again!")
				importRunningBinding.Set(false)
			}
		}()
	})

	importRunningBinding.AddListener(binding.NewDataListener(func() {
		importRunning, _ := importRunningBinding.Get()

		if importRunning {
			startImportBtn.Disable()
			progressBar.Start()
		} else {
			startImportBtn.Enable()
			progressBar.Stop()
		}
	}))

	openOutputDirBtn := widget.NewButtonWithIcon("Open Output Directory", theme.FolderOpenIcon(), func() {
		OpenPdfOutFolder()
	})

	bottomButtons := container.NewHBox(startImportBtn, openOutputDirBtn)
	bottom := container.NewVBox(progressBar, statusLabel, bottomButtons)

	// border layout
	border := layout.NewBorderLayout(top, bottom, nil, nil)
	var resContainer *fyne.Container

	// middle
	fileList := widget.NewListWithData(
		fileListBinding,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	fileList.OnSelected = func(id widget.ListItemID) {
		fileList.UnselectAll()
	}

	resContainer = container.New(border, top, bottom, fileList)

	go refreshImportDir()

	fmt.Println("PDF Import Panel initialized")

	return resContainer
}
