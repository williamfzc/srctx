package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/williamfzc/srctx/object"

	"github.com/urfave/cli/v2"
)

const (
	nodeLevelFunc = "func"
	nodeLevelFile = "file"

	// flags
	srcFlagName        = "src"
	repoRootFlagName   = "repoRoot"
	beforeFlagName     = "before"
	afterFlagName      = "after"
	lsifFlagName       = "lsif"
	scipFlagName       = "scip"
	nodeLevelFlagName  = "nodeLevel"
	outputJsonFlagName = "outputJson"
	outputCsvFlagName  = "outputCsv"
	outputDotFlagName  = "outputDot"
	outputHtmlFlagName = "outputHtml"
	withIndexFlagName  = "withIndex"
	cacheTypeFlagName  = "cacheType"
	langFlagName       = "lang"
	noDiffFlagName     = "noDiff"
	indexCmdFlagName   = "indexCmd"

	// config file
	DefaultConfigFile = "srctx_cfg.json"
)

type Options struct {
	// required
	Src      string `json:"src"`
	RepoRoot string `json:"repoRoot"`
	Before   string `json:"before"`
	After    string `json:"after"`
	LsifZip  string `json:"lsifZip"`
	ScipFile string `json:"scipFile"`

	// output
	OutputJson string `json:"outputJson"`
	OutputCsv  string `json:"outputCsv"`
	OutputDot  string `json:"outputDot"`
	OutputHtml string `json:"outputHtml"`

	// options
	NodeLevel string `json:"nodeLevel"`
	WithIndex bool   `json:"withIndex"`
	CacheType string `json:"cacheType"`
	Lang      string `json:"lang"`
	NoDiff    bool   `json:"noDiff"`
	IndexCmd  string `json:"indexCmd"`
}

func NewOptionsFromCliFlags(c *cli.Context) *Options {
	return &Options{
		Src:        c.String(srcFlagName),
		RepoRoot:   c.String(repoRootFlagName),
		Before:     c.String(beforeFlagName),
		After:      c.String(afterFlagName),
		LsifZip:    c.String(lsifFlagName),
		ScipFile:   c.String(scipFlagName),
		OutputJson: c.String(outputJsonFlagName),
		OutputCsv:  c.String(outputCsvFlagName),
		OutputDot:  c.String(outputDotFlagName),
		OutputHtml: c.String(outputHtmlFlagName),
		NodeLevel:  c.String(nodeLevelFlagName),
		WithIndex:  c.Bool(withIndexFlagName),
		CacheType:  c.String(cacheTypeFlagName),
		Lang:       c.String(langFlagName),
		NoDiff:     c.Bool(noDiffFlagName),
		IndexCmd:   c.String(indexCmdFlagName),
	}
}

func NewOptionsFromSrc(src string) (*Options, error) {
	return NewOptionsFromJSONFile(filepath.Join(src, DefaultConfigFile))
}

func NewOptionsFromJSONFile(fp string) (*Options, error) {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist at %s", fp)
	}

	jsonContent, err := os.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var opts Options
	if err := json.Unmarshal(jsonContent, &opts); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &opts, nil
}

type ImpactUnitWithFile struct {
	*object.ImpactUnit

	// line level impact
	AffectedLineCount int `csv:"affectedLineCount" json:"affectedLineCount"`
	TotalLineCount    int `csv:"totalLineCount" json:"totalLineCount"`
}

func WrapImpactUnitWithFile(impactUnit *object.ImpactUnit) *ImpactUnitWithFile {
	return &ImpactUnitWithFile{
		ImpactUnit:        impactUnit,
		AffectedLineCount: 0,
		TotalLineCount:    0,
	}
}
