package gravitino_test

import (
	"log"
	gravitino "patterns/creational/fluent_api"
	"testing"
)

func Test_FluentAPI(t *testing.T) {
	t.Parallel()

	client := gravitino.NewClient()

	clusterURL := "https://gravitino.cluster1.local"
	token := "abc.def.ghi"

	log.Println("--- TEST 1: (Get List Metalakes) ---")
	client.WithTarget(clusterURL, token).Metalakes().Get()

	log.Println("\n--- TEST 2: (Get List Catalogs of one Metalake) ---")
	client.WithTarget(clusterURL, token).Metalake("production_lake").Catalogs()
}
