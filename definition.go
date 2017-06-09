package onerest

import (
	"strings"
	"reflect"
	"fmt"
)

// Reports objects
type Serie struct {
	Date  string
	Value float64
}

type Node struct {
	OidToken string
	Name     string
	Series   []Serie
}

type Trend struct {
	Result [] *Node
}

// Assets objects

type Attribute struct {
	_type string
	Name  string
	Value interface{}
}

type Asset struct {
	Type       string `json:"_type"`
	Href       string
	Id         string
	Attributes map[string]Attribute
}

type OneResponse struct {
	_type     string
	Total     int64
	PageSize  int64
	PageStart int64
	Assets    [] Asset
}

func (asset *Asset) GetValue(attributename string) interface{} {
	for k, v := range asset.Attributes {
		if strings.Compare(k, attributename) == 0 {
			return v.Value
		}
	}
	return nil
}

func (asset *Asset) GetAsReflectValue(attributename string) reflect.Value {
	return reflect.ValueOf(asset.GetValue(attributename))
}

func (asset *Asset) GetAsStringValue(attributename string) string {
	return fmt.Sprintf("%v", reflect.ValueOf(asset.GetValue(attributename)))
}

func (trend *Trend) Merge(newtrend Trend) {
	if trend.Result == nil || len(trend.Result) == 0 {
		trend.Result = newtrend.Result
		return
	}
	for _, newnode := range newtrend.Result {
		currnodehassamenode := false
		for _, curnode := range trend.Result {
			if curnode.Name == newnode.Name {
				currnodehassamenode = true
				if len(curnode.Series) == 0 && len(newnode.Series) > 0 {
					curnode.Series = newnode.Series
				} else {
					for i, newserie := range newnode.Series {
						curnode.Series[i].Value += newserie.Value
					}
				}
			}
		}
		if !currnodehassamenode {
			trend.Result = append(trend.Result, newnode)
		}
	}
}

func (trend *Trend) GetNode(status string) *Node {
	for _, node := range trend.Result {
		if node.Name == status {
			return node
		}
	}
	return &Node{}
}

func (trend *Trend) GetValue(nodename string, serieindex int) Serie {
	n := trend.GetNode(nodename)
	if len(n.Series) > 0 {
		return n.Series[serieindex]
	}
	s := Serie{}
	s.Date = ""
	s.Value = 0

	return s
}

func (trend *Trend) GetFirstAvailableDateValue(serieindex int) string {
	for _, n := range trend.Result {
		if len(n.Series) > 0 {
			return n.Series[serieindex].Date
		}
	}
	return ""
}

func (trend *Trend) GetNodeLength() int {
	for _, n := range trend.Result {
		if len(n.Series) > 0 {
			return len(n.Series)
		}
	}
	return 0
}
