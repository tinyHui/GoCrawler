package HrefScraper

import (
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/tinyhui/GoCrawler/utils/log"
	"golang.org/x/net/html"
)

var logger = log.GetLogger()
var hashPattern = regexp.MustCompile(`#.*$`)

type HrefScraper interface {
	GrabLinks(parent url.URL, body io.Reader) []url.URL
}

func NewScraper() *hrefScraper {
	return &hrefScraper{}
}

type hrefScraper struct {
}

func (scraper *hrefScraper) GrabLinks(parent url.URL, body io.Reader) []url.URL {
	tokenizer := html.NewTokenizer(body)
	var links []url.URL

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return deduplicate(links)
		}

		token := tokenizer.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					link := strings.TrimSpace(attr.Val)
					link = trimHash(link)
					uri, err := complete(parent, link)
					if err != nil {
						logger.Warnf("Not able to convert %s into URL", attr.Val)
						continue
					}

					if IsValidURI(uri) {
						links = append(links, uri)
					} else {
						logger.Warnf("%s is not a valid URL", attr.Val)
					}
				}
			}
		}
	}
}
