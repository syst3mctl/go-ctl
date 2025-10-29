# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
main
{{.ProjectName}}

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Environment files
.env
.env.local
.env.*.local

# Log files
*.log
logs/

# Temporary files
tmp/
temp/
.tmp/

# Documentation and markdown files (optional)
*.md
!README.md

# Git files
.git/
.gitignore
.gitattributes

# CI/CD files
.github/
.gitlab-ci.yml
Jenkinsfile

# Docker files (don't include in build context)
Dockerfile*
docker-compose*.yml
.dockerignore

# Development tools
Makefile
{{if .HasFeature "air"}}.air.toml{{end}}

# Testing files
*_test.go
testdata/
coverage.out
coverage.html

# Build artifacts
dist/
build/
bin/

# Node.js (if using frontend assets)
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Configuration files that shouldn't be in container
config/local.*
config/development.*
config/test.*
