package main

import (
	"fmt"
	"net/http"
)

func main() {
	app := http.NewServeMux()
	
	app.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.Method)
	})

	http.ListenAndServe(":8888", app)
}