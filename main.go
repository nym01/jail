package main

import (
	"fmt"
	"net/http"
)

func runHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from jail")
}

func main() {
	http.HandleFunc("/run", runHandler)
	http.ListenAndServe(":8080", nil)
}
