package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("🧪 Debug Test for fetch-packages API")
	fmt.Println("=====================================")

	// Test the fallback function directly
	fmt.Println("\n1. Testing fallback function...")
	results, err := searchPackagesFallback("gin")
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d results:\n", len(results))
		for i, result := range results {
			if i >= 3 {
				break
			}
			fmt.Printf("   - %s: %s\n", result.Path, result.Synopsis)
		}
	}

	// Test the new fetch function
	fmt.Println("\n2. Testing fetchPackagesFromPkgGoDev...")
	results, err = fetchPackagesFromPkgGoDev("gin", 5, true)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d results:\n", len(results))
		for i, result := range results {
			if i >= 3 {
				break
			}
			fmt.Printf("   - %s: %s\n", result.Path, result.Synopsis)
		}
	}

	// Start a simple test server
	fmt.Println("\n3. Starting test server...")
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		response := PackageFetchResponse{
			Success:   true,
			Query:     "test",
			Provider:  "test",
			Count:     1,
			Results:   []PkgGoDevResult{{Path: "test/package", Synopsis: "Test package"}},
			Timestamp: time.Now().Unix(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	go func() {
		fmt.Println("   Server starting on :8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	// Wait and test
	time.Sleep(1 * time.Second)

	fmt.Println("\n4. Testing local server...")
	resp, err := http.Get("http://localhost:8081/test")
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("   Status: %d\n", resp.StatusCode)

		var result PackageFetchResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("   JSON Error: %v\n", err)
		} else {
			fmt.Printf("   Response: Success=%v, Count=%d\n", result.Success, result.Count)
		}
	}

	fmt.Println("\n✅ Debug test complete")
}
