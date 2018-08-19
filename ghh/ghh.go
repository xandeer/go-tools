package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/go-playground/webhooks.v5/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	port     int
	dir      string
	secret   string
	branches string
)

func main() {
	bindFlags()

	dir = fixRelativeDir(dir)

	hook, _ := github.New(github.Options.Secret(secret))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("start %s\n", start.Format("06-01-02 15:04:05"))
		payload, err := hook.Parse(r, github.PushEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
				fmt.Println("Event not found.")
			}
			fmt.Printf("parse failed, %+v\n", err)
			return
		}
		switch payload.(type) {
		case github.PushPayload:
			push := payload.(github.PushPayload)
			branch := strings.TrimPrefix(push.Ref, "refs/heads/")
			name := push.Repository.Name
			url := push.Repository.CloneURL

			if !strings.Contains(branches, branch) {
				fmt.Printf("enabled branches: %s, current branch: %s\n", branches, branch)
				return
			}

			clonedPath := filepath.Join(dir, name)
			_, err := git.PlainClone(clonedPath, false, &git.CloneOptions{
				URL:               url,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})

			if err == nil {
				make(clonedPath)
				return
			}

			if err == git.ErrRepositoryAlreadyExists {
				r, e := git.PlainOpen(name)

				fmt.Println("open repo " + name)
				if e != nil {
					fmt.Printf("open repo failed, %+v\n", e)
					return
				}
				w, err := r.Worktree()

				fmt.Println("work repo " + name)
				if err != nil {
					fmt.Printf("work repo failed, %+v\n", err)
					return
				}

				fmt.Println("start pull")
				err = w.Pull(&git.PullOptions{RemoteName: "origin"})

				if err != nil && err != git.NoErrAlreadyUpToDate {
					fmt.Printf("pull failed, %+v\n", err)
					return
				}

				fmt.Println("checkout branch " + branch)
				err = w.Checkout(&git.CheckoutOptions{
					Hash: plumbing.NewHash(push.After),
				})

				if err != nil {
					fmt.Printf("checkout branch failed, %+v\n", err)
					return
				}

				make(clonedPath)
				return
			}

			fmt.Printf("clone failed, %+v\n", err)
		}
	})
	fmt.Printf("server running at 0.0.0.0:%v\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func fixRelativeDir(dir string) string {
	if !(strings.HasPrefix(dir, "~") || strings.HasPrefix(dir, "/")) {
		d, _ := os.Getwd()
		return filepath.Join(d, dir)
	}
	return dir
}

func make(dir string) {
	cmd := exec.Command("/bin/sh", "-c", "cd "+dir+" && pwd && make")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("cmd exec failed: %+v\n", err)
	}
	fmt.Printf("cmd output: %s\n", out.String())
}

func bindFlags() {
	flag.IntVar(&port, "p", 3001, "Port")
	flag.StringVar(&dir, "d", "./", "Workspace directory")
	flag.StringVar(&secret, "s", "", "Secret")
	flag.StringVar(&branches, "b", "master", "Listen branches, split by ','")
	flag.Parse()

	if secret == "" {
		panic("secret not set")
	}
}
