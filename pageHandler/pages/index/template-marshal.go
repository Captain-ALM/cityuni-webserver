package index

import (
	"html/template"
	"sort"
)

type Marshal struct {
	Data           DataYaml
	Parameters     template.URL
	OrderStartDate int8
	OrderEndDate   int8
	OrderName      int8
	OrderDuration  int8
	Light          bool
	Counter        int
}

func (m *Marshal) GetEntries() (toReturn []EntryYaml) {
	toReturn = m.Data.Entries
	if m.OrderStartDate > 0 {
		sort.Slice(toReturn, func(i, j int) bool {
			return toReturn[i].StartDate.Before(toReturn[j].StartDate.Time)
		})
	}
	if m.OrderStartDate < 0 {
		sort.Slice(toReturn, func(i, j int) bool {
			return toReturn[i].StartDate.After(toReturn[j].StartDate.Time)
		})
	}
	if m.OrderEndDate > 0 {
		sort.Slice(toReturn, func(i, j int) bool {
			return toReturn[i].GetEndTime().Before(toReturn[j].GetEndTime())
		})
	}
	if m.OrderEndDate < 0 {
		sort.Slice(toReturn, func(i, j int) bool {
			return toReturn[i].GetEndTime().After(toReturn[j].GetEndTime())
		})
	}
	if m.OrderName > 0 {
		sort.Slice(toReturn, func(i, j int) bool {
			return toReturn[i].Name < toReturn[j].Name
		})
	}
	if m.OrderName < 0 {
		sort.Slice(toReturn, func(i, j int) bool {
			return toReturn[i].Name > toReturn[j].Name
		})
	}
	if m.OrderDuration > 0 {
		sort.Slice(toReturn, func(i, j int) bool {
			return toReturn[i].GetDuration() < toReturn[j].GetDuration()
		})
	}
	if m.OrderDuration < 0 {
		sort.Slice(toReturn, func(i, j int) bool {
			return toReturn[i].GetDuration() > toReturn[j].GetDuration()
		})
	}
	return toReturn
}

func (m *Marshal) CounterPlusPlus() int {
	toret := m.Counter
	m.Counter++
	return toret
}
