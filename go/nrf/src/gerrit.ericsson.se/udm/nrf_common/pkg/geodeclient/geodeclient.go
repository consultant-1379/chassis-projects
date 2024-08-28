package geodeclient

import (
	"gerrit.ericsson.se/udm/common/pkg/log"

	geode "github.com/gemfire/geode-go-client"
	. "github.com/gemfire/geode-go-client/query"
	"github.com/gemfire/geode-go-client/connector"
)

// GeodeClient struct
type GeodeClient struct {
	client *geode.Client
}

var instance *GeodeClient

// GetInstance for GeodeClient
func GetInstance() *GeodeClient {
	if instance == nil {
		instance = newGeodeClient()
	}
	return instance
}

func newGeodeClient() *GeodeClient {
	pool := connector.NewPool()
	pool.AddServer("eric-nrf-kvdb-ag-server", 40404)
	// Optionally add user credentials
	//pool.AddCredentials("jbloggs", "t0p53cr3t")

	conn := connector.NewConnector(pool)
	client := geode.NewGeodeClient(conn)

	return &GeodeClient{client}
}

// Put data into a region. key and value must be a supported type.
func (c *GeodeClient) Put(region string, key, value interface{}) {
	err := c.client.Put(region, key, value)
	if err != nil {
		log.Error(err)
	}
}

// Get an entry from a region using the specified key.
func (c *GeodeClient) Get(region string, key interface{}) (interface{}, error) {
	return c.client.Get(region, key)
}

// Remove an entry for a region.
func (c *GeodeClient) Remove(region string, key interface{}) error {
	return c.client.Remove(region, key)
}

// QueryForListResult Execute a query, returning a list of results
func (c *GeodeClient) QueryForListResult(query string) ([]interface{}, error){
	q := NewQuery(query)
	return c.client.QueryForListResult(q)
}
