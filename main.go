package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	output := flag.CommandLine.Output()
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Usage: "+os.Args[0]+" [OPTIONS] FILE [FILE...]")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "cut for tsv/csv")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Options:")
	flag.CommandLine.PrintDefaults()
}

func cut(w io.Writer, input string, keys []string, delimiter rune) error {
	var file io.ReadCloser
	var r *csv.Reader
	var err error
	if input == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(input)
		if err != nil {
			return err
		}
	}
	fmt.Println("hello")
	r = csv.NewReader(file)
	r.Comma = delimiter
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	keyMap := map[string]int{}
	header := records[0]
	for i, col := range header {
		keyMap[col] = i
}
	for _, line := range records {
		for i, key := range keys {
			if i != 0 {
				fmt.Print(string(delimiter))
			}
			fmt.Print(line[keyMap[key]])
		}
		fmt.Println()
	}
	return nil
}

func main() {
	flag.Usage = usage
	output := flag.CommandLine.Output()

	var fields string
	var delimiter string
	var version, help bool

	flag.StringVar(&fields, "f", "", "specify column name to select only these fields")
	flag.StringVar(&delimiter, "d", ",", "specify delimiter")
	flag.BoolVar(&version, "v", false, "show version")
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()

	if help {
		usage()
		return
	}

	if version {
		fmt.Fprintln(output, "1.0.0")
		return
	}

	args := flag.Args()
	if len(args) <= 0 {
		args = append(args, "-")
	}

	var keys = strings.Split(fields, ",")
	fmt.Println(keys)

	if delimiter == "\\t" {
		delimiter = "\t"
	}

	for _, arg := range args {
		err := cut(os.Stdout, arg, keys, []rune(delimiter)[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
