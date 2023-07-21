package object

type FileInfoPart struct {
	FileName string `csv:"fileName" json:"fileName"`

	// actually graph will not access the real file system
	// so of course it knows nothing about the real files
	// all the data we can access is from the indexing file
}

type UnitImpactPart struct {
	// unit level
	// if file level, UnitName == FileName
	// if func level, UnitName == FuncSignature
	// ...
	UnitName string `csv:"unitName" json:"unitName"`

	// Heat
	ImpactCount      int `csv:"impactCount" json:"impactCount"`
	TransImpactCount int `csv:"transImpactCount" json:"transImpactCount"`
	TotalUnitCount   int `csv:"totalUnitCount" json:"totalUnitCount"`

	// entries
	ImpactEntries     int `csv:"impactEntries" json:"impactEntries"`
	TotalEntriesCount int `csv:"totalEntriesCount" json:"totalEntriesCount"`
}

type ImpactDetails struct {
	// raw
	ReferencedIds           []string `json:"-" csv:"-"`
	ReferenceIds            []string `json:"-" csv:"-"`
	TransitiveReferencedIds []string `json:"-" csv:"-"`
	TransitiveReferenceIds  []string `json:"-" csv:"-"`
}

type ImpactUnit struct {
	*FileInfoPart
	*UnitImpactPart
	*ImpactDetails `json:"-" csv:"-"`
}

func NewImpactUnit() *ImpactUnit {
	return &ImpactUnit{
		FileInfoPart: &FileInfoPart{
			FileName: "",
		},
		UnitImpactPart: &UnitImpactPart{
			UnitName:          "",
			ImpactCount:       0,
			TransImpactCount:  0,
			TotalUnitCount:    0,
			ImpactEntries:     0,
			TotalEntriesCount: 0,
		},
		ImpactDetails: &ImpactDetails{
			ReferencedIds:           make([]string, 0),
			ReferenceIds:            make([]string, 0),
			TransitiveReferencedIds: make([]string, 0),
			TransitiveReferenceIds:  make([]string, 0),
		},
	}
}

func (iu *ImpactUnit) VisitedIds() []string {
	return append(iu.TransitiveReferenceIds, iu.TransitiveReferencedIds...)
}
