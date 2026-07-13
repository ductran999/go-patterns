package gravitino

import "fmt"

type CatalogAction interface {
	Get()
	Delete()
	Create()
}

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
	fmt.Printf("-> Getting all catalogs for metalake '%s'...\n", metalakeName)

	c.metalakeScope.target.client.doRequest(path)
}

func (c *catalogActionImpl) Create() {
	fmt.Printf("-> Creating a catalog in metalake '%s'...\n", c.metalakeScope.name)
}

func (c *catalogActionImpl) Delete() {
	fmt.Printf("-> Deleting a catalog in metalake '%s'...\n", c.metalakeScope.name)
}
