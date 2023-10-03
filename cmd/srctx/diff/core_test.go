package diff

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCollectLineMap(t *testing.T) {
	impactLineMap, err := collectLineMap(&Options{
		Src:    ".",
		Before: "HEAD~1",
		After:  "HEAD",
	})
	assert.Nil(t, err)
	log.Debugf("map: %v", impactLineMap)
}
