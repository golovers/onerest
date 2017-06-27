package onerest

import (
	"sort"
	"strings"
	"time"
	"errors"
)

type Scope struct {
	OneScopeBaseService
}

func (s *Scope) Theme(names ...string) IOneScopeBaseService {
	base := s.OneScopeBaseService.Clone(NewQuery().And("ParentAndUp.Name", names...))
	return &Theme{names: names, scope: s, OneScopeBaseService: *base}
}

func (s *Scope) Epic(names ...string) IOneScopeBaseService {
	base := s.OneScopeBaseService.Clone(NewQuery().And("SuperAndUp.Name", names...))
	return &Epic{names: names, scope: s, OneScopeBaseService: *base}
}

func (s *Scope) Group(typ string, names ...string) IOneScopeBaseService {
	if strings.ToLower(typ) == "epic" {
		return s.Epic(names...)
	} else {
		return s.Theme(names...)
	}
}

func (s *Scope) Self(selects ...string) ([]Asset, error) {
	scopes, err := s.find(NewQueryBuilder("Scope").Select(selects...).
		And("Name", s.scopeName))
	if err != nil {
		return  []Asset{}, err
	}
	if len(scopes) > 0 {
		return scopes, nil
	}
	return []Asset{}, errors.New("Scope not found")
}

func (s *Scope) Themes(selects ...string) ([]Asset, error) {
	return s.FindWith(NewQueryBuilder("Theme").Select(selects...))
}

func (s *Scope) Epics(selects ...string) ([]Asset, error) {
	return s.FindWith(NewQueryBuilder("Epic").Select(selects...))
}

func (s *Scope) Timeboxes(selects ...string) ([]Asset, error) {
	sel := selects
	sel = append(selects, "Name")

	scp, err := s.Self("Name", "Schedule.Name")
	if err != nil {
		return []Asset{}, err
	}
	rs, err := s.find(NewQueryBuilder("Timebox").Select(sel...).
		And("Schedule.Name", scp[0].GetAsStringValue("Schedule.Name")))
	if err != nil {
		return []Asset{}, err
	}
	// Sort result
	sort.Slice(rs, func(i, j int) bool {
		if strings.Compare(rs[i].GetAsStringValue("Name"),
			rs[j].GetAsStringValue("Name")) == 1 {
			return false
		}
		return true
	})
	return rs, nil
}

func (s *Scope) Timebox(name string, selects ...string) (Asset, error) {
	sel := selects
	sel = append(sel, "Name")
	sel = append(sel, "BeginDate")
	sel = append(sel, "EndDate")

	timeboxes, err := s.Timeboxes(sel...)
	if err != nil {
		return Asset{}, err
	}
	now := time.Now().UTC()
	for _, tbx := range timeboxes {
		if name == "current" {
			begin, _ := time.Parse("2006-1-2", tbx.GetAsStringValue("BeginDate"))
			end, _ := time.Parse("2006-1-2", tbx.GetAsStringValue("EndDate"))

			if now.After(begin.UTC()) && now.Before(end.UTC()) {
				return tbx, nil
			}
		} else {
			if tbx.GetAsStringValue("Name") == name {
				return tbx, nil
			}
		}
	}
	return Asset{},errors.New("Timebox not found")
}

func (s *Scope) Trend(params map[string]string) (Trend, error) {
	a, err := s.Self()
	if err != nil {
		return Trend{}, err
	}
	params["project"] = a[0].Id
	return s.OneScopeBaseService.Trend(params)
}
