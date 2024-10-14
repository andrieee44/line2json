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
	key, value, keyRep, valueRep string
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

	re, err = regexp.Compile(flags.key)
	exitIf(err)

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		array = append(array, re.ReplaceAllString(scanner.Text(), flags.keyRep))
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

	reKey, err = regexp.Compile(flags.key)
	exitIf(err)

	reValue, err = regexp.Compile(flags.value)
	exitIf(err)

	scanner = bufio.NewScanner(file)
	object = make(map[string]string)

	for scanner.Scan() {
		object[reKey.ReplaceAllString(scanner.Text(), flags.keyRep)] = reValue.ReplaceAllString(scanner.Text(), flags.valueRep)
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
		fmt.Fprintln(flag.CommandLine.Output(), `usage: line2json [OPTION]... [FILE]

line2json converts text lines to a JSON array.
if no FILE, read from STDIN.

example: line2json test.txt`)

		flag.PrintDefaults()
	}

	flag.BoolVar(&objFlag, "o", false, "convert each line to a key value object")
	flag.StringVar(&flags.key, "K", "", "regular expression to manipulate item if outputting an array, key if outputting an object")
	flag.StringVar(&flags.value, "V", "", "regular expression to manipulate value if outputting an object")
	flag.StringVar(&flags.keyRep, "k", "", "string to replace when key regular expression matches")
	flag.StringVar(&flags.valueRep, "v", "", "string to replace when value regular expression matches")
	flag.Parse()

	if objFlag {
		jsonObject(lineFile(), flags)

		return
	}

	jsonArray(lineFile(), flags)
}
