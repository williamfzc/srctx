package object

/*
StatGlobal

designed for offering a global view
todo: current struct is good for human reading but bad for data processing (too much storage cost
*/
type StatGlobal struct {
	UnitLevel string `json:"unitLevel,omitempty"`

	TotalUnits       []string `json:"totalUnits,omitempty"`
	ImpactUnits      []string `json:"impactUnits,omitempty"`
	TransImpactUnits []string `json:"transImpactUnits,omitempty"`

	TotalEntries  []string `json:"entries,omitempty"`
	ImpactEntries []string `json:"impactEntries,omitempty"`

	ImpactUnitsMap map[string]*ImpactUnit `json:"-"`
}
