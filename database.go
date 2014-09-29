package cloudant

type Database struct {
	name   string
	client *Client
}

func (db *Database) Name() string {
	return db.name
}

func (db *Database) Client() *Client {
	return db.client
}

func NewDatabase(name string, client *Client) *Database {
	return &Database{name: name, client: client}
}

func (client *Client) GetDatabase(name string) Database {
	return Database{name: name, client: client}
}

func (client *Client) ListDatabases() ([]string, error) {
	dbs := []string{}

	resp, err := client.doRequest("GET", "/_all_dbs", nil, nil)
	err = client.handleResponse(resp, err, 200, &dbs)

	return dbs, err
}

func (client *Client) CreateDatabase(name string) (CloudantDocumentResponse, error) {
	cdr := CloudantDocumentResponse{}
	resp, err := client.doRequest("PUT", "/"+name, nil, nil)
	err = client.handleResponse(resp, err, 201, &cdr)

	return cdr, err
}

func (client *Client) DeleteDatabase(name string) (CloudantDocumentResponse, error) {
	cdr := CloudantDocumentResponse{}
	resp, err := client.doRequest("DELETE", "/"+name, nil, nil)
	err = client.handleResponse(resp, err, 200, &cdr)

	return cdr, err
}
