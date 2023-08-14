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

	impactUnit.ImpactCount = len(referencedIds)
	impactUnit.TransImpactCount = len(transitiveReferencedIds)
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
		UnitLevel:   object.NodeLevelFile,
		UnitMapping: make(map[string]int),
	}

	// creating mapping
	curId := 0
	for each := range fg.IdCache {
		sg.UnitMapping[each] = curId
		curId++
	}

	entries := fg.ListEntries()
	totalEntries := make([]int, 0, len(entries))
	for _, each := range entries {
		eachId := sg.UnitMapping[each.Id()]
		totalEntries = append(totalEntries, eachId)
	}
	sg.TotalEntries = totalEntries

	stats := make(map[int]*object.ImpactUnit, 0)
	for _, each := range points {
		eachStat := fg.Stat(each)
		eachId := sg.UnitMapping[each.Id()]
		stats[eachId] = eachStat
	}
	sg.ImpactUnitsMap = stats

	// direct impact
	directImpactMap := make(map[int]struct{})
	for _, each := range stats {
		for _, eachReferenced := range each.ReferencedIds {
			eachId := sg.UnitMapping[eachReferenced]
			directImpactMap[eachId] = struct{}{}
		}
		// and itself
		itselfId := sg.UnitMapping[each.Self.Id()]
		directImpactMap[itselfId] = struct{}{}
	}
	directImpactList := make([]int, 0, len(directImpactMap))
	for each := range directImpactMap {
		directImpactList = append(directImpactList, each)
	}
	sg.ImpactUnits = directImpactList

	// in-direct impact
	indirectImpactMap := make(map[int]struct{})
	for _, each := range stats {
		for _, eachReferenced := range each.TransitiveReferencedIds {
			eachId := sg.UnitMapping[eachReferenced]
			indirectImpactMap[eachId] = struct{}{}
		}
		// and itself
		itselfId := sg.UnitMapping[each.Self.Id()]
		indirectImpactMap[itselfId] = struct{}{}
	}
	indirectImpactList := make([]int, 0, len(indirectImpactMap))
	for each := range indirectImpactMap {
		indirectImpactList = append(indirectImpactList, each)
	}
	sg.TransImpactUnits = indirectImpactList

	// entries
	entriesMap := make(map[int]struct{})
	for _, each := range stats {
		for _, eachEntry := range each.Entries {
			eachId := sg.UnitMapping[eachEntry]
			entriesMap[eachId] = struct{}{}
		}
	}
	entriesList := make([]int, 0, len(entriesMap))
	for each := range entriesMap {
		entriesList = append(entriesList, each)
	}
	sg.ImpactEntries = entriesList

	return sg
}
