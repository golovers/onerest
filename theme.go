package onerest

import "errors"

type Theme struct {
	names  []string
	scope *Scope
	OneScopeBaseService
}

func (theme *Theme) Self(selects ...string) ([]Asset, error) {
	themes, err := theme.find(NewQueryBuilder("Theme").Select(selects...).And("Scope.Name", theme.scopeName).And("Name", theme.names...))
	if err != nil {
		return []Asset{}, err
	}
	if len(themes) > 0 {
		return themes, nil
	}
	return []Asset{}, errors.New("Theme not found")
}

func (theme *Theme) Trend(params map[string]string) (Trend, error) {
	result := Trend{}
	selfs, err := theme.Self()
	if err != nil {
		return Trend{}, err
	}
	for _, t := range selfs {
		project, err := theme.scope.Self("ID")
		if err != nil {
			return Trend{}, err
		}
		params["project"] = project[0].Id
		params["theme"] = t.Id
		r, err := theme.OneScopeBaseService.Trend(params)
		if err != nil {
			return result, err
		}
		result.Merge(r)
	}
	return result, nil
}
