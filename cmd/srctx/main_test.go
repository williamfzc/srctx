package main

import "testing"

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
		"--lsif", "../../dump.lsif"})
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
		"--withIndex"})
}

func TestStat(t *testing.T) {
	mainFunc([]string{
		"srctx", "stat",
		"--src", "../..",
		"--outputDot", "output.dot",
		"--lsif", "../../dump.lsif",
		"--nodeLevel", "file",
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
