package DownloadEngine

import (
	"io"
	"net/http"
	"strings"
)

func fetch(uri string) *http.Response {
	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("%v", err)
		return nil
	}

	statusCode := resp.StatusCode
	if statusCode != 200 {
		logger.Warnf("Get %s, got status code %d", uri, statusCode)
	}

	return resp
}

func read(uri string, resp *http.Response) io.ReadCloser {
	if resp == nil {
		logger.Warnf("Get %s: Response is nil, return empty content", uri)
		return nil
	}

	if !isHtmlType(resp.Header.Get("Content-Type")) {
		logger.Infof("Get %s: Skip due to it is not a html", uri)
		return nil
	}

	return resp.Body
}

func isHtmlType(contentType string) bool {
	contentType = strings.TrimSpace(contentType)
	return strings.HasPrefix(contentType, "text/html")
}
