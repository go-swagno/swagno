#!/bin/bash

# Test script for both Swagger 2.0 and OpenAPI 3.0 versions

echo "Testing Swagno v2 (Swagger 2.0) and v3 (OpenAPI 3.0) implementations..."

echo ""
echo "=== Testing v3 (OpenAPI 3.0) ==="
cd v3
echo "Running tests for OpenAPI 3.0..."
if go test -v .; then
    echo "✅ OpenAPI 3.0 tests passed"
else
    echo "❌ OpenAPI 3.0 tests failed"
    exit 1
fi

echo ""
echo "Building v3 example..."
cd example/basic
if go build -o openapi3-example .; then
    echo "✅ OpenAPI 3.0 example builds successfully"
    rm -f openapi3-example
else
    echo "❌ OpenAPI 3.0 example build failed"
    exit 1
fi

cd ../..
echo ""
echo "=== Testing main project (Swagger 2.0) ==="
cd ..
echo "Running existing tests for Swagger 2.0..."
if go test -v .; then
    echo "✅ Swagger 2.0 tests passed"
else
    echo "❌ Swagger 2.0 tests failed"
    exit 1
fi

echo ""
echo "Building v2 examples..."
cd example/http
if go build -o swagger2-example .; then
    echo "✅ Swagger 2.0 example builds successfully"
    rm -f swagger2-example
else
    echo "❌ Swagger 2.0 example build failed"
    exit 1
fi

cd ../..
echo ""
echo "=== All tests passed! ==="
echo "✅ Both Swagger 2.0 and OpenAPI 3.0 versions are working correctly"
echo ""
echo "Usage summary:"
echo "- For Swagger 2.0: import \"github.com/go-swagno/swagno\""
echo "- For OpenAPI 3.0: import v3 \"github.com/go-swagno/swagno/v3\""
