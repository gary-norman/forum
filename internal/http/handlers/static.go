package handlers

import (
	"fmt"
	"net/http"

	"github.com/gary-norman/forum/internal/http/middleware"
)

func staticHandler(location string) http.Handler {
	dir := fmt.Sprintf(("./%v"), location)
	prefix := fmt.Sprintf("/%v/", location)
	fs := http.FileServer(http.Dir(dir))
	return http.StripPrefix(prefix, fs)
}

func MuxHandler(mux *http.ServeMux, address string) {
	// Logging middleware wrapper
	use := func(h http.Handler) http.Handler {
		return middleware.Logging(h)
	}
	address = fmt.Sprintf("/%v/", address)
	mux.Handle(address, use(staticHandler(address)))
}
