package main

import (
	"github.com/tinyhui/GoCrawler/ArgReader"
	"github.com/tinyhui/GoCrawler/CrawlerEngine"
	"github.com/tinyhui/GoCrawler/utils"
	"github.com/tinyhui/GoCrawler/utils/log"
)


func main() {
	var logger = log.GetLogger()

	initURL, err := ArgReader.GetURLFromArg()
	if err != nil {
		logger.Fatal(err)
	}

	parameters := utils.LoadParameters()

	engine := CrawlerEngine.NewCrawlerEngine(parameters, initURL)
	engine.Start()

	logger.Infoln("All urls have been covered.")
}