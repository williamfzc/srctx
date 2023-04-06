package main

type FileScope struct {
	FileName   string `json:"fileName"`
	LineNumber int    `json:"lineNumber"`
}

type ReferenceScope struct {
	TotalRefCount     int `json:"totalRefCount"`
	CrossFileRefCount int `json:"crossFileRefCount"`
	CrossDirRefCount  int `json:"crossDirRefCount"`
}

type LineStat struct {
	*FileScope
	RefScope *ReferenceScope `json:"ref"`
}

func NewLineStat(fileName string, lineNumber int) *LineStat {
	return &LineStat{
		FileScope: &FileScope{
			FileName:   fileName,
			LineNumber: lineNumber,
		},
		RefScope: &ReferenceScope{
			CrossFileRefCount: 0,
			CrossDirRefCount:  0,
		},
	}
}
