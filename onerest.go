package onerest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

const DATA_END_POINT = "/rest-1.v1/Data/"
const REPORTS_GROUPED_PRIMARY_WORKITEM_BY_DATE = "/api/Reports/GroupedPrimaryWorkitemByDate"

type OneRest struct {
	Host        string
	AccessToken string
	Username    string
	Password    string
	EnableLog   bool
}

func (one *OneRest) getRestRequest(url string) *http.Request {
	urlstr := one.Host + url
	one.EnableLog = true
	if one.EnableLog {
		log.Println(urlstr)
	}
	request, _ := http.NewRequest("GET", urlstr, nil)

	header := http.Header{}
	header.Add("Accept", "application/json")
	if one.Username != "" && one.Password != "" {
		header.Add("Authorization", "Basic  "+base64.StdEncoding.EncodeToString(bytes.NewBufferString(one.Username+":"+one.Password).Bytes()))
	} else {
		header.Add("Authorization", "Bearer "+one.AccessToken)
	}
	request.Header = header
	return request
}

// Trend report
func (one *OneRest) trend(query map[string]string) (Trend, error) {
	urlstr := bytes.NewBufferString(REPORTS_GROUPED_PRIMARY_WORKITEM_BY_DATE)

	// Seems the escape function executed within request.SetQueryParams(wheres) got issue
	// Because the "startDate":"2017-01-01T00%3A00%3A00" in the VersionOne API does not allow to escape the first part "2017-01-01T" but REQUIRED to escape the second part as "00%3A00%3A00"
	// We also need to escape for the project attribute with "project":"OneService%3A50282"
	// So we need to parse the url manually
	if len(query) > 0 {
		urlstr.WriteString("?")
	}
	keys := make([]string, 0, len(query))
	for k, v := range query {
		keyvalue := ""
		tv := reflect.ValueOf(v)
		if strings.Compare(k, "startDate") == 0 {

			keyvalue = k + "=" + tv.String() + "T00%3A00%3A00"
		} else {
			keyvalue = k + "=" + strings.Replace(tv.String(), "-", "%3A", 100)
		}
		keys = append(keys, keyvalue)
	}
	urlstr.WriteString(strings.Join(keys, "&"))
	// Execute the request
	client := &http.Client{}
	response, err := client.Do(one.getRestRequest(urlstr.String()))
	report := Trend{}
	if err != nil {
		log.Fatal(err)
		return report, err
	} else {
		err = json.NewDecoder(response.Body).Decode(&report)
		if err != nil {
			return report, nil
		}
	}
	return report, nil
}

// Find all assets base on given asset type, select and condition map
func (one *OneRest) find(query *Query) ([]Asset, error) {
	req := one.getRestRequest(DATA_END_POINT + query.Build())
	client := &http.Client{}
	resp, err := client.Do(req)

	oneResponse := OneResponse{
		Assets: make([]Asset, 0, 0),
	}
	if err != nil {
		return oneResponse.Assets, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		err2 := json.NewDecoder(resp.Body).Decode(&oneResponse)
		if err2 != nil {
			log.Printf("Error unmarshalling: %v, url: %v", err2, req.URL)
		}
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		log.Printf("Error %v, url: %v", string(body), req.URL)
		return oneResponse.Assets, err
	}

	return oneResponse.Assets, nil
}
