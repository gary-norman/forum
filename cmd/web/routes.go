package main

import (
  "net/http"
)

func routes() http.Handler {
  mux := http.NewServeMux()
  mux.HandleFunc("GET /", getHome)
  return mux
}
