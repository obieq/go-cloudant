package cloudant

import (
	"bytes"
	"encoding/json"
)

// Container for cloudant index information
// http://docs.cloudant.com/api/cloudant-query.html?highlight=query#list-all-indexes
type Index struct {
	DDocId     string          `json:"ddoc"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Definition indexDefinition `json:"def"`
}

type indexList struct {
	Indices []Index `json:"indexes"`
}

type indexDefinition struct {
	Fields []map[string]string `json:"fields"`
}

// Create an index on the specified field names
//  http://docs.cloudant.com/api/cloudant-query.html?highlight=query#creating-a-new-index
//
//  Options:
//    1) ddoc_name:  name of the design document.  auto-generated if not specified
//    2) index_name: name of the index auto-generate if not specified
//    3) type:       currently, only supported type is json (default).
//                   per the docs, full text and geospatial will be available at some point
func (db *Database) CreateIndex(fields []string, opts map[string]string) (CloudantDocumentResponse, error) {
	cdr := CloudantDocumentResponse{}

	json, _ := json.Marshal(fields) // given that fields is a string array, json.Marshal shouldn't raise an error

	// set fields names to be indexed
	str := `"index": {"fields": ` + string(json) + `}`

	// set design document name?
	if ddoc_name := opts["ddoc_name"]; ddoc_name != "" {
		str += `, "ddoc": "` + ddoc_name + `"`
	}

	// set index name?
	if index_name := opts["index_name"]; index_name != "" {
		str += `, "name": "` + index_name + `"`
	}

	body := []byte("{" + str + "}")

	resp, err := db.client.doRequest("POST", "/"+db.Name()+"/_index", nil, bytes.NewReader(body))
	err = db.client.handleResponse(resp, err, 200, &cdr)

	return cdr, err
}

func (db *Database) GetIndices() ([]Index, error) {
	il := indexList{}

	resp, err := db.client.doRequest("GET", "/"+db.Name()+"/_index", nil, nil)
	err = db.client.handleResponse(resp, err, 200, &il)

	return il.Indices, err
}

func (db *Database) GetIndexByName(name string) (*Index, error) {
	var index *Index = nil
	var err error
	var indices []Index

	if indices, err = db.GetIndices(); err != nil {
		return index, err
	}

	for _, idx := range indices {
		if idx.Name == name {
			index = &idx
			break
		}
	}

	// given that we don't have a cloudant error response for getting an index by a non-existent name,
	// we'll build an error response and return
	if index == nil {
		return index, NotFoundError()
	}

	return index, err
}

func (db *Database) GetIndexByDDocName(ddocId string) (*DesignDocument, error) {
	return db.GetDesignDocument(ddocId)
}

func (db *Database) DeleteIndexByName(name string) (CloudantDocumentResponse, error) {
	cdr := CloudantDocumentResponse{}
	var index *Index
	var err error

	// get index by name in order to grab the ddoc id
	if index, err = db.GetIndexByName(name); err != nil {
		return cdr, err
	}

	uri := "/" + db.Name() + "/_index/" + index.DDocId + "/" + index.Type + "/" + name
	resp, err := db.client.doRequest("DELETE", uri, nil, nil)
	err = db.client.handleResponse(resp, err, 200, &cdr)

	return cdr, err
}
