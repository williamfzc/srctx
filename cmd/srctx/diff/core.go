package diff

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/diff"
	"github.com/williamfzc/srctx/object"
	"github.com/williamfzc/srctx/parser"
	"github.com/williamfzc/srctx/parser/lsif"
)

// MainDiff allow accessing as a lib
func MainDiff(opts *Options) error {
	log.Infof("start diffing: %v", opts.Src)

	if opts.CacheType != lsif.CacheTypeFile {
		parser.UseMemCache()
	}

	// collect diff info
	lineMap, err := collectLineMap(opts)
	if err != nil {
		return err
	}

	// collect info from file (line number/size ...)
	totalLineCountMap, err := collectTotalLineCountMap(opts, opts.Src, lineMap)
	if err != nil {
		return err
	}

	err = createIndexFile(opts)
	if err != nil {
		return err
	}

	switch opts.NodeLevel {
	case object.NodeLevelFunc:
		err = funcLevelMain(opts, lineMap, totalLineCountMap)
		if err != nil {
			return err
		}
	case object.NodeLevelFile:
		err = fileLevelMain(opts, lineMap, totalLineCountMap)
		if err != nil {
			return err
		}
	}

	log.Infof("everything done.")
	return nil
}

func createIndexFile(opts *Options) error {
	if opts.IndexCmd == "" {
		return nil
	}

	log.Infof("create index file with cmd: %v", opts.IndexCmd)

	parts := strings.Split(opts.IndexCmd, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func collectLineMap(opts *Options) (diff.ImpactLineMap, error) {
	if !opts.NoDiff {
		lineMap, err := diff.GitDiff(opts.Src, opts.Before, opts.After)
		if err != nil {
			return nil, err
		}
		return lineMap, nil
	}
	log.Infof("noDiff enabled")
	return make(diff.ImpactLineMap), nil
}

func collectTotalLineCountMap(opts *Options, src string, lineMap diff.ImpactLineMap) (map[string]int, error) {
	totalLineCountMap := make(map[string]int)
	var err error

	if opts.RepoRoot != "" {
		repoRoot, err := filepath.Abs(opts.RepoRoot)
		if err != nil {
			return nil, err
		}

		log.Infof("path sync from %s to %s", repoRoot, src)
		lineMap, err = diff.PathOffset(repoRoot, src, lineMap)
		if err != nil {
			return nil, err
		}
	}

	for eachPath := range lineMap {
		totalLineCountMap[eachPath], err = lineCounter(filepath.Join(src, eachPath))
		if err != nil {
			// ignore this err, files can be removed/moved
			log.Infof("file has been removed: %s, set line counter to 0", eachPath)
		}
	}

	return totalLineCountMap, nil
}

// https://stackoverflow.com/a/24563853
func lineCounter(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount, nil
}
