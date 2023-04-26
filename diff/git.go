package diff

import (
	"bytes"
	"os/exec"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	log "github.com/sirupsen/logrus"
)

// file name -> lines
type AffectedLineMap = map[string][]int

func GitDiff(rootDir string, before string, after string) (AffectedLineMap, error) {
	// about why I use cmd rather than some libs
	// because go-git 's patch has some bugs ...
	gitDiffCmd := exec.Command("git", "diff", before, after)
	gitDiffCmd.Dir = rootDir
	data, err := gitDiffCmd.CombinedOutput()
	if err != nil {
		log.Errorf("git cmd error: %s", data)
		return nil, err
	}

	affected, err := Unified2Affected(data)
	if err != nil {
		return nil, err
	}
	return affected, nil
}

func Unified2Affected(patch []byte) (AffectedLineMap, error) {
	parsed, _, err := gitdiff.Parse(bytes.NewReader(patch))
	if err != nil {
		return nil, err
	}

	affectedMap := make(map[string][]int)
	for _, each := range parsed {
		if each.IsBinary || each.IsDelete {
			continue
		}
		affectedMap[each.NewName] = make([]int, 0)
		fragments := each.TextFragments
		for _, eachF := range fragments {
			left := int(eachF.NewPosition)

			for i, eachLine := range eachF.Lines {
				if eachLine.New() && eachLine.Op == gitdiff.OpAdd {
					affectedMap[each.NewName] = append(affectedMap[each.NewName], left+i-1)
				}
			}
		}
	}
	return affectedMap, nil
}
