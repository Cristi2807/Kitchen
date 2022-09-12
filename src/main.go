package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getOrder(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func main() {
	http.HandleFunc("/", getOrder)
	err := http.ListenAndServe(":8010", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed \n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s \n", err)
		os.Exit(1)
	}
}
