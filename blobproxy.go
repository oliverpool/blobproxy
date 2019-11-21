package blobproxy

import (
	"net/http"

	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

type Server struct {
	bucket *blob.Bucket
}

func (s Server) URLSignerHMACHandler(signer fileblob.URLSigner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := signer.KeyFromURL(r.Context(), r.URL)
		if err != nil {
			http.Error(w, "wrong signature", http.StatusUnauthorized)
			return
		}
		method := r.FormValue("method")
		if method != r.Method {
			http.Error(w, "http method should be the same as the one in the query", http.StatusMethodNotAllowed)
			return
		}
		s.dispatch(w, r, method, key)
	}
}
