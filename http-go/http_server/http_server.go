package main

import (
	"fmt"
	"io"
	"net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request:", r.Method, r.URL.Path)
	fmt.Println("Request headers:", r.Header)
	fmt.Println("Request Content-Length:", r.ContentLength)
	
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	w.Write(requestBody)
}

func main() {
	http.HandleFunc("/echo", echoHandler)
	http.ListenAndServe(":8080", nil)
}
