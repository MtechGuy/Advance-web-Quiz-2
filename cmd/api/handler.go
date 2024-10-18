package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func LoggingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi Welcome to Logging Middleware in Golang\n"))
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Authentication was successful, Welcome to the Authentication Middleware\n"))
}

// helloHandler returns a greeting
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// simulate a panic
	if r.Method != http.MethodPost {
		panic(fmt.Sprintf("Invalid HTTP method %s. Only GET requests are allowed.", r.Method))
	}

	fmt.Fprintln(w, "Hello, and welcome to Error Handling Middleware!")
}

// handleRequest processes incoming requests
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// decode the JSON payload into a User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// print the user's name and email to the console
	fmt.Println("Name:", user.Name)
	fmt.Println("Email:", user.Email)

	// write a response back to the client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, %s!\n", user.Name)
}
