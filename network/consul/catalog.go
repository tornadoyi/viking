package consul

import (
	_consul "github.com/hashicorp/consul/api"
	"sync"
	"time"
)

type Catalog struct {
	*_consul.Catalog
	client	*Client
	timer				*time.Timer
	focus				map[string]bool
	services			map[string][]*CatalogService
	mutex				*sync.RWMutex
}

func NewCatalog(client *Client) *Catalog {
	return &Catalog{
		Catalog:    client.client.Catalog(),
		client:   client,
		focus:    make(map[string]bool),
		services: make(map[string][]*CatalogService),
		mutex:    &sync.RWMutex{},
	}
}


type CatalogService = _consul.CatalogService
type CatalogDeregistration = _consul.CatalogDeregistration