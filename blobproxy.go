package blobproxy

import (
	"net/http"

	"gocloud.dev/blob"
)

type Server struct {
	bucket *blob.Bucket
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	if key[0] == '/' {
		key = key[1:]
	}
	if r.Method == "GET" {
		s.read(w, r, key)
	} else if r.Method == "HEAD" {
		s.exists(w, r, key)
	} else if r.Method == "POST" {
		s.write(w, r, key)
	} else if r.Method == "PUT" {
		s.write(w, r, key)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
