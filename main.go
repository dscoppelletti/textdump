// Copyright (C) 2022 Dario Scoppelletti, <http://www.scoppelletti.it/>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputFile string
var outputFile string
var overwrite bool

func init() {
	flag.StringVar(&inputFile, "input", "", "Input file (default stdin).")
	flag.StringVar(&outputFile, "output", "", "Output file (default stdout).")
	flag.BoolVar(&overwrite, "overwrite", false,
		"Whether output file can be overwritten (default false).")
}

func main() {
	flag.Parse()
	err := checkFlags()
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "%v\n", err)
		flag.PrintDefaults()
		os.Exit(2)
	}

	var reader *bufio.Reader
	var writer *bufio.Writer

	if inputFile == "" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		in, err := os.Open(inputFile)
		if err != nil {
			panic(err)
		}

		defer in.Close()
		reader = bufio.NewReader(in)
	}

	if outputFile == "" {
		writer = bufio.NewWriter(os.Stdout)
	} else {
		var mode int
		if overwrite {
			mode = os.O_TRUNC
		} else {
			mode = os.O_EXCL
		}

		out, err := os.OpenFile(outputFile, os.O_WRONLY | os.O_CREATE | mode,
			0644)
		if err != nil {
			panic(err)
		}

		defer out.Close()
		writer = bufio.NewWriter(out)
	}
	defer writer.Flush()

	count := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		count++
		line = strings.TrimRight(line, "\r\n")

		fmt.Fprintf(writer, "C%010d:", count)
		for _, c := range line {
			fmt.Fprintf(writer, " %-6s", string(c))
		}

		fmt.Fprintln(writer, "")

		fmt.Fprintf(writer, "X%010d:", count)
		for _, c := range line {
			fmt.Fprintf(writer, " 0x%04X", c)
		}

		fmt.Fprintln(writer, "")
	}
}

func checkFlags() error {
	if overwrite && outputFile != "" {
		return errors.New("Flag -overwrite is invalid without flag -output.")
	}

	return nil
}
