# Dynamic Fetch Packages API

This document describes the new dynamic `fetch-packages` API that replaces the legacy `search-packages` endpoint with enhanced functionality, HTML scraping capabilities, and flexible configuration options.

## Overview

The `fetch-packages` API provides a powerful and flexible way to search for Go packages from pkg.go.dev using real-time HTML scraping. It supports both JSON and HTML response formats, making it perfect for both API integration and HTMX-based frontend applications.

### Key Features

- **Real HTML Scraping**: Uses goquery to scrape pkg.go.dev search results in real-time
- **Dual Response Formats**: Supports both JSON (for APIs) and HTML (for HTMX)
- **Flexible Configuration**: Customizable via query parameters
- **Caching System**: Built-in caching for improved performance
- **Fallback Support**: Graceful fallback when pkg.go.dev is unavailable
- **Backward Compatibility**: Legacy `/search-packages` endpoint still works

## Endpoints

### Primary Endpoint
```
GET /fetch-packages
```

### Legacy Endpoint (Backward Compatible)
```
GET /search-packages
```
*Automatically redirects to `/fetch-packages?format=html`*

## Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `q` | string | "" | Search query for packages (required) |
| `provider` | string | "pkg.go.dev" | Package provider (`pkg.go.dev` or `fallback`) |
| `format` | string | "html" | Response format (`html` or `json`) |
| `limit` | integer | 15 | Maximum number of results (1-100) |
| `cache` | boolean | true | Enable/disable caching |

## HTML Scraping Implementation

The API uses the following HTML scraping strategy with goquery:

```go
// Build search URL
url := fmt.Sprintf("https://pkg.go.dev/search?q=%s", query)

// Fetch and parse HTML
resp, err := http.Get(url)
doc, err := goquery.NewDocumentFromReader(resp.Body)

// Extract package information using CSS selectors
doc.Find("div.SearchSnippet").Each(func(i int, s *goquery.Selection) {
    // Get package path from data-test-id="snippet-title"
    pathText := s.Find("a[data-test-id='snippet-title'] .SearchSnippet-header-path").Text()
    path := strings.Trim(pathText, "()")
    
    // Get synopsis from SearchSnippet-synopsis class
    synopsis := s.Find(".SearchSnippet-synopsis").Text()
    
    // Get version from published span
    version := s.Find("span[data-test-id='snippet-published']").Parent().Find("strong").First().Text()
})
```

## Response Formats

### JSON Response
Used for API integration and programmatic access.

**Request:**
```bash
curl "http://localhost:8080/fetch-packages?q=gin&format=json&limit=3"
```

**Response:**
```json
{
  "success": true,
  "query": "gin",
  "provider": "pkg.go.dev",
  "count": 3,
  "results": [
    {
      "path": "github.com/gin-gonic/gin",
      "synopsis": "Gin is a HTTP web framework written in Go (Golang)"
    },
    {
      "path": "github.com/gin-contrib/cors",
      "synopsis": "Official CORS gin's middleware"
    },
    {
      "path": "github.com/gin-contrib/sessions",
      "synopsis": "Gin middleware for session management"
    }
  ],
  "cache_hit": false,
  "timestamp": 1703123456
}
```

### HTML Response
Used for HTMX integration and web interface compatibility.

**Request:**
```bash
curl "http://localhost:8080/fetch-packages?q=gin&format=html"
```

**Response:**
```html
<div class="space-y-2">
    <div class="flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 transition duration-150">
        <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-2 mb-1">
                <code class="text-sm font-mono text-blue-600 truncate">github.com/gin-gonic/gin</code>
            </div>
            <p class="text-sm text-gray-600 line-clamp-2">Gin is a HTTP web framework written in Go (Golang)</p>
        </div>
        <button class="ml-3 px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600 transition duration-150"
                hx-post="/add-package" 
                hx-vals='{"pkgPath": "github.com/gin-gonic/gin"}'
                hx-target="#selected-packages" 
                hx-swap="beforeend">
            Add
        </button>
    </div>
</div>
```

## Usage Examples

### 1. Basic Package Search
```bash
# HTML format (default)
curl "http://localhost:8080/fetch-packages?q=gin"

# JSON format
curl "http://localhost:8080/fetch-packages?q=gin&format=json"
```

### 2. Custom Limit and Provider
```bash
# Limit results to 5 packages
curl "http://localhost:8080/fetch-packages?q=echo&format=json&limit=5"

# Use fallback provider
curl "http://localhost:8080/fetch-packages?q=web&provider=fallback&format=json"
```

### 3. Cache Control
```bash
# Force fresh results (bypass cache)
curl "http://localhost:8080/fetch-packages?q=fiber&format=json&cache=false"
```

### 4. HTMX Integration
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

### 5. JavaScript/Fetch Usage
```javascript
async function searchPackages(query, options = {}) {
    const params = new URLSearchParams({
        q: query,
        format: 'json',
        limit: options.limit || 15,
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

// Usage
searchPackages('gin').then(packages => {
    packages.forEach(pkg => {
        console.log(`${pkg.path}: ${pkg.synopsis}`);
    });
});
```

## Error Handling

The API implements comprehensive error handling with graceful fallbacks:

### Network Errors
```json
{
  "success": false,
  "query": "gin",
  "provider": "pkg.go.dev",
  "count": 0,
  "results": [],
  "error": "Failed to fetch packages: network timeout",
  "timestamp": 1703123456
}
```

### Empty Results
When no packages are found, the API returns an empty results array:
```json
{
  "success": true,
  "query": "nonexistentpackage",
  "provider": "pkg.go.dev",
  "count": 0,
  "results": [],
  "cache_hit": false,
  "timestamp": 1703123456
}
```

## Caching System

The API includes an intelligent caching system:

- **Cache Key**: Based on search query
- **Cache Duration**: In-memory until server restart
- **Cache Bypass**: Use `cache=false` parameter
- **Cache Status**: Indicated in JSON responses via `cache_hit` field

### Cache Benefits
- **Performance**: Cached results return instantly
- **Rate Limiting**: Reduces load on pkg.go.dev
- **Reliability**: Cached results available during network issues

## Fallback Provider

When pkg.go.dev is unavailable or returns errors, the API automatically falls back to a curated list of popular Go packages:

```bash
curl "http://localhost:8080/fetch-packages?q=web&provider=fallback&format=json"
```

### Fallback Packages Include:
- Web frameworks (Gin, Echo, Fiber, Chi)
- Database drivers (GORM, SQLx, MongoDB)
- Utilities (Viper, Zap, Testify)
- HTTP clients and middleware

## Integration with go-ctl

### Current Usage in go-ctl
The fetch-packages API is integrated into the go-ctl project generator:

1. **Package Search**: Users can search for Go packages in real-time
2. **Package Selection**: Found packages can be added to project dependencies
3. **Project Generation**: Selected packages are included in the generated `go.mod`

### Frontend Integration
```html
<!-- Search input with HTMX -->
<input hx-get="/fetch-packages" hx-trigger="keyup changed delay:500ms">

<!-- Results container -->
<div id="search-results"></div>

<!-- Selected packages -->
<div id="selected-packages"></div>
```

## Performance Considerations

### Optimal Settings
- **Interactive Search**: `limit=10`, `cache=true`
- **Comprehensive Search**: `limit=50`, `cache=false`
- **Fast Fallback**: `provider=fallback`, `limit=20`

### Rate Limiting
- **HTML Scraping**: Respects pkg.go.dev's rate limits
- **Caching**: Reduces API calls significantly
- **Debouncing**: Recommended 500ms delay for user input

## Migration Guide

### From Legacy search-packages
**Old Usage:**
```html
<input hx-get="/search-packages">
```

**New Usage (Recommended):**
```html
<input hx-get="/fetch-packages">
```

**No Changes Required**: The old endpoint automatically redirects to the new one.

### API Response Migration
**Old JSON (if you were using custom implementation):**
```json
["package1", "package2"]
```

**New JSON:**
```json
{
  "success": true,
  "results": [
    {"path": "package1", "synopsis": "Description"},
    {"path": "package2", "synopsis": "Description"}
  ]
}
```

## Testing

### Manual Testing
```bash
# Test basic functionality
curl "http://localhost:8080/fetch-packages?q=gin&format=json"

# Test fallback
curl "http://localhost:8080/fetch-packages?q=test&provider=fallback&format=json"

# Test caching
curl "http://localhost:8080/fetch-packages?q=echo&format=json"
curl "http://localhost:8080/fetch-packages?q=echo&format=json" # Should be cached

# Test legacy endpoint
curl "http://localhost:8080/search-packages?q=fiber"
```

### Performance Testing
```bash
#!/bin/bash
queries=("gin" "echo" "fiber" "chi" "mux")

for query in "${queries[@]}"; do
    echo "Testing: $query"
    time curl -s "http://localhost:8080/fetch-packages?q=$query&format=json" > /dev/null
done
```

## Best Practices

### For Frontend Developers
1. **Debounce Input**: Use 500ms delay for search input
2. **Handle Empty States**: Gracefully handle empty results
3. **Error Handling**: Always check the `success` field in JSON responses
4. **Loading States**: Show loading indicators during searches

### For API Consumers
1. **Use JSON Format**: For programmatic access
2. **Implement Caching**: Enable caching for better performance
3. **Handle Errors**: Implement fallback strategies
4. **Respect Limits**: Use reasonable limit values

### For Production
1. **Monitor Performance**: Track response times and cache hit rates
2. **Error Monitoring**: Log and monitor API errors
3. **Rate Limiting**: Implement client-side rate limiting
4. **HTTPS**: Always use HTTPS in production

## Troubleshooting

### Common Issues

**Empty Results**
- Check query spelling
- Try fallback provider
- Verify network connectivity

**Slow Responses**
- Network issues with pkg.go.dev
- Enable caching
- Use fallback provider for testing

**Cache Not Working**
- Cache is in-memory and resets on server restart
- Use `cache=false` to test fresh results

### Debug Endpoints
```bash
# Test server health
curl "http://localhost:8080/"

# Test with fallback
curl "http://localhost:8080/fetch-packages?q=test&provider=fallback&format=json"

# Test without cache
curl "http://localhost:8080/fetch-packages?q=test&cache=false&format=json"
```

## API Schema

### Request
```
GET /fetch-packages?q={query}&format={format}&limit={limit}&provider={provider}&cache={cache}
```

### Response (JSON)
```typescript
interface PackageFetchResponse {
  success: boolean;
  query: string;
  provider: string;
  count: number;
  results: PackageResult[];
  error?: string;
  cache_hit: boolean;
  timestamp: number;
}

interface PackageResult {
  path: string;
  synopsis: string;
}
```

### Response (HTML)
Returns HTML snippet compatible with HTMX and the existing go-ctl interface.

---

## Implementation Details

### HTML Scraping Strategy
The API uses specific CSS selectors to extract package information:

- **Package Path**: `a[data-test-id='snippet-title'] .SearchSnippet-header-path`
- **Synopsis**: `.SearchSnippet-synopsis`
- **Version**: `span[data-test-id='snippet-published'] parent strong:first`

### Error Recovery
1. Network timeout → Fallback provider
2. HTML parsing error → Fallback provider  
3. Empty results → Return empty array (not error)
4. Invalid parameters → Return error with details

### Security Considerations
- Input sanitization for all query parameters
- Rate limiting to prevent abuse
- No sensitive data exposure in error messages
- HTTPS recommended for production

---

*This API is part of the go-ctl project - a dynamic Go project generator inspired by Spring Boot Initializr.*