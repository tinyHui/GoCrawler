package utils

import (
	"io/ioutil"
	"os"

	"github.com/tinyhui/GoCrawler/DownloadEngine"
	"github.com/tinyhui/GoCrawler/utils/log"
	"gopkg.in/yaml.v2"
)

var logger = log.GetLogger()

type Parameters struct {
	DownloaderConfig *DownloadEngine.DownloaderConfig `yaml:"downloader"`
	SitemapFilePath  string                           `yaml:"sitemap_dir"`
}

func LoadParameters() *Parameters {
	parametersFile := os.Getenv("config")
	if parametersFile == "" {
		logger.Fatalln("Config file path missing")
	}

	yamlFile, err := ioutil.ReadFile(parametersFile)
	if err != nil {
		logger.Fatalf("configFile %s .Get err #%v", parametersFile, err)
	}

	parameters := Parameters{}
	err = yaml.Unmarshal(yamlFile, &parameters)
	if err != nil {
		logger.Fatalf("Unmarshal: %v", err)
	}

	return &parameters
}
