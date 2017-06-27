package onerest

import "errors"

type Epic struct {
	names  []string
	scope *Scope
	OneScopeBaseService
}

func (epic *Epic) Self(selects ...string) (Asset, error) {
	epics, err := epic.find(NewQueryBuilder("Epic").Select(selects...).And("Scope.Name", epic.scopeName).And("Name", epic.names...))
	if err != nil {
		return Asset{}, err
	}
	if len(epics) > 0 {
		return epics[0], nil
	}
	//TODO to be corrected in case multiple names are given
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
