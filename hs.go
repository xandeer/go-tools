package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	port int
	dir  string
)

var (
	green     = color.New(color.FgHiGreen)
	boldcyan  = color.New(color.FgCyan, color.Bold)
	boldred   = color.New(color.FgRed, color.Bold)
	boldwhite = color.New(color.FgHiWhite, color.Bold)
)

func main() {
	bindFlags()

	dir = fixRelativeDir(dir)

	if err := validateDir(dir); err != nil {
		boldred.Print(err)
		return
	}

	green.Printf("path: %s\n", http.Dir(dir))
	boldwhite.Print("Server running at ")
	boldcyan.Printf("http://0.0.0.0:%d\n", port)

	http.Handle("/", http.FileServer(http.Dir(dir)))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func bindFlags() {
	flag.IntVar(&port, "p", 9999, "Port")
	flag.StringVar(&dir, "d", "./", "Directory")
	flag.Parse()
}

func fixRelativeDir(dir string) string {
	if !(strings.HasPrefix(dir, "~") || strings.HasPrefix(dir, "/")) {
		d, _ := os.Getwd()
		return filepath.Join(d, dir)
	}
	return dir
}

func validateDir(dir string) error {
	fs, err := os.Stat(dir)
	if os.IsNotExist(err) {
		e := fmt.Sprintf("%s: no such directory\n", dir)
		return errors.New(e)
	}
	if !fs.IsDir() {
		e := fmt.Sprintf("%s: not a directory\n", dir)
		return errors.New(e)
	}
	return nil
}
