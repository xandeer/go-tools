package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"gopkg.in/go-playground/webhooks.v5/github"
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

			fmt.Printf("branch: %+v\n", branch)

			cmdStr := fmt.Sprintf("./deploy.sh ./ %s %s", name, url)
			cmd := exec.Command("/bin/sh", "-c", cmdStr)
			_, err := cmd.Output()

			if err != nil {
				println(err.Error())
				return
			}
		}
	})
	http.ListenAndServe(":3001", nil)
}
