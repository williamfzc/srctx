package object

/*
StatGlobal

designed for offering a global view
*/
type StatGlobal struct {
	UnitLevel   string         `json:"unitLevel,omitempty"`
	UnitMapping map[string]int `json:"unitMapping,omitempty"`

	ImpactUnits      []int `json:"impactUnits,omitempty"`
	TransImpactUnits []int `json:"transImpactUnits,omitempty"`

	TotalEntries  []int `json:"entries,omitempty"`
	ImpactEntries []int `json:"impactEntries,omitempty"`

	ImpactUnitsMap map[int]*ImpactUnit `json:"-"`
}
