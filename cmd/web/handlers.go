package main

import (
	"html/template"
	"net/http"
)

func getHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./assets/templates/index.html")
	if err != nil {
		//TODO log error
		http.Error(w, err.Error(), 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		//TODO log error
		return
	}
}
