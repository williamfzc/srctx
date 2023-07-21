package diff

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	log "github.com/sirupsen/logrus"
)

// ImpactLineMap file name -> lines
type ImpactLineMap = map[string][]int

func GitDiff(rootDir string, before string, after string) (ImpactLineMap, error) {
	// about why, I use cmd rather than some libs
	// because go-git 's patch has some bugs ...
	gitDiffCmd := exec.Command("git", "diff", before, after)
	gitDiffCmd.Dir = rootDir
	data, err := gitDiffCmd.CombinedOutput()
	if err != nil {
		log.Errorf("git cmd error: %s", data)
		return nil, err
	}

	affected, err := Unified2Impact(data)
	if err != nil {
		return nil, err
	}
	return affected, nil
}

func PathOffset(repoRoot string, srcRoot string, origin ImpactLineMap) (ImpactLineMap, error) {
	modifiedLineMap := make(map[string][]int)
	for file, lines := range origin {
		afterPath, err := PathOffsetOne(repoRoot, srcRoot, file)
		if err != nil {
			return nil, err
		}
		modifiedLineMap[afterPath] = lines
	}
	return modifiedLineMap, nil
}

func PathOffsetOne(repoRoot string, srcRoot string, target string) (string, error) {
	absFile := filepath.Join(repoRoot, target)
	return filepath.Rel(srcRoot, absFile)
}

func Unified2Impact(patch []byte) (ImpactLineMap, error) {
	parsed, _, err := gitdiff.Parse(bytes.NewReader(patch))
	if err != nil {
		return nil, err
	}

	impactLineMap := make(ImpactLineMap)
	for _, each := range parsed {
		if each.IsBinary || each.IsDelete {
			continue
		}
		impactLineMap[each.NewName] = make([]int, 0)
		fragments := each.TextFragments
		for _, eachF := range fragments {
			left := int(eachF.NewPosition)

			for i, eachLine := range eachF.Lines {
				if eachLine.New() && eachLine.Op == gitdiff.OpAdd {
					impactLineMap[each.NewName] = append(impactLineMap[each.NewName], left+i-1)
				}
			}
		}
	}
	return impactLineMap, nil
}
