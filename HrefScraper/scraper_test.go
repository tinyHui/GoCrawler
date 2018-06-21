package HrefScraper

import (
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type reader int

func (reader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func Test_GrabLinks(t *testing.T) {
	t.Run("Should get all urls from ahref tag when base url is a top level url", func(t *testing.T) {
		reader := strings.NewReader(`<p>
										  <a href="http://any.url">1</a>
										  <a href='http://www.google.com'>2</a>
										  <a style="" href=http://www.google.com/q=123>search 123</a>
										  <p href=https://www.google.com/q=abc>http://www.google.com/q=abc</p>
										  http://dontcollectthis.com
										  <a href="/abc">1</a>
										  <a href="mailto://someone"</a>
										  <a href='/abc?a=1'>2</a>
										  <a href='/abc?b=2&c=3#abc'>2</a>
										  <a href='/abc?b=2&c=3#cde'>3</a>
					                      <a href='#anyhashtag'>hastag</a>
										  <a href='   /cde   '></a>
										  <a href="OData.svc/Products-meta?$filter=Name gt 'Milk'"></a>
										</p>`)

		scraper := NewScraper()

		var baseURI, _ = url.Parse("http://any.base.url")
		links := scraper.GrabLinks(*baseURI, reader)

		assert.Equal(t, 9, len(links))
		assert.Equal(t, "http://any.url", links[0].String())
		assert.Equal(t, "http://www.google.com", links[1].String())
		assert.Equal(t, "http://www.google.com/q=123", links[2].String())
		assert.Equal(t, "http://any.base.url/abc", links[3].String())
		assert.Equal(t, "http://any.base.url/abc%3Fa=1", links[4].String())
		assert.Equal(t, "http://any.base.url/abc%3Fb=2&c=3", links[5].String())
		assert.Equal(t, "http://any.base.url", links[6].String())
		assert.Equal(t, "http://any.base.url/cde", links[7].String())
		assert.Equal(t, "http://any.base.url/OData.svc/Products-meta%3F$filter=Name%20gt%20%27Milk%27", links[8].String())
	})

	t.Run("Should get all urls from ahref tag when base url is not a top level url", func(t *testing.T) {
		baseURI, _ := url.Parse("https://blog.aixc.space/tags/record/")
		reader := strings.NewReader(`<div class="powered-by">
		Powered by <a class="theme-link" href="https://hexo.io">Hexo</a>
		<a href="OData.svc/Products-meta?$filter=Name gt 'Milk'"></a>
		</div>`)
		scraper := NewScraper()

		links := scraper.GrabLinks(*baseURI, reader)
		assert.Equal(t, 2, len(links))
		assert.Equal(t, "https://hexo.io", links[0].String())
		assert.Equal(t, "https://blog.aixc.space/tags/record/OData.svc/Products-meta%3F$filter=Name%20gt%20%27Milk%27", links[1].String())

	})

	t.Run("Should return empty list when reader got error while reading", func(t *testing.T) {
		reader := reader(0)

		scraper := NewScraper()

		var baseURI, _ = url.Parse("http://any.base.url")
		links := scraper.GrabLinks(*baseURI, reader)

		assert.Equal(t, 0, len(links))
	})
}
