package ArgReader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getURL(t *testing.T) {
	t.Run("should getURI give uri without error when arguments are sufficient", func(t *testing.T) {
		args := []string{"goprogram", "url1"}
		uri, err := getURI(args)
		assert.Nil(t, err, "Got Error")
		assert.Equal(t, "http://url1", uri.String())
	})

	t.Run("should getURI give uri without error when arguments are more then sufficient", func(t *testing.T) {
		args := []string{"goprogram", "url1", "somethingelse"}
		uri, err := getURI(args)
		assert.Nil(t, err, "Got Error")
		assert.Equal(t, "http://url1", uri.String())
	})

	t.Run("should getURI raise error when arguments are not sufficient", func(t *testing.T) {
		args := []string{"goprogram"}
		_, err := getURI(args)
		assert.Errorf(t, err, "Need error when no uri argument provided")
	})

	t.Run("should getURI raise error when argument url is wrong", func(t *testing.T) {
		args := []string{"goprogram", "htp/\\fdsa-wer%"}
		_, err := getURI(args)
		assert.Errorf(t, err, "Need error when no uri argument provided")
	})
}
