package cloudant

type ClusterInfo struct {
	Couchdb       string `json:"couchdb"`
	Version       string `json:"version"`
	CloudantBuild string `json:"cloudant_build"`
}

func (client *Client) GetClusterInfo() (ClusterInfo, error) {
	ci := ClusterInfo{}

	resp, err := client.doRequest("GET", "", nil, nil)
	err = client.handleResponse(resp, err, 200, &ci)
	return ci, err
}
