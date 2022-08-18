package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/Christian1984/pdf-importer/res"
)

const APP_NAME = "FSKneeboard PDF Import Tool"

var BuildVersion string

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

	fmt.Println("Loading icon...", false)
	iconAsset, err := res.Asset("icon.png")
	if err == nil {
		iconResource := fyne.NewStaticResource("icon.png", iconAsset)
		fmt.Println("Icon loaded", false)
		a.SetIcon(iconResource)
	} else {
		fmt.Println("Icon could not be loaded!", false)
	}

	w := a.NewWindow(APP_NAME + " " + BuildVersion)
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
