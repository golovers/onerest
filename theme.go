package onerest

import "errors"

type Theme struct {
	name  string
	scope *Scope
	OneScopeBaseService
}

func (theme *Theme) Self(selects ...string) (Asset, error) {
	themes, err := theme.find(NewQueryBuilder("Theme").Select(selects...).And("Scope.Name", theme.scopeName).And("Name", theme.name))
	if err != nil {
		return Asset{}, err
	}
	if len(themes) > 0 {
		return themes[0], nil
	}
	return Asset{}, errors.New("Theme not found")
}

func (theme *Theme) Trend(params map[string]string) (Trend, error) {
	self, err := theme.Self()
	if err != nil {
		return Trend{}, err
	}
	project, err := theme.scope.Self("ID")
	if err != nil {
		return Trend{}, err
	}
	params["project"] = project.Id
	params["theme"] = self.Id

	return theme.OneScopeBaseService.Trend(params)
}
