package main

import (
	"fmt"
	"net/http"

	"github.com/codercola034/notion-ai/notion"
)

func httpServer(port int) error {
	fmt.Println("Running http server on port", port)
	// a simple http server that returns the prompt response
	http.HandleFunc("/prompt", func(w http.ResponseWriter, r *http.Request) {
		// read the prompt from the request body
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		// get the response from the prompt
		res, err := notion.GetCompletion(string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// write the response to the http response
		for out, ended, err := res.Output(); !ended; out, ended, err = res.Output() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, out)
		}
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
