package HrefScraper

import (
	"net/url"
	"path"
	"strings"
)

func StrConvURL(urlString string) (url.URL, error) {
	uri, err := url.Parse(urlString)
	if err != nil {
		return url.URL{}, err
	}

	if uri.Host == "" && strings.Contains(uri.Path, ".") {
		slice := strings.SplitN(uri.Path, "/", 2)

		uri.Host = slice[0]
		uri.Path = "/"

		if len(slice) > 1 {
			uri.Path = "/" + slice[1]
		}
	}

	if uri.Scheme == "" {
		uri.Scheme = "http"
	}

	return *uri, nil
}

func IsValidURI(uri url.URL) bool {
	if uri.Host != "" && (
		uri.Scheme == "http" ||
			uri.Scheme == "https" ||
			uri.Scheme == "") {
		return true
	}
	return false
}

func trimHash(link string) string {
	if strings.Contains(link, "#") {
		return hashPattern.ReplaceAllString(link, "")
	}
	return link
}

func complete(parent url.URL, link string) (url.URL, error) {
	if link == "" {
		return parent, nil
	} else if strings.Contains(link, "://") {
		return StrConvURL(link)
	} else if strings.HasPrefix(link, "/") {
		newLink := url.URL(parent)
		newLink.Path = link
		return newLink, nil
	} else {
		newLink := url.URL(parent)
		newLink.Path = path.Join(parent.Path, link)
		return newLink, nil
	}
}

func deduplicate(links []url.URL) []url.URL {
	linksRecord := make(map[url.URL]bool)
	var newLinks []url.URL

	for _, link := range links {
		if _, value := linksRecord[link]; !value {
			linksRecord[link] = true
			newLinks = append(newLinks, link)
		}
	}
	return newLinks
}
