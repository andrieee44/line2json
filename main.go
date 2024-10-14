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
		scanner  *bufio.Scanner
		re       *regexp.Regexp
		array    []string
		jsonData []byte
		err      error
	)

	scanner = bufio.NewScanner(file)
	re = regexp.MustCompile(flags.key)

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
		scanner        *bufio.Scanner
		reKey, reValue *regexp.Regexp
		object         map[string]string
		jsonData       []byte
		err            error
	)

	scanner = bufio.NewScanner(file)
	reKey = regexp.MustCompile(flags.key)
	reValue = regexp.MustCompile(flags.value)
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
