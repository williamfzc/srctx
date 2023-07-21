package file

import (
	"github.com/williamfzc/srctx/object"
)

func (fg *Graph) Stat(f *Vertex) *object.ImpactUnit {
	referenceIds := fg.DirectReferenceIds(f)
	referencedIds := fg.DirectReferencedIds(f)

	transitiveReferencedIds := fg.TransitiveReferencedIds(f)
	transitiveReferenceIds := fg.TransitiveReferenceIds(f)

	entries := fg.EntryIds(f)

	impactUnit := object.NewImpactUnit()

	impactUnit.FileName = f.Path
	impactUnit.UnitName = f.Id()

	impactUnit.ImpactCount = len(referencedIds) + len(referenceIds)
	impactUnit.TransImpactCount = len(transitiveReferenceIds) + len(transitiveReferencedIds)

	impactUnit.ImpactEntries = len(entries)

	// details
	impactUnit.Self = f
	impactUnit.ReferenceIds = referenceIds
	impactUnit.ReferencedIds = referencedIds
	impactUnit.TransitiveReferenceIds = transitiveReferenceIds
	impactUnit.TransitiveReferencedIds = transitiveReferencedIds
	impactUnit.Entries = entries

	return impactUnit
}

func (fg *Graph) GlobalStat(points []*Vertex) *object.StatGlobal {
	sg := &object.StatGlobal{
		UnitLevel: object.NodeLevelFile,
	}

	totalUnits := make([]string, 0, len(fg.IdCache))
	for each := range fg.IdCache {
		totalUnits = append(totalUnits, each)
	}
	sg.TotalUnits = totalUnits

	entries := fg.ListEntries()
	totalEntries := make([]string, 0, len(entries))
	for _, each := range entries {
		totalEntries = append(totalUnits, each.Id())
	}
	sg.TotalEntries = totalEntries

	stats := make(map[string]*object.ImpactUnit, 0)
	for _, each := range points {
		eachStat := fg.Stat(each)
		stats[each.Id()] = eachStat
	}
	sg.ImpactUnitsMap = stats

	// direct impact
	directImpactMap := make(map[string]struct{})
	for _, each := range stats {
		for _, eachReferenced := range each.ReferencedIds {
			directImpactMap[eachReferenced] = struct{}{}
		}
		for _, eachReference := range each.ReferenceIds {
			directImpactMap[eachReference] = struct{}{}
		}
	}
	directImpactList := make([]string, 0, len(directImpactMap))
	for each := range directImpactMap {
		directImpactList = append(directImpactList, each)
	}
	sg.ImpactUnits = directImpactList

	// in-direct impact
	indirectImpactMap := make(map[string]struct{})
	for _, each := range stats {
		for _, eachReferenced := range each.TransitiveReferencedIds {
			indirectImpactMap[eachReferenced] = struct{}{}
		}
		for _, eachReference := range each.TransitiveReferenceIds {
			indirectImpactMap[eachReference] = struct{}{}
		}
	}
	indirectImpactList := make([]string, 0, len(indirectImpactMap))
	for each := range indirectImpactMap {
		indirectImpactList = append(indirectImpactList, each)
	}
	sg.TransImpactUnits = indirectImpactList

	// entries
	entriesMap := make(map[string]struct{})
	for _, each := range stats {
		for _, eachEntry := range each.Entries {
			entriesMap[eachEntry] = struct{}{}
		}
	}
	entriesList := make([]string, 0, len(entriesMap))
	for each := range entriesMap {
		entriesList = append(entriesList, each)
	}
	sg.ImpactEntries = entriesList

	return sg
}
