package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bensaufley/toSarif/formats/phpCsFixer"
	"github.com/bensaufley/toSarif/formats/pyright"
	"github.com/bensaufley/toSarif/formats/sarif"
	"github.com/bensaufley/toSarif/util"
	"github.com/pborman/getopt/v2"
)

type Format = string

var (
	Pyright    Format = "pyright"
	PhpCsFixer Format = "php-cs-fixer"
)

var formats = map[string]Format{
	"pyright":      Pyright,
	"php-cs-fixer": PhpCsFixer,
}

type Convertible interface {
	ToSarif() (*sarif.Sarif22SchemaJson, error)
}

func main() {
	helpFlag := getopt.BoolLong("help", 'h', "Display this help message")
	rawFormat := getopt.StringLong("format", 'f', "", fmt.Sprintf("Source format. Options: %v", strings.Join(util.Keys(formats), ", ")))
	inputPath := getopt.StringLong("input", 'i', "", "Input file")
	outputPath := getopt.StringLong("output", 'o', "", "Output file")

	getopt.Parse()

	if *helpFlag {
		getopt.Usage()
		os.Exit(0)
	}

	var format Format
	var ok bool
	if format, ok = formats[strings.ToLower(*rawFormat)]; !ok {
		os.Stderr.WriteString(fmt.Sprintf("Invalid format: \"%s\"\n\n", *rawFormat))
		getopt.Usage()
		os.Exit(1)
	}

	if inputPath == nil || *inputPath == "" {
		os.Stderr.WriteString("Input file is required\n\n")
		getopt.Usage()
		os.Exit(1)
	}

	input, err := os.ReadFile(*inputPath)
	if err != nil {
		handleError("Error opening input file", err)
	}

	var original Convertible

	switch format {
	case Pyright:
		pyright := pyright.Schema{}
		if err := json.Unmarshal(input, &pyright); err != nil {
			handleError("Error parsing input file", err)
		}
		original = &pyright
	case PhpCsFixer:
		phpCsFixer := phpCsFixer.PhpCsFixerSchemaJson{}
		if err := json.Unmarshal(input, &phpCsFixer); err != nil {
			handleError("Error parsing input file", err)
		}
		original = &phpCsFixer
	default:
		os.Stderr.WriteString(fmt.Sprintf("Unsupported format: %s\n\n", format))
		getopt.Usage()
		os.Exit(1)
	}

	sarif, err := original.ToSarif()
	if err != nil {
		handleError("Error converting to SARIF", err)
	}
	marshaled, err := json.Marshal(sarif)
	if err != nil {
		handleError("Error marshaling SARIF", err)
	}

	if outputPath == nil || *outputPath == "" {
		os.Stdout.WriteString(string(marshaled))
	} else {
		if err := os.WriteFile(*outputPath, marshaled, 0644); err != nil {
			handleError("Error writing output file", err)
		}
	}
}

func handleError(message string, err error) {
	if err == nil {
		return
	}

	os.Stderr.WriteString(fmt.Sprintf("%s: %s\n", message, err))
	os.Exit(1)
}
