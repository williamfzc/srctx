package function

import "github.com/williamfzc/srctx/object"

func (fg *Graph) Stat(f *Vertex) *object.ImpactUnit {
	referenceIds := fg.DirectReferenceIds(f)
	referencedIds := fg.DirectReferencedIds(f)

	transitiveReferencedIds := fg.TransitiveReferencedIds(f)
	transitiveReferenceIds := fg.TransitiveReferenceIds(f)

	impactUnit := object.NewImpactUnit()

	impactUnit.FileName = f.Path
	impactUnit.UnitName = f.Id()

	impactUnit.ImpactCount = len(referencedIds) + len(referenceIds)
	impactUnit.TransImpactCount = len(transitiveReferenceIds) + len(transitiveReferencedIds)
	impactUnit.TotalUnitCount = len(fg.IdCache)

	impactUnit.ImpactEntries = len(fg.EntryIds(f))
	impactUnit.TotalEntriesCount = len(fg.ListEntries())

	// details
	impactUnit.ReferenceIds = referenceIds
	impactUnit.ReferencedIds = referencedIds
	impactUnit.TransitiveReferenceIds = transitiveReferenceIds
	impactUnit.TransitiveReferencedIds = transitiveReferencedIds

	return impactUnit
}

func (fg *Graph) GlobalStat() {
}
