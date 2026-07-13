package gravitino

import (
	"fmt"
	"log"
)

type CatalogAction = MetalakeAction

type MetalakeScope interface {
	Catalogs() CatalogAction
}

func (t *targetImpl) Metalake(name string) MetalakeScope {
	return &metalakeScopeImpl{
		target: t,
		name:   name,
	}
}

type metalakeScopeImpl struct {
	target *targetImpl
	name   string
}

func (ms *metalakeScopeImpl) Catalogs() CatalogAction {
	return &catalogActionImpl{
		metalakeScope: ms,
	}
}

type catalogActionImpl struct {
	metalakeScope *metalakeScopeImpl
}

func (c *catalogActionImpl) Get() {
	baseURL := c.metalakeScope.target.baseURL
	metalakeName := c.metalakeScope.name

	path := fmt.Sprintf("%s/api/metalakes/%s/catalogs", baseURL, metalakeName)
	log.Println("-> Getting all catalogs for metalake", metalakeName)

	c.metalakeScope.target.client.doRequest(path)
}

func (c *catalogActionImpl) Create() {
	log.Println("-> Creating a catalog in metalake", c.metalakeScope.name)
}

func (c *catalogActionImpl) Delete() {
	log.Println("-> Deleting a catalog in metalake", c.metalakeScope.name)
}
