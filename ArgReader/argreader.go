package ArgReader

import (
	"net/url"
	"os"
	"strings"

	"github.com/tinyhui/GoCrawler/HrefScraper"
	"github.com/tinyhui/GoCrawler/utils/log"
)

var logger = log.GetLogger()

func GetURLFromArg() (url.URL, error) {
	args := os.Args

	return getURI(args)
}

func getURI(args []string) (url.URL, error) {
	if len(args) < 2 {
		return url.URL{}, NewInsufficientArgumentError()
	}

	if len(args) > 2 {
		logger.Warningln("More then url is provided inside arguments, only take the first argument")
	}

	urlString := args[1]
	if strings.TrimSpace(urlString) == "" {
		return url.URL{}, NewInsufficientArgumentError()
	}

	uri, err := HrefScraper.StrConvURL(urlString)
	if err != nil {
		return url.URL{}, err
	}

	return uri, nil
}