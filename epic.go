package onerest

import "errors"

type Epic struct {
	name  string
	scope *Scope
	OneScopeBaseService
}

func (epic *Epic) Self(selects ...string) (Asset, error) {
	epics, err := epic.find(NewQueryBuilder("Epic").Select(selects...).And("Scope.Name", epic.scopeName).And("Name", epic.name))
	if err != nil {
		return Asset{}, err
	}
	if len(epics) > 0 {
		return epics[0], nil
	}
	return Asset{}, errors.New("Epic not found")
}

func (epic *Epic) Trend(params map[string]string) (Trend, error) {
	self, err := epic.Self()
	if err != nil {
		return Trend{}, err
	}

	s, err := epic.scope.Self()
	if err != nil {
		return Trend{}, err
	}
	params["project"] = s.Id
	params["epic"] = self.Id

	return epic.OneScopeBaseService.Trend(params)
}
