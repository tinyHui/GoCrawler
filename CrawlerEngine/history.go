package CrawlerEngine

import "net/url"

func (engine *crawlerEngine) visited(uri url.URL) bool {
	_, v := engine.historyRecord.Load(uri.String())
	return v
}

func (engine *crawlerEngine) markVisited(uri url.URL) {
	engine.historyRecord.Store(uri.String(), true)
}
