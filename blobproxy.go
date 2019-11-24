package blobproxy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

type Server struct {
	Bucket *blob.Bucket
}

func (s Server) URLSignerHMACHandler(signer fileblob.URLSigner, publicGET bool) http.HandlerFunc {
	keyFromURL := func(_ context.Context, sURL *url.URL) (string, error) {
		q := sURL.Query()
		return q.Get("obj"), nil
	}
	if signer != nil {
		keyFromURL = signer.KeyFromURL
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var key string

		if publicGET && r.Method == "GET" {
			key = r.URL.Query().Get("obj")
		} else {
			var err error
			key, err = keyFromURL(r.Context(), r.URL)
			if err != nil {
				http.Error(w, "wrong signature", http.StatusUnauthorized)
				return
			}
			method := r.FormValue("method")
			if method != r.Method {
				http.Error(w, fmt.Sprintf("http method (%s) should be the same as in the query (%s)", r.Method, method), http.StatusMethodNotAllowed)
				return
			}
		}
		s.dispatch(w, r, r.Method, key)
	}
}
