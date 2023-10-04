package common

import "github.com/opensibyl/sibyl2/pkg/core"

type GraphOptions struct {
	Src      string        `json:"src"`
	LsifFile string        `json:"lsifFile"`
	ScipFile string        `json:"scipFile"`
	Lang     core.LangType `json:"lang"`

	// other options (like performance
	GenGolangIndex bool `json:"genGolangIndex"`
	NoEntries      bool `json:"noEntries"`
}

func DefaultGraphOptions() *GraphOptions {
	return &GraphOptions{
		Src:            ".",
		LsifFile:       "./dump.lsif",
		ScipFile:       "./index.scip",
		Lang:           core.LangUnknown,
		NoEntries:      false,
		GenGolangIndex: false,
	}
}
