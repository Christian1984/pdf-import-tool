package main

import (
	"os"

	"fyne.io/fyne/v2/dialog"
)

func ShowInformation(message string) {
	dialog.ShowInformation("Info", message, *ParentWindow)
}

func ShowError(message string) {
	dialog.ShowInformation("Something Went Wrong", message, *ParentWindow)
}

func ShowErrorAndExit(message string) {
	dialog.ShowConfirm("Something Went Wrong", message+"\nClick \"Yes\" to EXIT or \"No\" to CONTINUE (may result in an unstable experience)!", func(b bool) {
		if b {
			os.Exit(0)
		}
	}, *ParentWindow)
}

func ShowClearImportFolderPrompt(callback func(bool)) {
	dialog.ShowConfirm("PDF Importer", "This will delete all files and folders inside [FSKneeboard]\\charts\\!import - Do you want to proceed?", func(b bool) {
		callback(b)
	}, *ParentWindow)
}
