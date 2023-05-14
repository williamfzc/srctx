package main

import "testing"

func TestDiff(t *testing.T) {
	mainFunc([]string{
		"srctx", "diff",
		"--src", "../..",
		"--before", "HEAD~1",
		"--outputDot", "output.dot",
		"--outputCsv", "output.csv",
		"--outputJson", "output.json",
		"--lsif", "../../dump.lsif"})
}

func TestDiffRaw(t *testing.T) {
	mainFunc([]string{
		"srctx", "diff",
		"--src", "../..",
		"--before", "HEAD~1",
		"--outputDot", "output.dot",
		"--outputCsv", "output.csv",
		"--outputJson", "output.json",
		"--withIndex"})
}
