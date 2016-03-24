package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGitExternSyncFromBlob(t *testing.T) {
	defer chtmpdir(t)()
	wd, _ := os.Getwd()

	gitExternSyncFromBlob([]byte(``))
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Errorf("empty blob generated files")
	}

	needle := generateNeedle(32)
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, needle)
	}
	ts := httptest.NewServer(http.HandlerFunc(fn))
	defer ts.Close()

	fname := generateNeedle(32)
	gitExternSyncFromBlob([]byte("#sync:" + ts.URL + "\n" + fname))
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if string(data) != needle {
		t.Errorf("expected %q, got %q", needle, data)
	}
}

func generateNeedle(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// https://golang.org/src/os/os_test.go
func chtmpdir(t *testing.T) func() {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	d, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	if err := os.Chdir(d); err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	return func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("chtmpdir: %v", err)
		}
		os.RemoveAll(d)
	}
}
