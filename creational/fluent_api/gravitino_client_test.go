package gravitino_test

import (
	"fmt"
	gravitino "patterns/creational/fluent_api"
	"testing"
)

func Test_FluentAPI(t *testing.T) {
	client := gravitino.NewClient()

	clusterURL := "https://gravitino.cluster1.local"
	token := "abc.def.ghi"

	fmt.Println("--- TEST 1: (Get List Metalakes) ---")
	client.WithTarget(clusterURL, token).Metalakes().Get()

	fmt.Println("\n--- TEST 2: (Get List Catalogs of one Metalake) ---")
	client.WithTarget(clusterURL, token).Metalake("production_lake").Catalogs()
}
