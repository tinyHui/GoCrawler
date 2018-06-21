package DownloadEngine

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewDownloadEngine(t *testing.T) {
	t.Run("Should able to initial download engine with provided config", func(t *testing.T) {
		config := &DownloaderConfig{
			MaxConcurrentRequest: 100,
		}
		engine := NewDownloadEngine(config)
		queue := engine.GetURIQueue()

		assert.Equal(t, 1000, cap(queue))
	})

	t.Run("Should able to initial download engine with default worker number", func(t *testing.T) {
		config := &DownloaderConfig{}

		engine := NewDownloadEngine(config)
		queue := engine.GetURIQueue()

		assert.Equal(t, DEFAULT_MAX_CONCURRENT_REQUEST*10, cap(queue))
	})
}

func Test_FetchAndRead(t *testing.T) {
	t.Run("Should fetch and read response from server", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, "Any Response From Server")
		}))
		defer ts.Close()

		config := &DownloaderConfig{}
		engine := NewDownloadEngine(config)

		content, _ := ioutil.ReadAll(engine.FetchAndRead(ts.URL))
		assert.Equal(t, "Any Response From Server\n", string(content))
	})

	t.Run("Should get empty response when server not ready", func(t *testing.T) {
		config := &DownloaderConfig{}
		engine := NewDownloadEngine(config)

		content := engine.FetchAndRead("http://any.url.never.exist")
		assert.Nil(t, content)
	})
}