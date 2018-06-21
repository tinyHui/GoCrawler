package DownloadEngine

import (
	"io"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/tinyhui/GoCrawler/utils/log"
)

var logger = log.GetLogger()

type DownloadEngine interface {
	GetURIQueue() chan url.URL
	FetchAndRead(uri string) io.ReadCloser
}

func NewDownloadEngine(config *DownloaderConfig) *downloadEngine {
	maxConcurrentRequest := config.MaxConcurrentRequest
	if maxConcurrentRequest == 0 {
		maxConcurrentRequest = DEFAULT_MAX_CONCURRENT_REQUEST
	}

	return &downloadEngine{
		uriQueue:                 make(chan url.URL, maxConcurrentRequest*10),
		currentConcurrentRequest: int32(0),
		maxMaxConcurrentRequest:  int32(maxConcurrentRequest),
	}
}

type downloadEngine struct {
	uriQueue                 chan url.URL
	currentConcurrentRequest int32
	maxMaxConcurrentRequest  int32
}

func (engine *downloadEngine) GetURIQueue() chan url.URL {
	return engine.uriQueue
}

func (engine *downloadEngine) FetchAndRead(uri string) io.ReadCloser {
	for engine.currentConcurrentRequest > engine.maxMaxConcurrentRequest {
		time.Sleep(1 * time.Millisecond)
	}

	atomic.AddInt32(&engine.currentConcurrentRequest, 1)
	defer atomic.AddInt32(&engine.currentConcurrentRequest, -1)

	logger.Infof("Fetch and read %s", uri)

	bodyReader := fetch(uri)
	return read(uri, bodyReader)
}
