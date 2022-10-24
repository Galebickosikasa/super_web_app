package main

import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "kek")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "there will be a login page")
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/login/", loginPage)
	http.ListenAndServe(":3000", nil)
}

func main() {
	handleRequest()
}
