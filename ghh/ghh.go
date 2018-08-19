package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/go-playground/webhooks.v5/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	hook, _ := github.New(github.Options.Secret("yuEaNL74CuAw"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
				fmt.Println("Event not found.")
			}
		}
		switch payload.(type) {
		case github.PushPayload:
			push := payload.(github.PushPayload)
			branch := strings.TrimPrefix(push.Ref, "refs/heads/")
			name := push.Repository.Name
			url := push.Repository.CloneURL

			_, err := git.PlainClone(name, false, &git.CloneOptions{
				URL:               url,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})

			if err == nil {
				make(name)
				return
			}

			if err == git.ErrRepositoryAlreadyExists {
				r, e := git.PlainOpen(name)
				if e == nil {
					w, err := r.Worktree()

					if err == nil {
						err = w.Checkout(&git.CheckoutOptions{
							Hash: plumbing.NewHash(branch),
						})

						if err == nil {
							make(name)
						}
					}
				}
			}
		}
	})
	http.ListenAndServe(":3001", nil)
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
	fmt.Printf("cmd dir: %s\n", cmd.Dir)
	fmt.Printf("cmd: %v\n", cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		println(err)
	}
	fmt.Printf("Output: %s\n", out.String())
}
