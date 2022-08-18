package main

import (
	"flag"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const APP_NAME = "FSKneeboard PDF Importer"

var ParentWindow *fyne.Window
var AbsInputRoot string
var AbsOutputRoot string

func main() {
	flag.StringVar(&AbsInputRoot, "in", "", "specify input root")
	flag.StringVar(&AbsOutputRoot, "out", "", "specify output root")
	flag.Parse()

	a := app.New()
	w := a.NewWindow(APP_NAME)
	importPanel := PdfImportPanel()
	max := container.NewMax(importPanel)
	w.SetContent(max)
	w.Resize(fyne.NewSize(800, 600))
	ParentWindow = &w

	if AbsInputRoot == "" || AbsOutputRoot == "" {
		ShowErrorAndExit("Please specify root directories for input and output.")
	}

	w.ShowAndRun()
}
