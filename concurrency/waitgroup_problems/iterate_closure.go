package main

import (
	"fmt"
	"sync"
)

/**
 * TOPIC: Loop Variable Semantics & Goroutine Capture (Go 1.22+)
 *
 * HISTORICAL CONTEXT:
 * Prior to Go 1.22, the iteration variable (salutation) was reused across
 * all iterations. This led to a common pitfall where all goroutines captured
 * the same memory address, typically printing the final element of the slice.
 *
 * THE GO 1.22 CHANGE:
 * The language now implements "per-iteration" variable scoping. Each loop
 * iteration creates a new instance of the variable, making it safe to use
 * inside closures/goroutines without manual shadowing or passing parameters.
 */

func main() {
	var wg sync.WaitGroup
	salutations := []string{"hello", "greetings", "good day"}

	for _, salutation := range salutations {
		wg.Add(1)

		// MODERN APPROACH (Go 1.22+)
		// No manual parameter passing required. Each goroutine captures
		// its own unique 'salutation' instance.
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
		//
		// LEGACY APPROACH (Pre-Go 1.22):
		// To achieve the same result in older versions, you had to pass the
		// variable as an argument to avoid referencing the same memory address:
		//
		// go func(s string) {
		//     defer wg.Done()
		//     fmt.Println(s)
		// }(salutation)
		//
	}

	wg.Wait()
}
