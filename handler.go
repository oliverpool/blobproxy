package blobproxy

import (
	"io"
	"net/http"

	"gocloud.dev/gcerrors"
)

var errStatues = map[gcerrors.ErrorCode]int{
	gcerrors.NotFound: http.StatusNotFound,
}

func (Server) handleError(w http.ResponseWriter, r *http.Request, err error) {
	status := errStatues[gcerrors.Code(err)]
	if status == 0 {
		status = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), status)
	return

}

func (s Server) exists(w http.ResponseWriter, r *http.Request, key string) {
	ok, err := s.bucket.Exists(r.Context(), key)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}

func (s Server) read(w http.ResponseWriter, r *http.Request, key string) {
	reader, err := s.bucket.NewReader(r.Context(), key, nil)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	defer reader.Close()
	_, err = io.Copy(w, reader)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s Server) write(w http.ResponseWriter, r *http.Request, key string) {
	writer, err := s.bucket.NewWriter(r.Context(), key, nil)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	defer writer.Close()
	_, err = io.Copy(writer, r.Body)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	err = writer.Close()
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
