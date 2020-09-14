package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sb10/interpolate/interpolate"
)

const expectedArgs = 2

func main() {
	exit := 1

	defer func() {
		os.Exit(exit)
	}()

	if len(os.Args) != expectedArgs {
		fmt.Printf("\nusage: interpolation <pathto.csv>\n\n")

		return
	}

	csvPath := os.Args[1]

	f, err := os.Open(csvPath)
	if err != nil {
		fmt.Printf("could not open csv file '%s': %s\n", csvPath, err)

		return
	}
	defer f.Close()

	exit = parseCSV(f)
	exit = 0
}

func parseCSV(f io.ReadSeeker) int {
	// TODO: determine the number of decimal places from the input data
	c, err := interpolate.NewCSVInterpolator(f, 6)
	if err != nil {
		fmt.Printf("could not create an interpolator: %s\n", err)

		return 1
	}

	for {
		row, err := c.NextRow()
		if err != nil {
			break
		}

		fmt.Printf("%s\n", strings.Join(row, ","))
	}

	return 0
}
