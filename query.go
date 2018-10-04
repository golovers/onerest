package onerest

import (
	"bytes"
	"net/url"
)

var AND = ";"
var OR = "|"

func NewQueryBuilder(typ string) *QueryBuilder {
	return &QueryBuilder{asset: typ}
}

type QueryBuilder struct {
	asset string
}

func (s *QueryBuilder) Select(selects ...string) *Query {
	q := &Query{asset: s.asset, selects: bytes.NewBufferString(""), wheres: bytes.NewBufferString("")}
	for i, s := range selects {
		if i < len(selects) {
			if len(q.selects.String()) > 0 {
				q.selects.WriteString(",")
			}
			q.selects.WriteString(s)

		}
	}
	return q
}

type Query struct {
	asset   string
	selects *bytes.Buffer
	wheres  *bytes.Buffer
}

func NewQuery() *Query {
	return &Query{asset: "", selects: bytes.NewBufferString(""), wheres: bytes.NewBufferString("")}
}

func (q *Query) And(attribute string, value ...string) *Query {
	return q.AndWithOperator(attribute, "=", value...)
}

func (q *Query) AndWithOperator(attribute string, operator string, value ...string) *Query {
	if len(q.wheres.String()) > 0 {
		q.wheres.WriteString(AND)
	}
	q.wheres.WriteString(attribute)
	q.wheres.WriteString(operator)
	for i, v := range value {
		q.wheres.WriteString(url.QueryEscape("'" + v + "'"))
		if i < len(value)-1 {
			q.wheres.WriteString(",")
		}
	}

	return q
}

func (q *Query) Or(attribute string, values ...string) *Query {
	return q.OrWithOperator(attribute, "=", values...)
}

func (q *Query) OrWithOperator(attribute string, operator string, values ...string) *Query {
	if len(q.wheres.String()) > 0 {
		q.wheres.WriteString(OR)
	}
	q.wheres.WriteString(attribute)
	q.wheres.WriteString(operator)
	for i, v := range values {
		q.wheres.WriteString(url.QueryEscape("'" + v + "'"))
		if i < len(values)-1 {
			q.wheres.WriteString(",")
		}
	}

	return q
}

func (q *Query) AndWithQuery(query *Query) *Query {
	if len(query.wheres.String()) > 0 {
		if len(q.wheres.String()) > 0 {
			q.wheres.WriteString(AND)
		}
		q.wheres.WriteString("(")
		q.wheres.WriteString(query.wheres.String())
		q.wheres.WriteString(")")
	}

	return q
}

func (q *Query) AndWithQueryString(query string) *Query {
	if len(query) > 0 {
		if len(q.wheres.String()) > 0 {
			q.wheres.WriteString(AND)
		}

		q.wheres.WriteString("(")
		q.wheres.WriteString(url.QueryEscape(query))
		q.wheres.WriteString(")")
	}

	return q
}

func (q *Query) OrWithQuery(query *Query) *Query {
	if len(query.wheres.String()) > 0 {
		if len(q.wheres.String()) > 0 {
			q.wheres.WriteString(OR)
		}
		q.wheres.WriteString("(")
		q.wheres.WriteString(query.wheres.String())
		q.wheres.WriteString(")")
	}
	return q
}

func (q *Query) OrWithQueryString(query string) *Query {
	if len(query) > 0 {
		if len(q.wheres.String()) > 0 {
			q.wheres.WriteString(OR)
		}
		q.wheres.WriteString("(")
		q.wheres.WriteString(url.QueryEscape(query))
		q.wheres.WriteString(")")
	}
	return q
}

func (q *Query) Build() string {
	result := bytes.NewBufferString("")
	result.WriteString(q.asset)
	result.WriteString("?")
	sel := q.selects.String()
	if len(sel) > 0 {
		result.WriteString("sel=" + sel)
	}
	wheres := q.wheres.String()
	if len(wheres) > 0 {
		if len(sel) > 0 {
			result.WriteString("&")
		}
		result.WriteString("where=" + wheres)
	}
	return result.String()
}
