package onerest

type OneRestBuilder struct {
}

func (one *OneRestBuilder) Host(host string) *OneAuthBuilder {
	return &OneAuthBuilder{Host: host}
}

type OneAuthBuilder struct {
	Host string
}

func (auth *OneAuthBuilder) WithUserPassword(user, pass string) IOneService {
	return &OneService{OneBaseService{BaseConstrains: NewQuery(), OneRest: OneRest{Host: auth.Host, Username: user, Password: pass}}}
}

func (auth *OneAuthBuilder) WithAccessToken(token string) IOneService {
	return &OneService{OneBaseService{BaseConstrains: NewQuery(), OneRest: OneRest{Host: auth.Host, AccessToken: token}}}
}

///////////////////////////////////////////////////////////////////////////////////
type IOneBaseService interface {
	Stories(selects ...string) ([]Asset, error)
	Defects(selects ...string) ([]Asset, error)
	Find(typ string, selects ...string) ([]Asset, error)
	FindWith(query *Query) ([]Asset, error)
	Aggregations(typ string, aggregationAttribute string, groupAttributes ...string) (map[string]float64, error)
	Trend(params map[string]string) (Trend, error)
	BaseConstraints() *Query
}

///////////////////////////////////////////////////////////////////////////////////
type OneBaseService struct {
	BaseConstrains *Query
	OneRest
}

func (one *OneBaseService) Stories(selects ...string) ([]Asset, error) {
	return one.find(NewQueryBuilder("Story").Select(selects...).AndWithQuery(one.BaseConstraints()))
}

func (one *OneBaseService) Defects(selects ...string) ([]Asset, error) {
	return one.find(NewQueryBuilder("Defect").Select(selects...).AndWithQuery(one.BaseConstraints()))
}

func (one *OneBaseService) Scopes(selects ...string) ([]Asset, error) {
	return one.find(NewQueryBuilder("Scope").Select(selects...).AndWithQuery(one.BaseConstraints()))
}

func (one *OneBaseService) Find(typ string, selects ...string) ([]Asset, error) {
	return one.find(NewQueryBuilder(typ).Select(selects...).AndWithQuery(one.BaseConstraints()))
}

func (one *OneBaseService) FindWith(query *Query) ([]Asset, error) {
	return one.find(query.AndWithQuery(one.BaseConstraints()))
}

func (one *OneBaseService) BaseConstraints() *Query {
	return one.BaseConstrains
}

func (one *OneBaseService) Aggregations(typ string, aggregationAttribute string, groupAttributes ...string) (map[string]float64, error) {
	g := OneGroup{}
	var sel []string
	sel = append(sel, groupAttributes...)
	if aggregationAttribute != "count" {
		sel = append(sel, aggregationAttribute)
	}
	r, err := one.Find(typ, sel...)
	if err != nil {
		empty := make(map[string]float64)
		return empty, err
	}
	g.AddAssets(r)
	return g.Aggregations(aggregationAttribute, groupAttributes...), nil
}

func (one *OneBaseService) Trend(params map[string]string) (Trend, error) {
	t, errorOne := one.trend(params)
	return t, errorOne
}

///////////////////////////////////////////////////////////////////////////////////
type IOneService interface {
	IOneBaseService
	Scopes(selects ...string) ([]Asset, error)
	Scope(name string) IOneScopeService
}

type OneService struct {
	OneBaseService
}

func (one *OneService) Scopes(selects ...string) ([]Asset, error) {
	return one.Find("Scope", selects...)
}

func (one *OneService) Scope(name string) IOneScopeService {
	one.OneBaseService.BaseConstrains = NewQuery().And("Scope.Name", name)
	return &Scope{OneScopeBaseService{scopeName: name, OneBaseService: one.OneBaseService}}
}

///////////////////////////////////////////////////////////////////////////////////
type IOneScopeBaseService interface {
	IOneBaseService
	Self(selects ...string) ([]Asset, error)
	CreatedInTimeRange(typ string, startInUTC string, endInUTC string, selects ...string) ([]Asset, error)
	CreatedInTimebox(typ string, timebox Asset, selects ...string) ([]Asset, error)
	ScopeName() string
}

type OneScopeBaseService struct {
	scopeName string
	OneBaseService
}

func (one *OneScopeBaseService) Self(selects ...string) ([]Asset, error) {
	return []Asset{}, nil
}

func (one *OneScopeBaseService) ScopeName() string {
	return one.scopeName
}

// List all the assets that created in the given time range
// The given date must be in UTC: 2017-02-17T03:53:26.367
func (one *OneScopeBaseService) CreatedInTimeRange(typ string, startInUTC string, endInUTC string, selects ...string) ([]Asset, error) {
	return one.FindWith(NewQueryBuilder(typ).Select(selects...).AndWithOperator("CreateDateUTC", ">", startInUTC).AndWithOperator("CreateDateUTC", "<", endInUTC))
}

// List all the assets that created in the given timebox
func (one *OneScopeBaseService) CreatedInTimebox(typ string, timebox Asset, selects ...string) ([]Asset, error) {
	return one.CreatedInTimeRange(typ, timebox.GetAsStringValue("BeginDate")+"T00:00:00.000", timebox.GetAsStringValue("EndDate")+"T00:00:00.000", selects...)
}

func (one *OneScopeBaseService) Clone(constraints *Query) *OneScopeBaseService {
	constr := constraints.AndWithQuery(one.BaseConstrains)

	base := OneScopeBaseService{scopeName: one.scopeName,
		OneBaseService: OneBaseService{constr, one.OneRest}}

	return &base
}

///////////////////////////////////////////////////////////////////////////////////
type IOneScopeService interface {
	IOneScopeBaseService
	Themes(selects ...string) ([]Asset, error)
	Epics(selects ...string) ([]Asset, error)
	Timeboxes(selects ...string) ([]Asset, error)
	Timebox(name string, selects ...string) (Asset, error)
	Theme(names ...string) IOneScopeBaseService
	Epic(names ...string) IOneScopeBaseService
	Group(typ string, names ...string) IOneScopeBaseService
}
