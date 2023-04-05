package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	_, err := GitDiff("../", "HEAD~1", "HEAD")
	assert.Nil(t, err)
}
