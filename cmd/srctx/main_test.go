package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williamfzc/srctx/cmd/srctx/diff"
)

func TestDiff(t *testing.T) {
	t.Run("default diff", func(t *testing.T) {
		mainFunc([]string{
			"srctx", "diff",
			"--src", "../..",
			"--repoRoot", "../..",
			"--before", "HEAD~1",
			"--outputDot", "output.dot",
			"--outputCsv", "output.csv",
			"--outputJson", "output.json",
			"--cacheType", "mem",
			"--lsif", "../../dump.lsif",
		})
	})

	t.Run("raw diff", func(t *testing.T) {
		t.Skip("this case did not work in github action")
		mainFunc([]string{
			"srctx", "diff",
			"--src", "../..",
			"--before", "HEAD~1",
			"--outputDot", "output.dot",
			"--outputCsv", "output.csv",
			"--outputJson", "output.json",
			"--withIndex",
		})
	})

	t.Run("file level diff", func(t *testing.T) {
		mainFunc([]string{
			"srctx", "diff",
			"--src", "../..",
			"--before", "HEAD~1",
			"--outputDot", "output.dot",
			"--outputCsv", "output.csv",
			"--outputJson", "output.json",
			"--lsif", "../../dump.lsif",
			"--nodeLevel", "file",
		})
	})

	t.Run("specific language diff", func(t *testing.T) {
		mainFunc([]string{
			"srctx", "diff",
			"--src", "../..",
			"--before", "HEAD~1",
			"--outputDot", "output.dot",
			"--outputCsv", "output.csv",
			"--outputJson", "output.json",
			"--lsif", "../../dump.lsif",
			"--lang", "GOLANG",
		})
	})

	t.Run("no diff", func(t *testing.T) {
		mainFunc([]string{
			"srctx", "diff",
			"--src", "../..",
			"--outputDot", "output.dot",
			"--outputCsv", "output.csv",
			"--outputJson", "output.json",
			"--lsif", "../../dump.lsif",
			"--noDiff",
		})
	})

	t.Run("dump with existed file", func(t *testing.T) {
		mainFunc([]string{
			"srctx", "dump",
			"--src", "../..",
			"--lsif", "../../dump.lsif",
		})
	})
}

func TestRenderHtml(t *testing.T) {
	t.Run("render func html", func(t *testing.T) {
		mainFunc([]string{
			"srctx", "diff",
			"--src", "../..",
			"--before", "HEAD~1",
			"--outputHtml", "output.html",
			"--lsif", "../../dump.lsif",
			"--nodeLevel", "func",
		})
	})

	t.Run("render file html", func(t *testing.T) {
		mainFunc([]string{
			"srctx", "diff",
			"--src", "../..",
			"--before", "HEAD~1",
			"--outputHtml", "output.html",
			"--lsif", "../../dump.lsif",
			"--nodeLevel", "file",
		})
	})
}

func TestDiffCfg(t *testing.T) {
	t.Run("generate default config file", func(t *testing.T) {
		mainFunc([]string{
			"srctx", "diffcfg",
		})
		defer os.Remove(diff.DefaultConfigFile)
		assert.FileExists(t, diff.DefaultConfigFile)
	})
}

func TestDump(t *testing.T) {
	t.Run("dump", func(t *testing.T) {
		mainFunc([]string{"srctx", "dump", "--src", ".."})
	})
}
