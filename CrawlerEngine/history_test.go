package CrawlerEngine

import (
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_visited(t *testing.T) {
	t.Run("should give false when url is not visited", func(t *testing.T) {
		engine := &crawlerEngine{
			historyRecord: &sync.Map{},
		}

		uri, _ := url.Parse("http://monzo.com")
		assert.False(t, engine.visited(*uri))
	})

	t.Run("should give true when url is visited", func(t *testing.T) {
		engine := &crawlerEngine{
			historyRecord: &sync.Map{},
		}

		uri, _ := url.Parse("http://monzo.com")
		engine.historyRecord.Store(uri.String(), true)

		uriSame, _ := url.Parse("http://monzo.com")
		assert.True(t, engine.visited(*uriSame))
	})
}

func Test_markVisited(t *testing.T) {
	t.Run("should uri mark as visited", func(t *testing.T) {
		engine := &crawlerEngine{
			historyRecord: &sync.Map{},
		}

		uri, _ := url.Parse("http://monzo.com")
		v, _ := engine.historyRecord.Load(uri.String())
		assert.Nil(t, v)

		uriSame, _ := url.Parse("http://monzo.com")
		engine.markVisited(*uriSame)
		v, _ = engine.historyRecord.Load(uri.String())
		assert.True(t, v.(bool))
	})
}
