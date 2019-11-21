package blobproxy

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gocloud.dev/blob/memblob"
)

func memblobServer() Server {
	b := memblob.OpenBucket(nil)
	b.WriteAll(context.Background(), "pre-existing-file", []byte("pre-existing content"), nil)
	return Server{b}
}

func resReq(method, target string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(method, target, body)
}

func TestExists(t *testing.T) {
	is := is.New(t)
	s := memblobServer()

	res, req := resReq("", "/", nil)
	s.exists(res, req, "pre-existing-file")

	is.Equal(http.StatusOK, res.Code)
}

func TestNotExists(t *testing.T) {
	is := is.New(t)
	s := memblobServer()

	res, req := resReq("", "/", nil)
	s.exists(res, req, "non-existing")

	is.Equal(http.StatusNotFound, res.Code)
}

func TestRead(t *testing.T) {
	is := is.New(t)
	s := memblobServer()

	res, req := resReq("", "/", nil)
	s.read(res, req, "pre-existing-file")

	is.Equal(http.StatusOK, res.Code)
	is.Equal("pre-existing content", res.Body.String())
}

func TestNotRead(t *testing.T) {
	is := is.New(t)
	s := memblobServer()

	res, req := resReq("", "/", nil)
	s.read(res, req, "non-existing")

	is.Equal(http.StatusNotFound, res.Code)
}

func TestWrite(t *testing.T) {
	is := is.New(t)
	s := memblobServer()

	res, req := resReq("", "/", strings.NewReader("new content"))
	s.write(res, req, "new-file")

	is.Equal(http.StatusCreated, res.Code)

	res, req = resReq("", "/", nil)
	s.read(res, req, "new-file")

	is.Equal(http.StatusOK, res.Code)
	is.Equal("new content", res.Body.String())

}

func TestWritePreExisting(t *testing.T) {
	is := is.New(t)
	s := memblobServer()

	res, req := resReq("", "/", strings.NewReader("overwritten content"))
	s.write(res, req, "pre-existing-file")

	is.Equal(http.StatusCreated, res.Code)

	res, req = resReq("", "/", nil)
	s.read(res, req, "pre-existing-file")

	is.Equal(http.StatusOK, res.Code)
	is.Equal("overwritten content", res.Body.String())
}
