package function

type VertexStat struct {
	Referenced           int `json:"referenced" csv:"referenced"`
	Reference            int `json:"reference" csv:"reference"`
	TransitiveReferenced int `json:"transitiveReferenced" csv:"transitiveReferenced"`
	TransitiveReference  int `json:"transitiveReference" csv:"transitiveReference"`

	// raw
	Root                    *FuncVertex `json:"-" csv:"-"`
	ReferencedIds           []string    `json:"-" csv:"-"`
	ReferenceIds            []string    `json:"-" csv:"-"`
	TransitiveReferencedIds []string    `json:"-" csv:"-"`
	TransitiveReferenceIds  []string    `json:"-" csv:"-"`
}

func (v *VertexStat) VisitedIds() []string {
	return append(v.TransitiveReferenceIds, v.TransitiveReferencedIds...)
}

func (fg *FuncGraph) Stat(f *FuncVertex) *VertexStat {
	referenceIds := fg.ReferenceIds(f)
	referencedIds := fg.ReferencedIds(f)

	transitiveReferencedIds := fg.TransitiveReferencedIds(f)
	transitiveReferenceIds := fg.TransitiveReferenceIds(f)

	return &VertexStat{
		Referenced:           len(referencedIds),
		Reference:            len(referenceIds),
		TransitiveReferenced: len(transitiveReferencedIds),
		TransitiveReference:  len(transitiveReferenceIds),

		Root:                    f,
		ReferencedIds:           referencedIds,
		ReferenceIds:            referenceIds,
		TransitiveReferencedIds: transitiveReferencedIds,
		TransitiveReferenceIds:  transitiveReferenceIds,
	}
}
