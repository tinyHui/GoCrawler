package CrawlerEngine

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_underSameDomain(t *testing.T) {
	t.Run("should return true when given two urls under same baseDomain", func(t *testing.T) {
		engine := &crawlerEngine{
			baseDomain: "abc.url",
		}

		uri, _ := url.Parse("http://blog.abc.url")
		assert.True(t, engine.underSameDomain(*uri))

		uri, _ = url.Parse("http://www.abc.url:1234")
		assert.True(t, engine.underSameDomain(*uri))

		uri, _ = url.Parse("http://abc.url")
		assert.True(t, engine.underSameDomain(*uri))
	})

	t.Run("should return true when given two urls under same baseDomain but different schema", func(t *testing.T) {
		engine := &crawlerEngine{
			baseDomain: "abc.url",
		}

		uri2, _ := url.Parse("https://blog.abc.url")
		assert.True(t, engine.underSameDomain(*uri2))
	})
}

func Test_getDomain(t *testing.T) {
	t.Run("should get baseDomain", func(t *testing.T) {
		assert.Equal(t, "monzo.com", getDomain(url.URL{
			Host: "www.monzo.com",
		}))

		assert.Equal(t, "monzo.com", getDomain(url.URL{
			Host: "monzo.com",
		}))

		assert.Equal(t, "monzo.com", getDomain(url.URL{
			Host: "blog.monzo.com",
		}))
	})

	t.Run("should give empty baseDomain when host is empty", func(t *testing.T) {
		assert.Equal(t, "", getDomain(url.URL{
			Host: "",
		}))
	})
}
