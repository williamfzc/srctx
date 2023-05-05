package main

import "testing"

func TestDiff(t *testing.T) {
	mainFunc([]string{
		"srctx", "diff",
		"--src", "../..",
		"--before", "HEAD~1",
		"--outputDot", "output.dot",
		"--outputCsv", "output.csv",
		"--lsif", "../../dump.lsif"})
}
