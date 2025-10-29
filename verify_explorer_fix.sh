#!/bin/bash

# Project Explorer Template Fix Verification Script
# This script tests that the Project Explorer now shows real template content

set -e

echo "🧪 Project Explorer Fix Verification"
echo "===================================="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Start server in background
echo "🚀 Starting test server..."
go build -o test_server cmd/server/*.go
./test_server > server.log 2>&1 &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Test parameters
BASE_URL="http://localhost:8080/file-content"
PROJECT_NAME="test-app"
HTTP_PACKAGE="gin"
DATABASE="postgres"
DRIVER="gorm"
PARAMS="projectName=${PROJECT_NAME}&httpPackage=${HTTP_PACKAGE}&databases=${DATABASE}&driver_${DATABASE}=${DRIVER}"

# Files that were showing fake content before
declare -a TEST_FILES=(
    "internal/storage/db.go"
    "internal/handler/handler.go"
    "internal/domain/model.go"
    "internal/config/config.go"
    "internal/service/service.go"
    "internal/storage/postgres/repository.go"
    "go.mod"
    "main.go"
)

# Test each file
PASSED=0
FAILED=0

echo ""
echo "📋 Testing file content generation..."
echo ""

for file in "${TEST_FILES[@]}"; do
    echo -n "Testing $file... "

    # Make request
    response=$(curl -s "${BASE_URL}?path=${file}&${PARAMS}")

    # Check if response contains real code indicators
    if echo "$response" | grep -q "package\|import\|func\|type\|var"; then
        # Additional checks for specific content
        case "$file" in
            "internal/handler/handler.go")
                if echo "$response" | grep -q "gin.Context\|Handler struct\|service \*service.Service"; then
                    echo -e "${GREEN}✓ PASS${NC}"
                    ((PASSED++))
                else
                    echo -e "${RED}✗ FAIL (missing expected handler content)${NC}"
                    ((FAILED++))
                fi
                ;;
            "internal/storage/db.go")
                if echo "$response" | grep -q "gorm.DB\|InitDatabase\|HealthCheck"; then
                    echo -e "${GREEN}✓ PASS${NC}"
                    ((PASSED++))
                else
                    echo -e "${RED}✗ FAIL (missing expected database content)${NC}"
                    ((FAILED++))
                fi
                ;;
            "internal/domain/model.go")
                if echo "$response" | grep -q "type.*struct\|User\|domain"; then
                    echo -e "${GREEN}✓ PASS${NC}"
                    ((PASSED++))
                else
                    echo -e "${RED}✗ FAIL (missing expected model content)${NC}"
                    ((FAILED++))
                fi
                ;;
            "internal/config/config.go")
                if echo "$response" | grep -q "Config.*struct\|Load.*func\|ServerConfig"; then
                    echo -e "${GREEN}✓ PASS${NC}"
                    ((PASSED++))
                else
                    echo -e "${RED}✗ FAIL (missing expected config content)${NC}"
                    ((FAILED++))
                fi
                ;;
            "internal/service/service.go")
                if echo "$response" | grep -q "Service.*struct\|interface\|CreateUser\|ListUsers"; then
                    echo -e "${GREEN}✓ PASS${NC}"
                    ((PASSED++))
                else
                    echo -e "${RED}✗ FAIL (missing expected service content)${NC}"
                    ((FAILED++))
                fi
                ;;
            "internal/storage/postgres/repository.go")
                if echo "$response" | grep -q "postgres\|gorm\|repository\|Storage.*struct"; then
                    echo -e "${GREEN}✓ PASS${NC}"
                    ((PASSED++))
                else
                    echo -e "${RED}✗ FAIL (missing expected repository content)${NC}"
                    ((FAILED++))
                fi
                ;;
            "go.mod")
                if echo "$response" | grep -q "module.*${PROJECT_NAME}\|go 1\.\|require"; then
                    echo -e "${GREEN}✓ PASS${NC}"
                    ((PASSED++))
                else
                    echo -e "${RED}✗ FAIL (missing expected go.mod content)${NC}"
                    ((FAILED++))
                fi
                ;;
            "main.go")
                if echo "$response" | grep -q "package main\|func main\|${HTTP_PACKAGE}"; then
                    echo -e "${GREEN}✓ PASS${NC}"
                    ((PASSED++))
                else
                    echo -e "${RED}✗ FAIL (missing expected main.go content)${NC}"
                    ((FAILED++))
                fi
                ;;
            *)
                echo -e "${GREEN}✓ PASS${NC}"
                ((PASSED++))
                ;;
        esac
    else
        echo -e "${RED}✗ FAIL (no valid code structure found)${NC}"
        echo "  Response preview: $(echo "$response" | head -c 100)..."
        ((FAILED++))
    fi
done

echo ""
echo "🧪 Testing template variable substitution..."
echo ""

# Test specific template variables
echo -n "Testing project name substitution... "
response=$(curl -s "${BASE_URL}?path=go.mod&${PARAMS}")
if echo "$response" | grep -q "module ${PROJECT_NAME}"; then
    echo -e "${GREEN}✓ PASS${NC}"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL${NC}"
    ((FAILED++))
fi

echo -n "Testing HTTP framework selection... "
response=$(curl -s "${BASE_URL}?path=internal/handler/handler.go&${PARAMS}")
if echo "$response" | grep -q "github.com/gin-gonic/gin"; then
    echo -e "${GREEN}✓ PASS${NC}"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL${NC}"
    ((FAILED++))
fi

echo -n "Testing database driver selection... "
response=$(curl -s "${BASE_URL}?path=internal/storage/db.go&${PARAMS}")
if echo "$response" | grep -q "gorm.io/gorm"; then
    echo -e "${GREEN}✓ PASS${NC}"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL${NC}"
    ((FAILED++))
fi

echo ""
echo "🧪 Testing error handling..."
echo ""

# Test non-existent file
echo -n "Testing non-existent file handling... "
response=$(curl -s "${BASE_URL}?path=non/existent/file.go&${PARAMS}")
if [ ! -z "$response" ]; then
    echo -e "${GREEN}✓ PASS (returns fallback content)${NC}"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL (empty response)${NC}"
    ((FAILED++))
fi

# Clean up
kill $SERVER_PID 2>/dev/null || true
rm -f test_server

echo ""
echo "📊 Test Results"
echo "==============="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"
echo -e "Total:  $((PASSED + FAILED))"

if [ $FAILED -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 All tests passed! Project Explorer fix is working correctly.${NC}"
    echo ""
    echo "✅ The Project Explorer now shows:"
    echo "   - Real template-generated content"
    echo "   - Framework-specific implementations"
    echo "   - Database-specific code"
    echo "   - Proper import statements"
    echo "   - Actual project structure"
    echo ""
    echo "Users will now see exactly what they'll get in the downloaded project!"
    exit 0
else
    echo ""
    echo -e "${RED}❌ Some tests failed. Please check the implementation.${NC}"
    echo ""
    echo "Common issues to check:"
    echo "   - Template files exist and are properly formatted"
    echo "   - Generator.GenerateFileContent method is working"
    echo "   - Template data structure includes all required fields"
    echo "   - handleFileContent is using the generator correctly"
    exit 1
fi
