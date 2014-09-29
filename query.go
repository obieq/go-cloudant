package cloudant

import (
	"bytes"
	"encoding/json"
)

// Container for cloudant query selectors
// http://docs.cloudant.com/api/cloudant-query.html?highlight=query#selector-syntax
type Query struct {
	Selector map[string]interface{} `json:"selector"`
	Sorters  []map[string]string    `json:"sort"`
	Limit    int                    `json:"limit"`
	Skip     int                    `json:"skip"`
}

type QueryResults struct {
	Docs interface{} `json:"docs"`
}

func NewQuery() *Query {
	return &Query{Selector: make(map[string]interface{}), Sorters: []map[string]string{}, Limit: 25, Skip: 0}
}

func (q *Query) Sort(field string, isAscending bool) {
	sort := make(map[string]string)
	direction := "asc"

	if !isAscending {
		direction = "desc"
	}

	sort[field] = direction
	q.Sorters = append(q.Sorters, sort)
}

func (db *Database) Query(q *Query, results interface{}) error {
	j, err := json.Marshal(q)
	if err != nil {
		return err
	}

	body := []byte(string(j))
	resp, err := db.client.doRequest("POST", "/"+db.Name()+"/_find", nil, bytes.NewReader(body))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return newCloudantError(resp)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	//fmt.Printf(buf.String()) // NOTE: uncomment this line when debugging results

	// for test coverage purposes, performing nil check using non-idiomatic pattern
	var objmap map[string]*json.RawMessage
	if err == nil {
		err = json.Unmarshal(buf.Bytes(), &objmap)
	}

	if err == nil {
		err = json.Unmarshal(*objmap["docs"], results)
	}

	return err
}
