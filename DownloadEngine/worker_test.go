package DownloadEngine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetch(t *testing.T) {
	t.Run("Should able to fetch response from given url when server ready", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, world")
		}))
		defer ts.Close()

		response := fetch(ts.URL)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, int64(len("Hello, world\n")), response.ContentLength)
	})

	t.Run("Should get nil when server not found", func(t *testing.T) {
		response := fetch("http://someurl.never.exist")

		assert.Nil(t, response)
	})
}

func Test_read(t *testing.T) {
	t.Run("Should read content from response", func(t *testing.T) {
		header := make(map[string][]string)
		header["Content-Type"] = []string{"text/html; charset=ISO-8859-1"}

		testResp := &http.Response{
			Header: header,
			Body:   ioutil.NopCloser(bytes.NewBufferString("Hello, world")),
		}

		body := read("http://any.url", testResp)
		content, _ := ioutil.ReadAll(body)

		assert.Equal(t, "Hello, world", string(content))
	})

	t.Run("Should return nil when send in response is nil", func(t *testing.T) {
		body := read("http://any.url", nil)
		assert.Nil(t, body)
	})

	t.Run("Should return empty content when response type is not HTML", func(t *testing.T) {
		header := make(map[string][]string)
		header["Content-Type"] = []string{"application/javascript"}
		testResp := &http.Response{
			Header: header,
			Body:   ioutil.NopCloser(bytes.NewBufferString("Hello, world")),
		}

		body := read("http://any.url", testResp)
		assert.Nil(t, body)
	})

	t.Run("Should return nil when response type is empty in header", func(t *testing.T) {
		testResp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("Hello, world")),
		}

		body := read("http://any.url", testResp)
		assert.Nil(t, body)
	})
}

func Test_isHtmlType(t *testing.T) {
	t.Run("Should return true when type is text/html", func(t *testing.T) {
		assert.True(t, isHtmlType("text/html"))
	})

	t.Run("Should return true when type is text/html with charset", func(t *testing.T) {
		assert.True(t, isHtmlType(" text/html; charset=utf-8"))
	})

	t.Run("Should return false when type is zip", func(t *testing.T) {
		assert.False(t, isHtmlType("application/zip, application/octet-stream"))
	})
}

