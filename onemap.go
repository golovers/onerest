package onerest

import (
	"fmt"
)

type OneMap struct {
	Total  int
	Values map[string][]Asset
}

func (ones *OneMap) Get(key string) []Asset {
	return ones.Values[key]
}

// Put the value into the map, the key of the map is the value of the given keyAttribute
func (ones *OneMap) Put(keyAttribute string, value Asset) {
	f := func(asset Asset) string {
		return asset.GetAsStringValue(keyAttribute)
	}
	ones.PutWithKeyTransformer(f, value)
}

type KeyTransformer func(asset Asset) string

// Calculate the key by applying the KeyTransformer to the given keyAttribute
func (ones *OneMap) PutWithKeyTransformer(keyTransformer KeyTransformer, value Asset) {
	ones.Total += 1

	key := keyTransformer(value)
	if ones.Values == nil {
		ones.Values = make(map[string][]Asset)
	}

	values := ones.Values[key]
	if values == nil {
		values = make([]Asset, 0, 0)
	}
	values = append(values, value)

	ones.Values[key] = values
}

// Put the value to the map.
// The key will be a string that combine values of the given keyAttributes
func (ones *OneMap) Puts(value Asset, keyAttributes ...string) {
	ones.Total += 1
	key := ""
	for i, k := range keyAttributes {
		s := value.GetValue(k)
		if i == 0 {
			key = fmt.Sprintf("%v", s)
		} else {
			key = key + "-" + fmt.Sprintf("%v", s)
		}
	}
	if ones.Values == nil {
		ones.Values = make(map[string][]Asset)
	}

	values := ones.Values[key]
	if values == nil {
		values = make([]Asset, 0, 0)
	}
	values = append(values, value)

	ones.Values[key] = values
}
