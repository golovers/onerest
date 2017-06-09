package onerest

import "strings"

type OneGroup struct {
	Assets [] *Asset
}

func (onegroup *OneGroup) AddAsset(asset Asset) {
	if onegroup.Assets == nil {
		onegroup.Assets = make([] *Asset, 0, 0)
	}
	onegroup.Assets = append(onegroup.Assets, &asset)
}

func (onegroup *OneGroup) AddAssets(assets []Asset) {
	for _, asset := range assets {
		onegroup.AddAsset(asset)
	}
}

// Group items by attributes
// If multiple keys are given, the key of the map will be follow this pattern: key1->key2->key3
func (onegroup *OneGroup) Aggregations(aggregationAttribute string, keyAttributes ...string) map[string]float64 {
	m := OneMap{}
	for _, asset := range onegroup.Assets {
		m.Puts(*asset, keyAttributes...)
	}
	result := make(map[string]float64)

	for key, assets := range m.Values {
		_, exist := result[key]
		if exist {
			result[key] = result[key] + AggerationByType(aggregationAttribute, assets)
		} else {
			result[key] = AggerationByType(aggregationAttribute, assets)
		}
	}
	RemoveUnexpectedCharactersInKeys(result)
	return result
}

// Group items by attributes
// If multiple keys are given, the key of the map will be follow this pattern: key1->key2->key3
func (onegroup *OneGroup) AggregationsWithKeyTransformer(aggregationAttribute string, keyTransformer KeyTransformer) map[string]float64 {
	m := OneMap{}
	for _, asset := range onegroup.Assets {
		m.PutWithKeyTransformer (keyTransformer, *asset)
	}
	result := make(map[string]float64)

	for key, assets := range m.Values {
		_, exist := result[key]
		if exist {
			result[key] = result[key] + AggerationByType(aggregationAttribute, assets)
		} else {
			result[key] = AggerationByType(aggregationAttribute, assets)
		}
	}
	RemoveUnexpectedCharactersInKeys(result)
	return result
}

func AggerationByType(attribute string, assets []Asset) float64 {
	if attribute == "count" {
		return float64(len(assets))
	} else {
		return SumByAttributeValue(attribute, assets)
	}
}

func SumByAttributeValue(attributename string, assets [] Asset) float64 {
	sum := float64(0)
	for _, asset := range assets {
		value := asset.GetAsReflectValue(attributename)
		if (value.IsValid()) {
			sum += float64(value.Float())
		}
	}
	return sum
}

func RemoveUnexpectedCharactersInKeys(data map[string] float64) {
	unexpectedKeys := []string{">", "<"}
	for k, _ := range data {
		newk := k
		for _, ch := range unexpectedKeys {
		newk = strings.Replace(newk, ch, "", 100)
		}

		if newk != k {
			data[newk] = data[k]
			delete(data, k)
		}
	}
}