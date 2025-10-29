# Fetch Packages API Documentation

The `fetch-packages` API provides a dynamic and flexible way to search for Go packages from various providers. This API replaces the legacy `search-packages` endpoint with enhanced functionality, better error handling, and support for multiple response formats.

## Overview

The fetch-packages API is designed to be:
- **Dynamic**: Configurable query parameters for different use cases
- **Flexible**: Support for both JSON and HTML response formats
- **Cached**: Built-in caching mechanism for improved performance
- **Fallback-ready**: Graceful degradation when external services are unavailable
- **Backward Compatible**: Works with existing HTMX frontend implementations

## Endpoints

### Primary Endpoint
```
GET /fetch-packages
```

### Legacy Endpoint (Backward Compatibility)
```
GET /search-packages
```
*Note: This endpoint redirects to `/fetch-packages` for backward compatibility*

## Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `q` | string | "" | Search query for packages |
| `provider` | string | "pkg.go.dev" | Package provider to search |
| `format` | string | "html" | Response format (`html` or `json`) |
| `limit` | integer | 15 | Maximum number of results (1-100) |
| `cache` | boolean | true | Enable/disable caching |

### Provider Options

| Provider | Description |
|----------|-------------|
| `pkg.go.dev` | Official Go package registry (default) |
| `fallback` | Local fallback with popular packages |

## Response Formats

### HTML Response (Default)
Used for HTMX integration and web interface compatibility.

**Request:**
```
GET /fetch-packages?q=gin&format=html
```

**Response:**
```html
<div class="space-y-2">
    <div class="p-3 border border-gray-200 rounded-lg hover:border-blue-300 transition-colors">
        <div class="flex justify-between items-start mb-2">
            <code class="text-sm font-mono text-blue-600">github.com/gin-gonic/gin</code>
            <button class="text-blue-500 hover:text-blue-700 text-sm font-medium" 
                    hx-post="/add-package" 
                    hx-vals='{"pkgPath": "github.com/gin-gonic/gin"}'
                    hx-target="#selected-packages" 
                    hx-swap="beforeend">
                Add
            </button>
        </div>
        <p class="text-sm text-gray-600">Gin is a HTTP web framework written in Go (Golang)</p>
    </div>
</div>
```

### JSON Response
Used for API integration and programmatic access.

**Request:**
```
GET /fetch-packages?q=gin&format=json&limit=5
```

**Response:**
```json
{
  "success": true,
  "query": "gin",
  "provider": "pkg.go.dev",
  "count": 5,
  "results": [
    {
      "path": "github.com/gin-gonic/gin",
      "synopsis": "Gin is a HTTP web framework written in Go (Golang)"
    },
    {
      "path": "github.com/gin-contrib/cors",
      "synopsis": "Official CORS gin's middleware"
    }
  ],
  "cache_hit": false,
  "timestamp": 1703123456
}
```

### Error Response (JSON)
```json
{
  "success": false,
  "query": "invalid-query",
  "provider": "pkg.go.dev",
  "count": 0,
  "results": [],
  "error": "Failed to fetch packages: network timeout",
  "timestamp": 1703123456
}
```

## Usage Examples

### Basic Search
```bash
curl "http://localhost:8080/fetch-packages?q=gin"
```

### JSON API Usage
```bash
curl "http://localhost:8080/fetch-packages?q=gin&format=json&limit=10"
```

### Custom Provider with Caching Disabled
```bash
curl "http://localhost:8080/fetch-packages?q=web&provider=fallback&cache=false"
```

### HTMX Frontend Integration
```html
<input type="text" 
       name="q"
       placeholder="Search for packages..."
       hx-get="/fetch-packages"
       hx-trigger="keyup changed delay:500ms"
       hx-target="#search-results"
       hx-swap="innerHTML">

<div id="search-results"></div>
```

## JavaScript/Fetch Usage
```javascript
async function searchPackages(query) {
  const response = await fetch(`/fetch-packages?q=${encodeURIComponent(query)}&format=json&limit=20`);
  const data = await response.json();
  
  if (data.success) {
    console.log(`Found ${data.count} packages:`, data.results);
    return data.results;
  } else {
    console.error('Search failed:', data.error);
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

## Caching Mechanism

The API includes a built-in caching system:

- **Cache Key**: Based on the search query
- **Cache Duration**: Configurable (default: in-memory until restart)
- **Cache Bypass**: Use `cache=false` parameter
- **Cache Hit Indicator**: Available in JSON response (`cache_hit` field)

## Error Handling

The API implements comprehensive error handling:

1. **Network Errors**: Automatic fallback to local package list
2. **Invalid Queries**: Returns empty results gracefully
3. **Rate Limiting**: Respects pkg.go.dev rate limits
4. **Timeout Handling**: 10-second timeout with fallback
5. **Malformed Responses**: Graceful parsing error handling

## Rate Limiting and Best Practices

### Recommended Usage Patterns

1. **Debounced Search**: Use 500ms delay for user input
2. **Reasonable Limits**: Keep limit under 50 for good performance
3. **Cache Utilization**: Enable caching for repeated queries
4. **Error Handling**: Always handle both success and error states

### Rate Limiting
- **pkg.go.dev**: Respects upstream rate limits
- **Fallback**: No rate limiting on local fallback
- **Caching**: Reduces API calls significantly

## Integration with go-ctl

### Current Integration Points

1. **Main Form**: Package search in the project generator
2. **HTMX**: Dynamic search results without page refresh  
3. **Package Selection**: Add/remove packages from project dependencies
4. **Template Generation**: Selected packages included in go.mod

### Backward Compatibility

The legacy `/search-packages` endpoint is maintained for backward compatibility:
- Automatically redirects to `/fetch-packages?format=html`
- Maintains same response format
- No breaking changes to existing frontend code

## Configuration

### Environment Variables
Currently, no specific environment variables are required. The API uses:
- Default timeout: 10 seconds
- Default cache: In-memory
- Default provider: pkg.go.dev

### Future Enhancements
- Redis-based caching
- Configurable timeouts
- Additional package providers
- Authentication for higher rate limits

## API Response Schema

### PackageFetchResponse (JSON)
```go
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
```

### PkgGoDevResult
```go
type PkgGoDevResult struct {
    Path     string `json:"path"`
    Synopsis string `json:"synopsis"`
}
```

## Troubleshooting

### Common Issues

1. **Empty Results**: Check query spelling and try fallback provider
2. **Slow Response**: Network issues, fallback will be used automatically
3. **Cache Issues**: Use `cache=false` to bypass cache
4. **Format Issues**: Ensure `format` parameter is `html` or `json`

### Debug Information

Enable debug logging to see:
- API call details
- Cache hit/miss information
- Fallback usage
- Error details

### Health Check
```bash
# Test basic functionality
curl "http://localhost:8080/fetch-packages?q=test&format=json"

# Test fallback
curl "http://localhost:8080/fetch-packages?q=test&provider=fallback&format=json"
```

## Migration Guide

### From search-packages to fetch-packages

**Old Usage:**
```html
<input hx-get="/search-packages" hx-trigger="keyup changed delay:500ms">
```

**New Usage (Recommended):**
```html
<input hx-get="/fetch-packages" hx-trigger="keyup changed delay:500ms">
```

**No changes required** - the old endpoint still works but redirects to the new one.

### Adding JSON Support

**New JSON API Usage:**
```javascript
// Replace old implementation with:
fetch('/fetch-packages?q=gin&format=json')
  .then(response => response.json())
  .then(data => {
    if (data.success) {
      // Handle results: data.results
    } else {
      // Handle error: data.error
    }
  });
```

## Performance Considerations

- **Caching**: First search hits network, subsequent searches use cache
- **Fallback**: Local fallback is faster than network calls
- **Limits**: Lower limits improve response time
- **Debouncing**: Implement client-side debouncing for better UX

## Security Considerations

- **Input Validation**: All query parameters are validated and sanitized
- **Rate Limiting**: Built-in protection against abuse
- **Error Messages**: No sensitive information exposed in errors
- **HTTPS**: Use HTTPS in production for secure communication

---

*This API is part of the go-ctl project - a web-based Go project generator inspired by Spring Boot Initializr.*