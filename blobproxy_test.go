package blobproxy

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/matryer/is"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

func TestSigner(t *testing.T) {
	is := is.New(t)

	bucketFolder := "./testdata/temp"
	os.RemoveAll(bucketFolder)
	os.MkdirAll(bucketFolder, os.ModePerm)

	baseURL, err := url.Parse("http://localhost/")
	is.NoErr(err)

	signer := fileblob.NewURLSignerHMAC(baseURL, []byte("secret"))
	bucket, err := fileblob.OpenBucket(bucketFolder, &fileblob.Options{
		URLSigner: signer,
	})
	is.NoErr(err)

	signedURL, err := bucket.SignedURL(context.Background(), "123", &blob.SignedURLOptions{
		Method: "PUT",
	})
	is.NoErr(err)

	// Correct: 201 created
	res, req := resReq("PUT", signedURL, nil)
	Server{bucket}.URLSignerHMACHandler(signer, false).ServeHTTP(res, req)

	is.Equal(201, res.Code)

	// Wrong method: 405
	res, req = resReq("GET", signedURL, nil)
	Server{bucket}.URLSignerHMACHandler(signer, false).ServeHTTP(res, req)

	is.Equal(405, res.Code)

	// Wrong method, but public GET: 200
	res, req = resReq("GET", signedURL, nil)
	Server{bucket}.URLSignerHMACHandler(signer, true).ServeHTTP(res, req)

	is.Equal(200, res.Code)

	// Wrong signature: 401
	res, req = resReq("PUT", "/some-url", nil)
	Server{bucket}.URLSignerHMACHandler(signer, false).ServeHTTP(res, req)

	is.Equal(401, res.Code)
}

func TestNilSigner(t *testing.T) {
	is := is.New(t)

	bucketFolder := "./testdata/temp"
	os.RemoveAll(bucketFolder)
	os.MkdirAll(bucketFolder, os.ModePerm)

	bucket, err := fileblob.OpenBucket(bucketFolder, nil)
	is.NoErr(err)

	// Correct: 201 created
	res, req := resReq("POST", "/some-url?method=POST&obj=345", nil)
	Server{bucket}.URLSignerHMACHandler(nil, false).ServeHTTP(res, req)

	is.Equal(201, res.Code)

	// Wrong method: 405
	res, req = resReq("PUT", "/some-url?method=POST&obj=345", nil)
	Server{bucket}.URLSignerHMACHandler(nil, false).ServeHTTP(res, req)

	is.Equal(405, res.Code)
}
