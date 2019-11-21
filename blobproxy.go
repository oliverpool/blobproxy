package blobproxy

import (
	"context"
	"net/http"
	"net/url"

	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

type Server struct {
	Bucket *blob.Bucket
}

func (s Server) URLSignerHMACHandler(signer fileblob.URLSigner) http.HandlerFunc {
	keyFromURL := func(_ context.Context, sURL *url.URL) (string, error) {
		q := sURL.Query()
		return q.Get("obj"), nil
	}

	if signer != nil {
		keyFromURL = signer.KeyFromURL
	}

	return func(w http.ResponseWriter, r *http.Request) {
		key, err := keyFromURL(r.Context(), r.URL)
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
