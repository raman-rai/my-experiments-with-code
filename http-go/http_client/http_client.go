package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {

	user := User{Name: "Alice", Age: 30}
	userJSON, _ := json.Marshal(user)
	
	response, err := http.Post("http://localhost:8080/echo", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	println("Response status:", response.Status)
	println("Response body:", string(responseBody))
}