package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getOrder(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/order" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Printf("got /order request\n")

	postBody, _ := json.Marshal(map[string]string{
		"name":  "Toby",
		"email": "Toby@example.com",
	})
	responseBody := bytes.NewBuffer(postBody)

	http.Post("http://dinninghall:8020/distribution", "application/json", responseBody)

}

func main() {
	http.HandleFunc("/order", getOrder)

	fmt.Printf("Server Kitchen started on PORT 8010\n")
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal(err)
	}
}
