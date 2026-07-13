package gravitino

import (
	"log"
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
	path := m.target.baseURL + "/api/metalakes"
	log.Println("-> Getting all metalakes...")
	m.target.client.doRequest(path)
}

func (m *metalakeActionImpl) Create() {
	log.Println("-> Creating a metalake...")
}

func (m *metalakeActionImpl) Delete() {
	log.Println("-> Deleting a metalake...")
}
