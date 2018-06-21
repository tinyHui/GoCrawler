package DownloadEngine

type DownloaderConfig struct {
	MaxConcurrentRequest uint `yaml:"maxConcurrentRequest"`
}

const DEFAULT_MAX_CONCURRENT_REQUEST = 10