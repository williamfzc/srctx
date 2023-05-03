package main

import "testing"

func TestDiff(t *testing.T) {
	t.Skip()
	mainFunc([]string{
		"srctx", "diff",
		"--src", "../..",
		"--before", "HEAD~5",
		"--outputDot", "d.dot",
		"--lsif", "../../dump.lsif"})
}
