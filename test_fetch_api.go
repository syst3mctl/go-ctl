package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PackageFetchResponse represents the API response structure
type PackageFetchResponse struct {
	Success   bool             `json:"success"`
	Query     string           `json:"query"`
	Provider  string           `json:"provider"`
	Count     int              `json:"count"`
	Results   []PkgGoDevResult `json:"results"`
	Error     string           `json:"error,omitempty"`
	CacheHit  bool             `json:"cache_hit"`
	Timestamp int64            `json:"timestamp"`
}

// PkgGoDevResult represents a package search result
type PkgGoDevResult struct {
	Path     string `json:"path"`
	Synopsis string `json:"synopsis"`
}

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("🧪 Testing fetch-packages API")
	fmt.Println("==============================")

	// Test cases
	tests := []struct {
		name        string
		endpoint    string
		params      map[string]string
		expectJSON  bool
		expectError bool
	}{
		{
			name:     "Basic HTML search (legacy compatibility)",
			endpoint: "/search-packages",
			params:   map[string]string{"q": "gin"},
		},
		{
			name:     "Basic HTML search (new endpoint)",
			endpoint: "/fetch-packages",
			params:   map[string]string{"q": "gin"},
		},
		{
			name:       "JSON API search",
			endpoint:   "/fetch-packages",
			params:     map[string]string{"q": "gin", "format": "json"},
			expectJSON: true,
		},
		{
			name:       "JSON with custom limit",
			endpoint:   "/fetch-packages",
			params:     map[string]string{"q": "echo", "format": "json", "limit": "5"},
			expectJSON: true,
		},
		{
			name:       "Fallback provider",
			endpoint:   "/fetch-packages",
			params:     map[string]string{"q": "web", "provider": "fallback", "format": "json"},
			expectJSON: true,
		},
		{
			name:       "Cache disabled",
			endpoint:   "/fetch-packages",
			params:     map[string]string{"q": "fiber", "format": "json", "cache": "false"},
			expectJSON: true,
		},
		{
			name:       "Empty query",
			endpoint:   "/fetch-packages",
			params:     map[string]string{"q": "", "format": "json"},
			expectJSON: true,
		},
		{
			name:        "Invalid provider",
			endpoint:    "/fetch-packages",
			params:      map[string]string{"q": "test", "provider": "invalid", "format": "json"},
			expectJSON:  true,
			expectError: true,
		},
	}

	// Run tests
	successCount := 0
	for i, test := range tests {
		fmt.Printf("\n%d. %s\n", i+1, test.name)
		fmt.Println("   " + strings.Repeat("-", len(test.name)))

		if runTest(test) {
			successCount++
			fmt.Println("   ✅ PASSED")
		} else {
			fmt.Println("   ❌ FAILED")
		}
	}

	// Performance test
	fmt.Printf("\n🚀 Performance Test\n")
	fmt.Println("   ----------------")
	runPerformanceTest()

	// Summary
	fmt.Printf("\n📊 Test Summary\n")
	fmt.Println("   =============")
	fmt.Printf("   Passed: %d/%d tests\n", successCount, len(tests))
	if successCount == len(tests) {
		fmt.Println("   🎉 All tests passed!")
	} else {
		fmt.Printf("   ⚠️  %d tests failed\n", len(tests)-successCount)
	}
}

func runTest(test struct {
	name        string
	endpoint    string
	params      map[string]string
	expectJSON  bool
	expectError bool
}) bool {
	// Build URL
	u := fmt.Sprintf("%s%s", baseURL, test.endpoint)
	if len(test.params) > 0 {
		params := url.Values{}
		for k, v := range test.params {
			params.Add(k, v)
		}
		u += "?" + params.Encode()
	}

	fmt.Printf("   URL: %s\n", u)

	// Make request
	start := time.Now()
	resp, err := http.Get(u)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	fmt.Printf("   Status: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("   Duration: %v\n", duration)

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("   Error reading body: %v\n", err)
		return false
	}

	// Check response
	if test.expectJSON {
		return validateJSONResponse(body, test.expectError)
	} else {
		return validateHTMLResponse(body)
	}
}

func validateJSONResponse(body []byte, expectError bool) bool {
	var response PackageFetchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("   JSON parse error: %v\n", err)
		return false
	}

	fmt.Printf("   Query: '%s'\n", response.Query)
	fmt.Printf("   Provider: %s\n", response.Provider)
	fmt.Printf("   Success: %t\n", response.Success)
	fmt.Printf("   Count: %d\n", response.Count)
	fmt.Printf("   Cache Hit: %t\n", response.CacheHit)

	if expectError {
		if !response.Success && response.Error != "" {
			fmt.Printf("   Expected Error: %s\n", response.Error)
			return true
		} else {
			fmt.Println("   Expected error but got success")
			return false
		}
	}

	if !response.Success {
		fmt.Printf("   Unexpected Error: %s\n", response.Error)
		return false
	}

	// Show first few results
	if len(response.Results) > 0 {
		fmt.Println("   Sample Results:")
		for i, result := range response.Results {
			if i >= 3 { // Show max 3 results
				fmt.Printf("   ... and %d more\n", len(response.Results)-3)
				break
			}
			fmt.Printf("     - %s: %s\n", result.Path, truncate(result.Synopsis, 60))
		}
	}

	return true
}

func validateHTMLResponse(body []byte) bool {
	html := string(body)
	fmt.Printf("   Response Length: %d bytes\n", len(body))

	// Basic HTML validation
	if len(html) == 0 {
		fmt.Println("   Empty response (acceptable for empty query)")
		return true
	}

	// Check for basic HTML structure
	if containsAny(html, []string{"<div", "class=", "hx-post"}) {
		fmt.Println("   Valid HTML structure detected")
		return true
	}

	fmt.Println("   Warning: Response doesn't look like expected HTML")
	if len(html) < 200 {
		fmt.Printf("   Response preview: %s\n", html)
	}
	return false
}

func runPerformanceTest() {
	queries := []string{"gin", "echo", "fiber", "chi", "mux"}

	fmt.Println("   Testing cache performance...")

	totalDuration := time.Duration(0)
	for i, query := range queries {
		// First request (cache miss)
		start := time.Now()
		resp, err := http.Get(fmt.Sprintf("%s/fetch-packages?q=%s&format=json", baseURL, query))
		duration1 := time.Since(start)
		totalDuration += duration1

		if err != nil {
			fmt.Printf("   Error with query '%s': %v\n", query, err)
			continue
		}
		resp.Body.Close()

		// Second request (cache hit)
		start = time.Now()
		resp, err = http.Get(fmt.Sprintf("%s/fetch-packages?q=%s&format=json", baseURL, query))
		duration2 := time.Since(start)

		if err != nil {
			fmt.Printf("   Error with cached query '%s': %v\n", query, err)
			continue
		}
		resp.Body.Close()

		improvement := duration1 - duration2
		improvementPct := float64(improvement) / float64(duration1) * 100

		fmt.Printf("   %s: %v → %v (%.1f%% faster)\n", query, duration1, duration2, improvementPct)
	}

	avgDuration := totalDuration / time.Duration(len(queries))
	fmt.Printf("   Average response time: %v\n", avgDuration)

	if avgDuration < 2*time.Second {
		fmt.Println("   ✅ Performance: Good")
	} else if avgDuration < 5*time.Second {
		fmt.Println("   ⚠️  Performance: Acceptable")
	} else {
		fmt.Println("   ❌ Performance: Poor")
	}
}

// Helper functions
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if contains(s, substr) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
