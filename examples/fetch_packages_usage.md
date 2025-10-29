# Fetch Packages API - Usage Examples

This document provides practical examples of how to use the new dynamic `fetch-packages` API in various scenarios.

## Basic Usage

### 1. Simple Package Search (HTML Response)
Perfect for HTMX integration in web interfaces.

```bash
curl "http://localhost:8080/fetch-packages?q=gin"
```

**Response:** HTML snippet ready for insertion into the DOM.

### 2. JSON API Search
For programmatic access and API integrations.

```bash
curl "http://localhost:8080/fetch-packages?q=gin&format=json"
```

**Response:**
```json
{
  "success": true,
  "query": "gin",
  "provider": "pkg.go.dev",
  "count": 15,
  "results": [
    {
      "path": "github.com/gin-gonic/gin",
      "synopsis": "Gin is a HTTP web framework written in Go (Golang)"
    }
  ],
  "cache_hit": false,
  "timestamp": 1703123456
}
```

## Advanced Usage

### 3. Custom Result Limit
Limit the number of results returned.

```bash
curl "http://localhost:8080/fetch-packages?q=echo&format=json&limit=5"
```

### 4. Using Fallback Provider
When you want to use the local fallback instead of pkg.go.dev.

```bash
curl "http://localhost:8080/fetch-packages?q=web&provider=fallback&format=json"
```

### 5. Disable Caching
Force a fresh search without using cached results.

```bash
curl "http://localhost:8080/fetch-packages?q=fiber&format=json&cache=false"
```

## Frontend Integration Examples

### 6. HTMX Integration (Current Implementation)
The existing go-ctl interface uses this pattern:

```html
<input type="text" 
       name="q"
       placeholder="Search pkg.go.dev for packages..."
       hx-get="/fetch-packages"
       hx-trigger="keyup changed delay:500ms"
       hx-target="#search-results"
       hx-swap="innerHTML">

<div id="search-results"></div>
```

### 7. JavaScript/Fetch API
For custom frontend implementations:

```javascript
async function searchPackages(query, options = {}) {
  const params = new URLSearchParams({
    q: query,
    format: 'json',
    limit: options.limit || 20,
    provider: options.provider || 'pkg.go.dev',
    cache: options.cache !== false ? 'true' : 'false'
  });

  try {
    const response = await fetch(`/fetch-packages?${params}`);
    const data = await response.json();
    
    if (data.success) {
      return data.results;
    } else {
      throw new Error(data.error || 'Search failed');
    }
  } catch (error) {
    console.error('Package search error:', error);
    return [];
  }
}

// Usage examples
searchPackages('gin').then(packages => {
  console.log('Found packages:', packages);
});

searchPackages('echo', { limit: 5, cache: false }).then(packages => {
  packages.forEach(pkg => {
    console.log(`${pkg.path}: ${pkg.synopsis}`);
  });
});
```

### 8. React/Vue.js Integration
Example React hook for package searching:

```javascript
import { useState, useEffect, useCallback } from 'react';

function usePackageSearch(query, options = {}) {
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const searchPackages = useCallback(async (searchQuery) => {
    if (!searchQuery.trim()) {
      setResults([]);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const params = new URLSearchParams({
        q: searchQuery,
        format: 'json',
        limit: options.limit || 15,
        provider: options.provider || 'pkg.go.dev'
      });

      const response = await fetch(`/fetch-packages?${params}`);
      const data = await response.json();

      if (data.success) {
        setResults(data.results);
      } else {
        setError(data.error || 'Search failed');
        setResults([]);
      }
    } catch (err) {
      setError(err.message);
      setResults([]);
    } finally {
      setLoading(false);
    }
  }, [options.limit, options.provider]);

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      searchPackages(query);
    }, 500); // Debounce

    return () => clearTimeout(timeoutId);
  }, [query, searchPackages]);

  return { results, loading, error };
}

// Usage in component
function PackageSearch() {
  const [query, setQuery] = useState('');
  const { results, loading, error } = usePackageSearch(query);

  return (
    <div>
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search for Go packages..."
      />
      
      {loading && <div>Searching...</div>}
      {error && <div>Error: {error}</div>}
      
      <ul>
        {results.map((pkg, index) => (
          <li key={index}>
            <strong>{pkg.path}</strong>
            <p>{pkg.synopsis}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}
```

## Testing Examples

### 9. Performance Testing
Test cache performance and response times:

```bash
# First request (cache miss)
time curl -s "http://localhost:8080/fetch-packages?q=gin&format=json" > /dev/null

# Second request (cache hit)
time curl -s "http://localhost:8080/fetch-packages?q=gin&format=json" > /dev/null
```

### 10. Batch Testing
Test multiple queries in sequence:

```bash
#!/bin/bash
queries=("gin" "echo" "fiber" "chi" "mux" "gorilla")

for query in "${queries[@]}"; do
  echo "Testing query: $query"
  response=$(curl -s "http://localhost:8080/fetch-packages?q=$query&format=json")
  count=$(echo "$response" | grep -o '"count":[0-9]*' | cut -d':' -f2)
  echo "  Found $count packages"
  echo
done
```

## Error Handling Examples

### 11. Graceful Error Handling
Handle various error scenarios:

```javascript
async function robustPackageSearch(query) {
  try {
    const response = await fetch(`/fetch-packages?q=${encodeURIComponent(query)}&format=json`);
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    
    const data = await response.json();
    
    if (data.success) {
      return {
        success: true,
        packages: data.results,
        cached: data.cache_hit,
        count: data.count
      };
    } else {
      return {
        success: false,
        error: data.error || 'Unknown error occurred',
        packages: []
      };
    }
  } catch (error) {
    console.error('Network or parsing error:', error);
    return {
      success: false,
      error: 'Network error or server unavailable',
      packages: []
    };
  }
}

// Usage with error handling
robustPackageSearch('gin').then(result => {
  if (result.success) {
    console.log(`Found ${result.count} packages (cached: ${result.cached})`);
    result.packages.forEach(pkg => {
      console.log(`- ${pkg.path}: ${pkg.synopsis}`);
    });
  } else {
    console.error('Search failed:', result.error);
    // Show user-friendly error message
    showErrorMessage('Unable to search packages. Please try again.');
  }
});
```

## Migration Examples

### 12. Migrating from Legacy search-packages
Old code using the legacy endpoint:

```javascript
// OLD - still works but deprecated
fetch('/search-packages?q=gin')
  .then(response => response.text())
  .then(html => {
    document.getElementById('results').innerHTML = html;
  });
```

New recommended approach:

```javascript
// NEW - using the dynamic API
fetch('/fetch-packages?q=gin&format=json')
  .then(response => response.json())
  .then(data => {
    if (data.success) {
      renderResults(data.results);
    } else {
      showError(data.error);
    }
  });
```

### 13. Backward Compatibility Test
Ensure your migration doesn't break existing functionality:

```bash
# Test that legacy endpoint still works
curl "http://localhost:8080/search-packages?q=gin"

# Test that new endpoint produces same HTML output
curl "http://localhost:8080/fetch-packages?q=gin&format=html"

# Compare responses (they should be identical)
diff <(curl -s "http://localhost:8080/search-packages?q=gin") \
     <(curl -s "http://localhost:8080/fetch-packages?q=gin&format=html")
```

## Best Practices

### 14. Optimal Configuration
Recommended settings for different use cases:

```javascript
// For interactive search (user typing)
const interactiveConfig = {
  format: 'json',
  limit: 10,
  cache: true // Use cache for better UX
};

// For comprehensive search
const comprehensiveConfig = {
  format: 'json',
  limit: 50,
  cache: false // Always fresh results
};

// For fast fallback search
const fallbackConfig = {
  format: 'json',
  provider: 'fallback',
  limit: 20,
  cache: true
};
```

### 15. Debouncing Implementation
Prevent excessive API calls during user input:

```javascript
class DebouncedPackageSearch {
  constructor(delay = 500) {
    this.delay = delay;
    this.timeoutId = null;
  }

  search(query, callback) {
    // Clear previous timeout
    if (this.timeoutId) {
      clearTimeout(this.timeoutId);
    }

    // Set new timeout
    this.timeoutId = setTimeout(async () => {
      try {
        const response = await fetch(`/fetch-packages?q=${encodeURIComponent(query)}&format=json`);
        const data = await response.json();
        callback(data);
      } catch (error) {
        callback({ success: false, error: error.message, results: [] });
      }
    }, this.delay);
  }
}

// Usage
const searcher = new DebouncedPackageSearch(500);
searcher.search('gin', (results) => {
  console.log('Search results:', results);
});
```

## Integration with go.mod Generation

### 16. Adding Selected Packages to Project
Example of how the API integrates with the project generation:

```javascript
let selectedPackages = [];

function addPackage(packagePath) {
  if (!selectedPackages.includes(packagePath)) {
    selectedPackages.push(packagePath);
    updateUI();
  }
}

function generateProject() {
  const projectData = {
    name: document.getElementById('project-name').value,
    packages: selectedPackages,
    // ... other project options
  };

  fetch('/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(projectData)
  })
  .then(response => response.blob())
  .then(blob => {
    // Download generated project
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${projectData.name}.zip`;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
  });
}
```

---

*These examples demonstrate the flexibility and power of the new fetch-packages API. The API is designed to be backward-compatible while providing enhanced functionality for modern web applications.*