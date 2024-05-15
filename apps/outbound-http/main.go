package main

import (
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		url := "https://random-data-api.fermyon.app/animals/json"
		resp, err := spinhttp.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, resp.Body)
	})
}

func main() {}
