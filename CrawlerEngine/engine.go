package CrawlerEngine

import (
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tinyhui/GoCrawler/DownloadEngine"
	"github.com/tinyhui/GoCrawler/HrefScraper"
	"github.com/tinyhui/GoCrawler/utils"
	"github.com/tinyhui/GoCrawler/utils/log"
)

var logger = log.GetLogger()

type CrawlerEngine interface {
	Start()
	startCrawler(done chan bool, uriQueue chan url.URL)
	requestURIPage(uri url.URL)
	processURIPage(parent url.URL)
	visited(uri url.URL) bool
	markVisited(uri url.URL)
}

func NewCrawlerEngine(parameters *utils.Parameters, initURL url.URL) *crawlerEngine {
	if !HrefScraper.IsValidURI(initURL) {
		logger.Errorln("Please specify a valid url")
		return &crawlerEngine{}
	}

	engine := &crawlerEngine{
		historyRecord:   &sync.Map{},
		downloadEngine:  DownloadEngine.NewDownloadEngine(parameters.DownloaderConfig),
		linkScraper:     HrefScraper.NewScraper(),
		sitemapStreamer: NewSitemapStreamer(parameters.SitemapFilePath),
	}

	engine.baseDomain = getDomain(initURL)

	engine.downloadEngine.GetURIQueue() <- initURL

	return engine
}

type crawlerEngine struct {
	baseDomain       string
	historyRecord    *sync.Map
	downloadEngine   DownloadEngine.DownloadEngine
	linkScraper      HrefScraper.HrefScraper
	sitemapStreamer  SitemapStreamer
	currentWorkerNum int64
}

func (engine *crawlerEngine) Start() {
	logger.Infoln("Crawler Engine Started")

	engine.sitemapStreamer.Init()
	defer engine.sitemapStreamer.End()

	uriQueue := engine.downloadEngine.GetURIQueue()
	defer close(uriQueue)

	done := make(chan bool)
	defer close(done)

	go engine.startCrawler(done, uriQueue)

	<-done
}

func (engine *crawlerEngine) startCrawler(done chan bool, uriQueue chan url.URL) {
	defer func() { done <- true }()

loop:
	for {
		select {
		case uri := <-uriQueue:
			engine.requestURIPage(uri)
		case <-time.After(1 * time.Second):
			if engine.currentWorkerNum > 0 {
				logger.Warnf("Still waiting %d workers to finish", engine.currentWorkerNum)
				continue
			}
			break loop
		}
	}
}

func (engine *crawlerEngine) requestURIPage(uri url.URL) {
	if !engine.underSameDomain(uri) ||
		engine.visited(uri) {
		return
	}
	engine.markVisited(uri)
	engine.sitemapStreamer.NewLoc(uri.String())

	go engine.processURIPage(uri)
}

func (engine *crawlerEngine) processURIPage(parent url.URL) {
	engine.currentWorkerNum = atomic.AddInt64(&engine.currentWorkerNum, int64(1))

	bodyReader := engine.downloadEngine.FetchAndRead(parent.String())
	if bodyReader == nil {
		logger.Errorf("not able to get content from %s, skip", parent)
		return
	}
	defer bodyReader.Close()

	childLinks := engine.linkScraper.GrabLinks(parent, bodyReader)
	for _, childURI := range childLinks {
		if !engine.visited(childURI) {
			engine.downloadEngine.GetURIQueue() <- childURI
		}

		engine.sitemapStreamer.AppendChildLink(childURI.String())
	}

	engine.currentWorkerNum = atomic.AddInt64(&engine.currentWorkerNum, int64(-1))
}
