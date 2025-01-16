package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type regexFlags struct {
	keyRegex, keyReplace, valueRegex, valueReplace string
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func exitIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "line2json:", err)
		os.Exit(1)
	}
}

func lineFile() *os.File {
	var (
		file *os.File
		argv []string
		err  error
	)

	argv = flag.Args()
	if len(argv) == 0 {
		return os.Stdin
	}

	file, err = os.Open(argv[0])
	exitIf(err)

	return file
}

func jsonArray(file *os.File, flags *regexFlags) {
	var (
		re       *regexp.Regexp
		scanner  *bufio.Scanner
		array    []string
		jsonData []byte
		err      error
	)

	re, err = regexp.Compile(flags.keyRegex)
	exitIf(err)

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		array = append(array, re.ReplaceAllString(scanner.Text(), flags.keyReplace))
	}

	exitIf(scanner.Err())

	jsonData, err = json.Marshal(array)
	panicIf(err)

	fmt.Println(string(jsonData))
}

func jsonObject(file *os.File, flags *regexFlags) {
	var (
		reKey, reValue *regexp.Regexp
		scanner        *bufio.Scanner
		object         map[string]string
		jsonData       []byte
		err            error
	)

	reKey, err = regexp.Compile(flags.keyRegex)
	exitIf(err)

	reValue, err = regexp.Compile(flags.valueRegex)
	exitIf(err)

	scanner = bufio.NewScanner(file)
	object = make(map[string]string)

	for scanner.Scan() {
		object[reKey.ReplaceAllString(scanner.Text(), flags.keyReplace)] = reValue.ReplaceAllString(scanner.Text(), flags.valueReplace)
	}

	exitIf(scanner.Err())

	jsonData, err = json.Marshal(object)
	panicIf(err)

	fmt.Println(string(jsonData))
}

func main() {
	var (
		objFlag bool
		flags   *regexFlags
	)

	flags = &regexFlags{}

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: line2json [--keyRegex <regex>] [--keyReplace <replacement>] [--object] [--valueRegex <regex>] [--valueReplace <replacement>] [FILE]")
		flag.PrintDefaults()
	}

	flag.StringVar(&flags.keyRegex, "keyRegex", "", "The regular expression to manipulate item if outputting an array or manipulate key if outputting an object.")
	flag.StringVar(&flags.keyReplace, "keyReplace", "", "The replacement string when key regular expression matches.")
	flag.BoolVar(&objFlag, "object", false, "Whether to convert each line to a key-value object.")
	flag.StringVar(&flags.valueRegex, "valueRegex", "", "The regular expression to manipulate value if outputting an object.")
	flag.StringVar(&flags.valueReplace, "valueReplace", "", "The replacement string when value regular expression matches.")
	flag.Parse()

	if objFlag {
		jsonObject(lineFile(), flags)

		return
	}

	jsonArray(lineFile(), flags)
}
