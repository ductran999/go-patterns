package gravitino

import (
	"fmt"
)

type TargetAPI interface {
	Metalakes() MetalakeAction
	Metalake(name string) MetalakeScope
}

type MetalakeAction interface {
	Get()
	Create()
	Delete()
}

func (t *targetImpl) Metalakes() MetalakeAction {
	return &metalakeActionImpl{
		target: t,
	}
}

type metalakeActionImpl struct {
	target *targetImpl
}

func (m *metalakeActionImpl) Get() {
	path := fmt.Sprintf("%s/api/metalakes", m.target.baseURL)
	fmt.Println("-> Getting all metalakes...")
	m.target.client.doRequest(path)
}

func (m *metalakeActionImpl) Create() {
	fmt.Println("-> Creating a metalake...")
}

func (m *metalakeActionImpl) Delete() {
	fmt.Println("-> Deleting a metalake...")
}
