package main

import (
	"fmt"
	"net/http"

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
			fmt.Printf("%+v\n", push)
		}
	})
	http.ListenAndServe(":3001", nil)
}
