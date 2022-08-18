package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const APP_NAME = "FSKneeboard PDF Import Tool"

var ParentWindow *fyne.Window

var relLibPath string
var relInputRoot string
var relOutputRoot string

var AbsLibPath string
var AbsInputRoot string
var AbsOutputRoot string

var ImporterExePath string
var ImporterDllPath string

func main() {
	flag.StringVar(&relLibPath, "lib", ".\\lib", "specify lib path")
	flag.StringVar(&relInputRoot, "in", ".\\in", "specify input root")
	flag.StringVar(&relOutputRoot, "out", ".\\out", "specify output root")
	flag.Parse()

	AbsLibPath, _ = filepath.Abs(relLibPath)
	AbsInputRoot, _ = filepath.Abs(relInputRoot)
	AbsOutputRoot, _ = filepath.Abs(relOutputRoot)

	ImporterExePath, _ = filepath.Abs(AbsLibPath + "\\" + "gswin64c.exe")
	ImporterDllPath, _ = filepath.Abs(AbsLibPath + "\\" + "gsdll64.dll")

	fmt.Println("AbsLibPath: " + AbsLibPath)
	fmt.Println("AbsInputRoot: " + AbsInputRoot)
	fmt.Println("AbsOutputRoot: " + AbsOutputRoot)

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
