package srctx

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williamfzc/srctx/parser"
)

func TestParser(t *testing.T) {
	file, err := os.Open("./dump.lsif.zip")
	assert.Nil(t, err)
	defer file.Close()
	newParser, err := parser.NewParser(context.Background(), file)
	assert.Nil(t, err)
	assert.NotEmpty(t, newParser.Docs)

	// files
	for id, each := range newParser.Docs.Entries {
		log.Printf("%d %s", id, each)
	}
	// defs
}
