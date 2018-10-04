package onerest

import "errors"

type Epic struct {
	names []string
	scope *Scope
	OneScopeBaseService
}

func (epic *Epic) Self(selects ...string) ([]Asset, error) {
	epics, err := epic.find(NewQueryBuilder("Epic").Select(selects...).And("Scope.Name", epic.scopeName).And("Name", epic.names...))
	if err != nil {
		return []Asset{}, err
	}
	if len(epics) > 0 {
		return epics, nil
	}
	return []Asset{}, errors.New("epic not found")
}

func (epic *Epic) Trend(params map[string]string) (Trend, error) {
	result := Trend{}
	selfs, err := epic.Self()
	if err != nil {
		return Trend{}, err
	}
	for _, e := range selfs {
		s, err := epic.scope.Self()
		if err != nil {
			return Trend{}, err
		}
		params["project"] = s[0].Id
		params["epic"] = e.Id
		r, err := epic.OneScopeBaseService.Trend(params)
		if err != nil {
			return result, err

		}
		result.Merge(r)
	}
	return result, nil
}
