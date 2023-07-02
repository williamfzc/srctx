package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const (
	nodeLevelFunc = "func"
	nodeLevelFile = "file"
	nodeLevelDir  = "dir"

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

	// config file
	DefaultConfigFile = "srctx_cfg.json"
)

type Options struct {
	// required
	Src      string `json:"src,omitempty"`
	RepoRoot string `json:"repoRoot,omitempty"`
	Before   string `json:"before,omitempty"`
	After    string `json:"after,omitempty"`
	LsifZip  string `json:"lsifZip,omitempty"`
	ScipFile string `json:"scipFile,omitempty"`

	// output
	OutputJson string `json:"outputJson,omitempty"`
	OutputCsv  string `json:"outputCsv,omitempty"`
	OutputDot  string `json:"outputDot,omitempty"`
	OutputHtml string `json:"outputHtml,omitempty"`

	// options
	NodeLevel string `json:"nodeLevel,omitempty"`
	WithIndex bool   `json:"withIndex,omitempty"`
	CacheType string `json:"cacheType,omitempty"`
	Lang      string `json:"lang,omitempty"`
	NoDiff    bool   `json:"noDiff,omitempty"`
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

type FileVertex struct {
	FileName                 string  `csv:"fileName" json:"fileName"`
	AffectedLinePercent      float32 `csv:"affectedLinePercent" json:"affectedLinePercent"`
	AffectedFunctionPercent  float32 `csv:"affectedFunctionPercent" json:"affectedFunctionPercent"`
	AffectedReferencePercent float32 `csv:"affectedReferencePercent" json:"affectedReferencePercent"`

	AffectedLines int `csv:"affectedLines" json:"affectedLines"`
	TotalLines    int `csv:"totalLines" json:"totalLines"`

	AffectedFunctions int `csv:"affectedFunctions" json:"affectedFunctions"`
	TotalFunctions    int `csv:"totalFunctions" json:"totalFunctions"`

	AffectedReferences   int      `csv:"affectedReferences" json:"affectedReferences"`
	AffectedReferenceIds []string `csv:"-" json:"-"`
	TotalReferences      int      `csv:"totalReferences" json:"totalReferences"`
}
