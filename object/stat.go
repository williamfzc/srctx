package object

type FileImpactPart struct {
	// file level
	FileName string `csv:"fileName" json:"fileName"`

	// line level impact
	AffectedLineCount int `csv:"affectedLineCount" json:"affectedLineCount"`
	TotalLineCount    int `csv:"totalLineCount" json:"totalLineCount"`
}

type UnitImpactPart struct {
	// unit level
	// if file level, UnitName == FileName
	// if func level, UnitName == FuncSignature
	// ...
	UnitName string `csv:"unitName" json:"unitName"`

	// Heat
	DirectConnectCount   int `csv:"directConnectCount" json:"directConnectCount"`
	InDirectConnectCount int `csv:"inDirectConnectCount" json:"inDirectConnectCount"`
	TotalUnitCount       int `csv:"totalUnitCount" json:"totalUnitCount"`

	// entries
	AffectedEntries   int `csv:"affectedEntries" json:"affectedEntries"`
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
	*FileImpactPart
	*UnitImpactPart
	*ImpactDetails `json:"-" csv:"-"`
}

func NewImpactUnit() *ImpactUnit {
	return &ImpactUnit{
		FileImpactPart: &FileImpactPart{
			FileName:          "",
			AffectedLineCount: 0,
			TotalLineCount:    0,
		},
		UnitImpactPart: &UnitImpactPart{
			UnitName:             "",
			DirectConnectCount:   0,
			InDirectConnectCount: 0,
			TotalUnitCount:       0,
			AffectedEntries:      0,
			TotalEntriesCount:    0,
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
