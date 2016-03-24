package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const Name string = "git-extern-sync"

func main() {
	if err := gitExternSyncFromPath(".gitignore"); err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %v\n", Name, err)
		os.Exit(1)
	}
}

func gitExternSyncFromPath(p string) (err error) {
	var file *os.File
	if file, _ = os.Open(p); file == nil {
		fmt.Fprintf(os.Stderr, "%s: %s not found\n", Name, p)
		return
	}
	defer file.Close()

	var blob []byte
	if blob, err = ioutil.ReadAll(file); err != nil {
		return
	}

	err = gitExternSyncFromBlob(blob)
	return
}

func gitExternSyncFromBlob(blob []byte) (err error) {
	var wd string
	if wd, err = os.Getwd(); err != nil {
		return
	}

	var regex *regexp.Regexp
	if regex, err = regexp.Compile(`^#\s*sync:(.*)$`); err != nil {
		return
	}

	var uri *string
	scanner := bufio.NewScanner(bytes.NewBuffer(blob))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case uri != nil:
			if err = installUri(*uri, wd, line, false); err != nil {
				return
			}
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

func installUri(uri, wd, p string) (err error) {
	fpath := filepath.Join(wd, p)

	// ensure nested directories are created (mkdir -p)
	dirpath := filepath.Dir(fpath)
	if err = os.MkdirAll(dirpath, os.ModeDir|0755); err != nil {
		return
	}

	var file *os.File
	var prevData, data []byte
	var prevMD5 *string
	var MD5 string

	// if file exists, compare MD5s and prompt before overwriting
	if file, _ = os.Open(fpath); file != nil {
		defer file.Close()
		if prevData, err = ioutil.ReadAll(file); err != nil {
			return
		}
		_prevMD5 := fmt.Sprintf("%x", md5.Sum(prevData))
		prevMD5 = &_prevMD5
	}

	var resp *http.Response
	if resp, err = http.Get(uri); err != nil {
		return
	}
	defer resp.Body.Close()
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	MD5 = fmt.Sprintf("%x", md5.Sum(data))

	printed := false
	// force type variable would help for testing this
	if prevMD5 != nil && *prevMD5 != MD5 {
		promptFmt := "%s: mismatch: %s (data from %s). Overwrite? [y/N] "
		fmt.Printf(promptFmt, Name, p, uri)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		resp := scanner.Text()
		if resp == "y" || resp == "Y" {
			fmt.Printf("%s: overwrote: %s\n", Name, p)
			printed = true
		} else {
			fmt.Printf("%s: skipped: %s\n", Name, p)
			return
		}

	}

	if file, err = os.Create(fpath); err != nil {
		return
	}
	defer file.Close()

	if _, err = file.Write(data); err != nil {
		return
	}
	if !printed {
		fmt.Printf("%s: synchronized: %s\n", Name, p)
	}

	return
}
