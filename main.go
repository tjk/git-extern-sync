package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const Name string = "git-extern-sync"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var file *os.File
	if file, _ = os.Open("./.gitignore"); file == nil {
		fmt.Fprintf(os.Stderr, "%s: no ./gitignore found\n", Name)
		return
	}
	defer file.Close()

	var wd string
	if wd, err = os.Getwd(); err != nil {
		return
	}

	var regex *regexp.Regexp
	if regex, err = regexp.Compile(`^#\s*sync:(.*)$`); err != nil {
		return
	}

	var uri *string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case uri != nil:
			fpath := filepath.Join(wd, line)
			dirpath := filepath.Dir(fpath)
			if err = os.MkdirAll(dirpath, os.ModeDir|0755); err != nil {
				return
			}
			var file *os.File
			if file, err = os.Create(fpath); err != nil {
				return
			}
			defer file.Close()
			var resp *http.Response
			if resp, err = http.Get(*uri); err != nil {
				return
			}
			defer resp.Body.Close()
			var body []byte
			if body, err = ioutil.ReadAll(resp.Body); err != nil {
				return
			}
			if _, err = file.Write(body); err != nil {
				return
			}
			fmt.Printf("%s: synchronized %s\n", Name, fpath)
			uri = nil
		default:
			if res := regex.FindStringSubmatch(line); res != nil {
				uri = &res[1]
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	return
}
