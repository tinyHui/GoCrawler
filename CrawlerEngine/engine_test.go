package CrawlerEngine

import (
	"io"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var queue = make(chan url.URL, 10)
var requestURL, _ = url.Parse("http://any.url")
var r = reader(0)

type reader int

func (reader) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (reader) Close() (err error) {
	return nil
}

type mockDownloadEngine struct {
	mock.Mock
}

func (e *mockDownloadEngine) GetURIQueue() chan url.URL {
	return queue
}

func (e *mockDownloadEngine) FetchAndRead(uri string) io.ReadCloser {
	return r
}

type mockLongWaitDownloadEngine struct {
	mock.Mock
}

func (e *mockLongWaitDownloadEngine) GetURIQueue() chan url.URL {
	return queue
}

func (e *mockLongWaitDownloadEngine) FetchAndRead(uri string) io.ReadCloser {
	time.Sleep(5 * time.Second)
	return r
}

type mockHrefScraper struct {
	mock.Mock
}

func (c *mockHrefScraper) GrabLinks(uri url.URL, body io.Reader) []url.URL {
	return []url.URL{
		{Scheme: "http", Host: "monzo.com"},
		{Scheme: "https", Host: "blog.monzo.com"},
	}
}

func Test_processURIPage(t *testing.T) {
	t.Run("should request page", func(t *testing.T) {
		downloadEngine := &mockDownloadEngine{}
		linkScraper := &mockHrefScraper{}
		streamer := &sitemapStreamer{}

		engine := &crawlerEngine{
			downloadEngine:  downloadEngine,
			linkScraper:     linkScraper,
			sitemapStreamer: streamer,
			historyRecord:   &sync.Map{},
		}

		engine.processURIPage(*requestURL)

		uri := <-queue
		assert.Equal(t, "http://monzo.com", uri.String())

		uri = <-queue
		assert.Equal(t, "https://blog.monzo.com", uri.String())
	})

	t.Run("should increment current worker count", func(t *testing.T) {
		downloadEngine := &mockLongWaitDownloadEngine{}
		linkScraper := &mockHrefScraper{}
		streamer := &sitemapStreamer{}

		engine := &crawlerEngine{
			downloadEngine:  downloadEngine,
			linkScraper:     linkScraper,
			sitemapStreamer: streamer,
			historyRecord:   &sync.Map{},
		}

		proceed := make(chan bool, 1)
		go func() {
			proceed <- true
			engine.processURIPage(url.URL{})
			proceed <- true
		}()

		<-proceed
		assert.Equal(t, int64(1), engine.currentWorkerNum)

		<-proceed
		assert.Equal(t, int64(0), engine.currentWorkerNum)
	})
}

func Test_startCrawler(t *testing.T) {
	t.Run("should long waiting time in requesting not cause URI queue timeout", func(t *testing.T) {
		downloadEngine := &mockLongWaitDownloadEngine{}
		linkScraper := &mockHrefScraper{}
		streamer := &sitemapStreamer{}

		engine := &crawlerEngine{
			downloadEngine:  downloadEngine,
			linkScraper:     linkScraper,
			sitemapStreamer: streamer,
			historyRecord:   &sync.Map{},
		}

		done := make(chan bool, 1)
		uriQueue := make(chan url.URL, 1)
		uriQueue <- url.URL{}

		go func() { engine.startCrawler(done, uriQueue) }()

		select {
		case <-done:
			assert.Fail(t, "Exit too quick")
		case <-time.After(4*time.Second + 800*time.Millisecond):
		}
	})
}

func Test_requestURIPage(t *testing.T) {
	t.Run("should mark visited if under same domain", func(t *testing.T) {
		downloadEngine := &mockDownloadEngine{}
		linkScraper := &mockHrefScraper{}
		streamer := &sitemapStreamer{}

		engine := &crawlerEngine{
			downloadEngine:  downloadEngine,
			linkScraper:     linkScraper,
			sitemapStreamer: streamer,
			historyRecord:   &sync.Map{},
			baseDomain:      "monzo.com",
		}

		uri, _ := url.Parse("http://monzo.com")
		engine.requestURIPage(*uri)

		assert.True(t, engine.visited(*uri))
	})

	t.Run("should not mark visited if not under same domain", func(t *testing.T) {
		downloadEngine := &mockDownloadEngine{}
		linkScraper := &mockHrefScraper{}
		streamer := &sitemapStreamer{}

		engine := &crawlerEngine{
			downloadEngine:  downloadEngine,
			linkScraper:     linkScraper,
			sitemapStreamer: streamer,
			historyRecord:   &sync.Map{},
			baseDomain:      "monzo-other.com",
		}

		uri, _ := url.Parse("http://monzo.com")
		engine.requestURIPage(*uri)

		assert.False(t, engine.visited(*uri))
	})
}