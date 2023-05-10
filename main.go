package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	app := mux.NewRouter()
	
	app.HandleFunc("/hello/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Println(r.Method, id)
	}).Methods(http.MethodGet)

	http.ListenAndServe(":8888", app)
}