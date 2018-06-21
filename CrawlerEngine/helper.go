package CrawlerEngine

import (
	"net/url"
	"regexp"
	"strings"
)

var portPattern, _ = regexp.Compile(":.*$")

func (engine *crawlerEngine) underSameDomain(uri url.URL) bool {
	return getDomain(uri) == engine.baseDomain
}

func getDomain(uri url.URL) string {
	if strings.Contains(uri.Host, ":") {
		uri.Host = portPattern.ReplaceAllString(uri.Host, "")
	}

	slice := strings.Split(uri.Host, ".")
	if len(slice) < 2 {
		return ""
	}
	return strings.Join(slice[len(slice)-2:], ".")
}

