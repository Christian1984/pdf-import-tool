package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var in string
	var out string
	var verbose bool

	flag.StringVar(&in, "in", "", "the input file as a relative or absolute path, e.g. .\\in\\test.pdf")
	flag.StringVar(&out, "out", "", "the output file path (pattern) as a relative or absolute path, e.g. .\\out\\out--%03d.png")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose mode")
	flag.Parse()

	if verbose {
		fmt.Println("in: " + in)
		fmt.Println("out: " + out)
	}

	absIn, inErr := filepath.Abs(in)
	absOut, outErr := filepath.Abs(out)

	if strings.TrimSpace(in) == "" || strings.TrimSpace(out) == "" || inErr != nil || outErr != nil {
		fmt.Println("The provided input and/or output paths could not be parsed!")
		os.Exit(1)
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
		"-r150x150",

		"-o",
		absOut,
		absIn,
	}

	cmd := exec.Command(".\\gswin64c.exe", cmdParams...)

	fmt.Println("Starting import for file " + absIn)
	s, importErr := cmd.Output()

	result := string(s)

	if verbose {
		fmt.Println("Import output:\n" + result)
	}

	if importErr != nil {
		fmt.Println("Could not import PDF file, reason: " + importErr.Error())
		os.Exit(1)
	} else if strings.Contains(strings.ToLower(result), "error") {
		fmt.Println("Could not import PDF file, errors occured!")
		os.Exit(1)
	} else {
		fmt.Println("Import successful!")
	}

	os.Exit(0)
}
