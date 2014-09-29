package cloudant

import (
	"bytes"
	"encoding/json"
)

type CloudantDocumentInterfacer interface {
	Id() string
	SetId(string)
	Revision() string
}

type cloudantDocument struct {
	DocId       string `json:"_id,omitempty"`
	DocRevision string `json:"_rev,omitempty"`
}

type CloudantDocument struct {
	cloudantDocument
}

type CloudantDocumentResponse struct {
	Id       string `json:"id"`
	Revision string `json:"rev"`
	Ok       string `json:"ok"`
	Result   string `json:"result"`
}

func (doc *CloudantDocument) Id() string {
	return doc.DocId
}

func (doc *CloudantDocument) SetId(id string) {
	doc.DocId = id
}

func (doc *CloudantDocument) Revision() string {
	return doc.DocRevision
}

func (db *Database) GetDocument(id string, doc interface{}) error {
	resp, err := db.client.doRequest("GET", "/"+db.Name()+"/"+id, nil, nil)
	return db.client.handleResponse(resp, err, 200, doc)
}

func (db *Database) CreateDocument(doc interface{}, isBatch bool) (CloudantDocumentResponse, error) {
	cdr := CloudantDocumentResponse{}
	uri := "/" + db.Name()
	successStatusCode := 201

	j, err := json.Marshal(doc)
	if err != nil {
		return cdr, err
	}

	if isBatch {
		uri += "?batch=ok"
		successStatusCode = 202
	}

	resp, err := db.client.doRequest("POST", uri, nil, bytes.NewReader(j))
	err = db.client.handleResponse(resp, err, successStatusCode, &cdr)
	return cdr, err
}

func (db *Database) UpdateDocument(doc CloudantDocumentInterfacer, isBatch bool) (CloudantDocumentResponse, error) {
	cdr := CloudantDocumentResponse{}
	uri := "/" + db.Name() + "/" + doc.Id() + "?rev=" + doc.Revision()
	successStatusCode := 201

	j, err := json.Marshal(doc)
	if err != nil {
		return cdr, err
	}

	if isBatch {
		uri += "&batch=ok"
		successStatusCode = 202
	}

	resp, err := db.client.doRequest("PUT", uri, nil, bytes.NewReader(j))
	err = db.client.handleResponse(resp, err, successStatusCode, &cdr)
	return cdr, err
}

func (db *Database) DeleteDocument(id string, revision string) (CloudantDocumentResponse, error) {
	cdr := CloudantDocumentResponse{}

	resp, err := db.client.doRequest("DELETE", "/"+db.Name()+"/"+id+"?rev="+revision, nil, nil)
	err = db.client.handleResponse(resp, err, 200, &cdr)
	return cdr, err
}
