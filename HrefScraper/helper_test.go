package HrefScraper

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidURI(t *testing.T) {
	t.Run("should return true when given http link", func(t *testing.T) {
		uri, _ := url.Parse("http://abc.url")
		assert.True(t, IsValidURI(*uri))
	})

	t.Run("should return true when given https link", func(t *testing.T) {
		uri, _ := url.Parse("https://abc.url")
		assert.True(t, IsValidURI(*uri))
	})

	t.Run("should return false when given no specific schema link", func(t *testing.T) {
		uri := url.URL{
			Scheme: "",
			Host:   "www.google.com",
		}
		assert.True(t, IsValidURI(uri))
	})

	t.Run("should return false when given other schema link", func(t *testing.T) {
		uri, _ := url.Parse("mailto://someone")
		assert.False(t, IsValidURI(*uri))
	})

	t.Run("should return false when link is empty", func(t *testing.T) {
		uri, _ := url.Parse("")
		assert.False(t, IsValidURI(*uri))
	})
}

func Test_strConvURL(t *testing.T) {
	t.Run("should convert string to url", func(t *testing.T) {
		uri, _ := StrConvURL("http://monzo.com")
		assert.Equal(t, "monzo.com", uri.Host)

		uri, _ = StrConvURL("http://www.monzo.com")
		assert.Equal(t, "www.monzo.com", uri.Host)

		uri, _ = StrConvURL("http://www.monzo.com/path")
		assert.Equal(t, "www.monzo.com", uri.Host)
		assert.Equal(t, "/path", uri.Path)
	})

	t.Run("should convert string to url when no schema defined", func(t *testing.T) {
		uri, _ := StrConvURL("monzo.com")
		assert.Equal(t, "monzo.com", uri.Host)
		assert.Equal(t, "", "")
		assert.Equal(t, "http", uri.Scheme)

		uri, _ = StrConvURL("www.monzo.com")
		assert.Equal(t, "www.monzo.com", uri.Host)
		assert.Equal(t, "", "")
		assert.Equal(t, "http", uri.Scheme)

		uri, _ = StrConvURL("www.monzo.com/path")
		assert.Equal(t, "www.monzo.com", uri.Host)
		assert.Equal(t, "/path", uri.Path)
		assert.Equal(t, "http", uri.Scheme)
	})
}

func Test_trimHash(t *testing.T) {
	t.Run("Should trim hash from url", func(t *testing.T) {
		assert.Equal(t, "http://any.url", trimHash("http://any.url#tag"))
		assert.Equal(t, "/any.url", trimHash("/any.url#tag_"))
		assert.Equal(t, "http://any.url", trimHash("http://any.url#tag_#23"))
		assert.Equal(t, "http://any.url", trimHash("http://any.url#tag_#!@#@!"))
		assert.Equal(t, "http://any.url/1?a=2", trimHash("http://any.url/1?a=2"))
		assert.Equal(t, "/1?a=2", trimHash("/1?a=2#tag"))
	})
}

func Test_complete(t *testing.T) {
	t.Run("Should complete with base url when given link gives relative path", func(t *testing.T) {
		baseURI, _ := url.Parse("http://any.url/first/second")

		uri, _ := complete(*baseURI, "http://monzo.com")
		assert.Equal(t, "http://monzo.com", uri.String())
		uri, _ = complete(*baseURI, "/abc")
		assert.Equal(t, "http://any.url/abc", uri.String())
		uri, _ = complete(*baseURI, "../fgh")
		assert.Equal(t, "http://any.url/first/fgh", uri.String())
		uri, _ = complete(*baseURI, "../../lmn")
		assert.Equal(t, "http://any.url/lmn", uri.String())
		uri, _ = complete(*baseURI, "./abc")
		assert.Equal(t, "http://any.url/first/second/abc", uri.String())
		uri, _ = complete(*baseURI, "/abc?a=123")
		assert.Equal(t, "http://any.url/abc%3Fa=123", uri.String())
		uri, _ = complete(*baseURI, "/abc?")
		assert.Equal(t, "http://any.url/abc%3F", uri.String())
		uri, _ = complete(*baseURI, "OData.svc/Products-meta?$filter=Name gt 'Milk'")
		assert.Equal(t, "http://any.url/first/second/OData.svc/Products-meta%3F$filter=Name%20gt%20%27Milk%27", uri.String())
	})

	t.Run("Should give new path if given new url", func(t *testing.T) {
		baseURI, _ := url.Parse("http://any.url/first/second")

		uri, _ := complete(*baseURI, "http://another.url")
		assert.Equal(t, "http://another.url", uri.String())
		uri, _ = complete(*baseURI, "https://another.https.url")
		assert.Equal(t, "https://another.https.url", uri.String())
		uri, _ = complete(*baseURI, "mailto://someone")
		assert.Equal(t, "mailto://someone", uri.String())
	})
}

func Test_deduplicate(t *testing.T) {
	t.Run("Should deduplicate urls", func(t *testing.T) {
		source := []url.URL{
			{Host: "url1"}, {Host: "url2"}, {Host: "url3"}, {Host: "url3"}, {Host: "url2"}, {Host: "url1"},
		}

		deduplicated := deduplicate(source)

		assert.Equal(t, 3, len(deduplicated))
		assert.Equal(t, "url1", deduplicated[0].Host)
		assert.Equal(t, "url2", deduplicated[1].Host)
		assert.Equal(t, "url3", deduplicated[2].Host)
	})
}
