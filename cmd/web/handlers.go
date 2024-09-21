package main

import "net/http"

func getHome(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("home"))
}
