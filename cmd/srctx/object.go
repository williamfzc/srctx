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

type FuncReferenceScope struct {
	TotalFuncRefCount     int `json:"totalFuncRefCount"`
	CrossFuncFileRefCount int `json:"crossFuncFileRefCount"`
	CrossFuncDirRefCount  int `json:"crossFuncDirRefCount"`
}

func (r *ReferenceScope) IsSafe() bool {
	return r.TotalRefCount == 0 && r.CrossFileRefCount == 0 && r.CrossDirRefCount == 0
}

func (fr *FuncReferenceScope) IsSafe() bool {
	return fr.TotalFuncRefCount == 0 && fr.CrossFuncFileRefCount == 0 && fr.CrossFuncDirRefCount == 0
}

type LineStat struct {
	*FileScope
	RefScope     *ReferenceScope     `json:"ref"`
	FuncRefScope *FuncReferenceScope `json:"funcRef"`
}

func (ls *LineStat) IsSafe() bool {
	return ls.RefScope.IsSafe() && ls.FuncRefScope.IsSafe()
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
		FuncRefScope: &FuncReferenceScope{
			TotalFuncRefCount:     0,
			CrossFuncFileRefCount: 0,
			CrossFuncDirRefCount:  0,
		},
	}
}
