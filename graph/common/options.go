package common

type GraphOptions struct {
	Src      string `json:"src"`
	LsifFile string `json:"lsifFile"`
	ScipFile string `json:"scipFile"`

	// other options (like performance
	GenGolangIndex bool `json:"genGolangIndex"`
	NoEntries      bool `json:"noEntries"`
}

func DefaultGraphOptions() *GraphOptions {
	return &GraphOptions{
		Src:            ".",
		GenGolangIndex: false,
		LsifFile:       "./dump.lsif",
		ScipFile:       "./index.scip",
		NoEntries:      false,
	}
}
