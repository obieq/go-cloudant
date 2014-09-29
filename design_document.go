package cloudant

// Container for cloudant design document information
// http://docs.cloudant.com/api/design-documents-get-put-delete-copy.html
type DesignDocument struct {
	CloudantDocument
	Language string                        `json:"language"`
	Views    map[string]designDocumentView `json:"views"`
}

type designDocumentView struct {
	Map     designDocumentViewMap    `json:"map"`
	Reduce  string                   `json:"reduce"`
	Options designDocumentViewOption `json:"options"`
}

type designDocumentViewMap struct {
	Fields map[string]string `json:"fields"`
}

type designDocumentViewOption struct {
	Definition designDocumentViewOptionDefinition `json:"def"`
	W          int                                `json:"w"`
}

type designDocumentViewOptionDefinition struct {
	Fields []string `json:"fields"`
}

func (db *Database) GetDesignDocument(id string) (*DesignDocument, error) {
	doc := &DesignDocument{}
	resp, err := db.client.doRequest("GET", "/"+db.Name()+"/_design/"+id, nil, nil)
	err = db.client.handleResponse(resp, err, 200, doc)

	return doc, err
}

func (ddoc *DesignDocument) ViewKeys() []string {
	keys := make([]string, 0, len(ddoc.Views))
	for k := range ddoc.Views {
		keys = append(keys, k)
	}

	return keys
}
