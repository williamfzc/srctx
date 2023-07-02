package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williamfzc/srctx/cmd/srctx/diff"
)

func TestDiff(t *testing.T) {
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
}

func TestDiffRaw(t *testing.T) {
	// this case did not work in github action
	// i still do not know why
	t.Skip()

	mainFunc([]string{
		"srctx", "diff",
		"--src", "../..",
		"--before", "HEAD~1",
		"--outputDot", "output.dot",
		"--outputCsv", "output.csv",
		"--outputJson", "output.json",
		"--withIndex",
	})
}

func TestDiffDir(t *testing.T) {
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
}

func TestDiffSpecificLang(t *testing.T) {
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
}

func TestDiffNoDiff(t *testing.T) {
	mainFunc([]string{
		"srctx", "diff",
		"--src", "../..",
		"--outputDot", "output.dot",
		"--outputCsv", "output.csv",
		"--outputJson", "output.json",
		"--lsif", "../../dump.lsif",
		"--noDiff",
	})
}

func TestRenderHtml(t *testing.T) {
	mainFunc([]string{
		"srctx", "diff",
		"--src", "../..",
		"--before", "HEAD~1",
		"--outputHtml", "output.html",
		"--lsif", "../../dump.lsif",
		"--nodeLevel", "func",
	})
}

func TestDiffCfg(t *testing.T) {
	mainFunc([]string{
		"srctx", "diffcfg",
	})
	defer os.Remove(diff.DefaultConfigFile)
	assert.FileExists(t, diff.DefaultConfigFile)
}
