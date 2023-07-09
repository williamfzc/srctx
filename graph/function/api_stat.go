package function

import "github.com/williamfzc/srctx/object"

func (fg *FuncGraph) Stat(f *FuncVertex) *object.ImpactUnit {
	referenceIds := fg.DirectReferenceIds(f)
	referencedIds := fg.DirectReferencedIds(f)

	transitiveReferencedIds := fg.TransitiveReferencedIds(f)
	transitiveReferenceIds := fg.TransitiveReferenceIds(f)

	impactUnit := object.NewImpactUnit()

	impactUnit.FileName = f.Path
	impactUnit.UnitName = f.Id()

	impactUnit.DirectConnectCount = len(referencedIds) + len(referenceIds)
	impactUnit.InDirectConnectCount = len(transitiveReferenceIds) + len(transitiveReferencedIds)
	impactUnit.TotalUnitCount = len(fg.IdCache)

	impactUnit.AffectedEntries = len(fg.EntryIds(f))
	impactUnit.TotalEntriesCount = len(fg.ListEntries())

	// details
	impactUnit.ReferenceIds = referenceIds
	impactUnit.ReferencedIds = referencedIds
	impactUnit.TransitiveReferenceIds = transitiveReferenceIds
	impactUnit.TransitiveReferencedIds = transitiveReferencedIds

	return impactUnit
}
