package lexer

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	tokens, err := File2Tokens("./api.go")
	assert.Nil(t, err)

	for _, each := range tokens {
		log.Infof("each: %v", each)
	}
}
