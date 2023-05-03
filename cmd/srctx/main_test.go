package main

import "testing"

func TestDiff(t *testing.T) {
	mainFunc([]string{
		"srctx", "diff",
		"--src", "../..",
		"--before", "HEAD~1",
		"--outputDot", "d.dot",
		"--lsif", "../../dump.lsif"})
}
