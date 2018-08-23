package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	ip   string
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

	fileHandler := http.FileServer(http.Dir(dir))
	http.Handle("/", withLogging(fileHandler))
	http.Handle("/u", withLogging(http.HandlerFunc(upload)))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func withLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		boldwhite.Printf(
			"%s\t%s\t%s\t%s\n",
			start.Format("06-01-02 15:04:05"),
			r.RemoteAddr,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func bindFlags() {
	flag.IntVar(&port, "p", 9999, "Port")
	flag.StringVar(&dir, "d", "./", "Directory")
	flag.StringVar(&ip, "i", "127.0.0.1", "IP or domain name")
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

func upload(w http.ResponseWriter, r *http.Request) {
	tpl := `<html>
			<head>
				<title>Upload a File</title>
			</head>
			<body>
				<form enctype="multipart/form-data" action="http://%v:%v/u" method="post">
					<input type="file" name="uploadfile" required="true" />
					<input type="submit" value="upload file" />
				</form>
			</body>
		</html>`
	html := []byte(fmt.Sprintf(tpl, ip, port))
	if r.Method == "GET" {
		w.Write(html)
	} else if r.Method == "POST" {
		start := time.Now()

		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			boldred.Println(err)
			return
		}
		defer file.Close()

		path := dir + "/upload/"
		os.MkdirAll(path, os.ModePerm)
		f, err := os.OpenFile(path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			boldred.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		boldwhite.Printf(
			"%s\t%s\t%s\tupload file: %s\t%s\n",
			start.Format("06-01-02 15:04:05"),
			r.RemoteAddr,
			r.RequestURI,
			handler.Filename,
			time.Since(start),
		)
		w.Write(html)
	}
}
