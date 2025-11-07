package generator

import (
	"fmt"
	"strings"

	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// GenerateReactProject generates a React project structure
func (g *Generator) GenerateReactProject(config metadata.ProjectConfig) map[string]string {
	if config.FrontendConfig == nil {
		return make(map[string]string)
	}

	files := make(map[string]string)
	frontendConfig := config.FrontendConfig
	isTypeScript := frontendConfig.Language.ID == "typescript"

	// package.json
	files["package.json"] = g.generatePackageJson(config)

	// Vite config
	if isTypeScript {
		files["vite.config.ts"] = g.generateViteConfigTS(config)
	} else {
		files["vite.config.js"] = g.generateViteConfigJS(config)
	}

	// TypeScript config
	if isTypeScript {
		files["tsconfig.json"] = g.generateTSConfig(config)
		files["tsconfig.node.json"] = g.generateTSConfigNode(config)
		files["src/vite-env.d.ts"] = g.generateViteEnvDTS(config)
	}

	// ESLint config
	if frontendConfig.Linter.ID == "eslint" {
		if isTypeScript {
			files[".eslintrc.cjs"] = g.generateESLintConfigTS(config)
		} else {
			files[".eslintrc.cjs"] = g.generateESLintConfigJS(config)
		}
	}

	// Prettier config
	hasPrettier := false
	for _, feature := range frontendConfig.Features {
		if feature.ID == "prettier" {
			hasPrettier = true
			break
		}
	}
	if hasPrettier {
		files[".prettierrc"] = g.generatePrettierConfig(config)
		files[".prettierignore"] = g.generatePrettierIgnore(config)
	}

	// Tailwind config
	hasTailwind := false
	for _, feature := range frontendConfig.Features {
		if feature.ID == "tailwind" {
			hasTailwind = true
			break
		}
	}
	if hasTailwind {
		files["tailwind.config.js"] = g.generateTailwindConfig(config)
		files["postcss.config.js"] = g.generatePostCSSConfig(config)
	}

	// index.html
	files["index.html"] = g.generateIndexHtml(config)

	// Main entry point
	if isTypeScript {
		files["src/main.tsx"] = g.generateMainTSX(config)
	} else {
		files["src/main.jsx"] = g.generateMainJSX(config)
	}

	// App component
	if isTypeScript {
		files["src/App.tsx"] = g.generateAppTSX(config)
	} else {
		files["src/App.jsx"] = g.generateAppJSX(config)
	}

	// Example component
	if isTypeScript {
		files["src/components/HelloWorld.tsx"] = g.generateHelloWorldTSX(config)
	} else {
		files["src/components/HelloWorld.jsx"] = g.generateHelloWorldJSX(config)
	}

	// Create directory structure with placeholder files
	// Assets directory
	files["src/assets/.gitkeep"] = ""
	
	// Lib directory structure
	files["src/lib/hooks/.gitkeep"] = ""
	files["src/lib/types/.gitkeep"] = ""
	files["src/lib/validations/.gitkeep"] = ""
	
	// Store directory
	files["src/store/.gitkeep"] = ""
	
	// Generate example files for lib structure if TypeScript
	if isTypeScript {
		files["src/lib/types/index.ts"] = g.generateTypesIndexTS(config)
		files["src/lib/validations/index.ts"] = g.generateValidationsIndexTS(config)
		files["src/lib/hooks/index.ts"] = g.generateHooksIndexTS(config)
		files["src/store/index.tsx"] = g.generateStoreIndexTSX(config)
	} else {
		files["src/lib/types/index.js"] = g.generateTypesIndexJS(config)
		files["src/lib/validations/index.js"] = g.generateValidationsIndexJS(config)
		files["src/lib/hooks/index.js"] = g.generateHooksIndexJS(config)
		files["src/store/index.jsx"] = g.generateStoreIndexJSX(config)
	}

	// Styles
	files["src/styles/index.css"] = g.generateIndexCSS(config, hasTailwind)

	// .gitignore
	files[".gitignore"] = g.generateFrontendGitignore(config)

	// README
	files["README.md"] = g.generateReactREADME(config)

	return files
}

// generatePackageJson generates package.json
func (g *Generator) generatePackageJson(config metadata.ProjectConfig) string {
	frontendConfig := config.FrontendConfig
	isTypeScript := frontendConfig.Language.ID == "typescript"

	dependencies := map[string]string{
		"react":       "^18.2.0",
		"react-dom":   "^18.2.0",
		"vite":        "^5.0.0",
		"@vitejs/plugin-react": "^4.2.0",
	}

	devDependencies := map[string]string{
		"@types/react":     "^18.2.0",
		"@types/react-dom": "^18.2.0",
	}

	// Add feature dependencies
	for _, feature := range frontendConfig.Features {
		switch feature.ID {
		case "react-router":
			dependencies["react-router-dom"] = "^6.20.0"
		case "react-query":
			dependencies["@tanstack/react-query"] = "^5.8.0"
		case "axios":
			dependencies["axios"] = "^1.6.0"
		case "vitest":
			devDependencies["vitest"] = "^1.0.0"
			devDependencies["@testing-library/react"] = "^14.1.0"
			devDependencies["@testing-library/jest-dom"] = "^6.1.0"
		}
	}

	// Add linter
	if frontendConfig.Linter.ID == "eslint" {
		devDependencies["eslint"] = "^8.54.0"
		devDependencies["eslint-plugin-react"] = "^7.33.2"
		devDependencies["eslint-plugin-react-hooks"] = "^4.6.0"
		devDependencies["eslint-plugin-react-refresh"] = "^0.4.4"
		if isTypeScript {
			devDependencies["@typescript-eslint/eslint-plugin"] = "^6.12.0"
			devDependencies["@typescript-eslint/parser"] = "^6.12.0"
		}
	}

	// Add Prettier
	hasPrettier := false
	for _, feature := range frontendConfig.Features {
		if feature.ID == "prettier" {
			hasPrettier = true
			break
		}
	}
	if hasPrettier {
		devDependencies["prettier"] = "^3.1.0"
		devDependencies["eslint-config-prettier"] = "^9.0.0"
	}

	// Add Tailwind
	hasTailwind := false
	for _, feature := range frontendConfig.Features {
		if feature.ID == "tailwind" {
			hasTailwind = true
			break
		}
	}
	if hasTailwind {
		devDependencies["tailwindcss"] = "^3.3.6"
		devDependencies["postcss"] = "^8.4.32"
		devDependencies["autoprefixer"] = "^10.4.16"
	}

	// Add custom npm packages
	for _, pkg := range frontendConfig.CustomPackages {
		// Extract package name (handle scoped packages)
		pkgName := pkg
		if strings.Contains(pkg, "@") {
			parts := strings.Split(pkg, "@")
			if len(parts) >= 2 {
				pkgName = parts[0] + "@" + parts[1]
			}
		}
		dependencies[pkgName] = "latest"
	}

	// Build dependencies string
	depsStr := ""
	for pkg, version := range dependencies {
		depsStr += fmt.Sprintf(`    "%s": "%s",
`, pkg, version)
	}

	devDepsStr := ""
	for pkg, version := range devDependencies {
		devDepsStr += fmt.Sprintf(`    "%s": "%s",
`, pkg, version)
	}

	// Remove trailing comma
	if len(depsStr) > 0 {
		depsStr = depsStr[:len(depsStr)-2] + "\n"
	}
	if len(devDepsStr) > 0 {
		devDepsStr = devDepsStr[:len(devDepsStr)-2] + "\n"
	}

	return fmt.Sprintf(`{
  "name": "%s",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "lint": "eslint . --ext js,jsx%s --report-unused-disable-directives --max-warnings 0",
    "preview": "vite preview"%s
  },
  "dependencies": {
%s  },
  "devDependencies": {
%s  }
}
`, config.ProjectName, func() string {
		if isTypeScript {
			return ",ts,tsx"
		}
		return ""
	}(), func() string {
		hasVitest := false
		for _, feature := range frontendConfig.Features {
			if feature.ID == "vitest" {
				hasVitest = true
				break
			}
		}
		if hasVitest {
			return `,
    "test": "vitest"`
		}
		return ""
	}(), depsStr, devDepsStr)
}

// Helper functions for generating individual files
func (g *Generator) generateViteConfigJS(config metadata.ProjectConfig) string {
	return `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
})
`
}

func (g *Generator) generateViteConfigTS(config metadata.ProjectConfig) string {
	return `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
})
`
}

func (g *Generator) generateTSConfig(config metadata.ProjectConfig) string {
	return `{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,

    /* Bundler mode */
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",

    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
`
}

func (g *Generator) generateTSConfigNode(config metadata.ProjectConfig) string {
	return `{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true
  },
  "include": ["vite.config.ts"]
}
`
}

func (g *Generator) generateViteEnvDTS(config metadata.ProjectConfig) string {
	return `/// <reference types="vite/client" />
`
}

func (g *Generator) generateESLintConfigJS(config metadata.ProjectConfig) string {
	return `module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    'eslint:recommended',
    'plugin:react/recommended',
    'plugin:react/jsx-runtime',
    'plugin:react-hooks/recommended',
  ],
  ignorePatterns: ['dist', '.eslintrc.cjs'],
  parserOptions: { ecmaVersion: 'latest', sourceType: 'module' },
  settings: { react: { version: '18.2' } },
  plugins: ['react-refresh'],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
  },
}
`
}

func (g *Generator) generateESLintConfigTS(config metadata.ProjectConfig) string {
	return `module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
  ],
  ignorePatterns: ['dist', '.eslintrc.cjs'],
  parser: '@typescript-eslint/parser',
  plugins: ['react-refresh'],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
  },
}
`
}

func (g *Generator) generatePrettierConfig(config metadata.ProjectConfig) string {
	return `{
  "semi": true,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
`
}

func (g *Generator) generatePrettierIgnore(config metadata.ProjectConfig) string {
	return `dist
node_modules
build
`
}

func (g *Generator) generateTailwindConfig(config metadata.ProjectConfig) string {
	return `/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
`
}

func (g *Generator) generatePostCSSConfig(config metadata.ProjectConfig) string {
	return `export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
`
}

func (g *Generator) generateIndexHtml(config metadata.ProjectConfig) string {
	frontendConfig := config.FrontendConfig
	isTypeScript := frontendConfig.Language.ID == "typescript"
	ext := "jsx"
	if isTypeScript {
		ext = "tsx"
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>%s</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.%s"></script>
  </body>
</html>
`, config.ProjectName, ext)
}

func (g *Generator) generateMainJSX(config metadata.ProjectConfig) string {
	hasReactQuery := false
	for _, feature := range config.FrontendConfig.Features {
		if feature.ID == "react-query" {
			hasReactQuery = true
			break
		}
	}

	imports := "import React from 'react'\nimport ReactDOM from 'react-dom/client'\nimport App from './App.jsx'\nimport './styles/index.css'\n"
	
	var providerWraps string
	if hasReactQuery {
		imports = "import React from 'react'\nimport ReactDOM from 'react-dom/client'\nimport { QueryClient, QueryClientProvider } from '@tanstack/react-query'\nimport App from './App.jsx'\nimport './styles/index.css'\n"
		setup := "const queryClient = new QueryClient()\n\n"
		providerWraps = "<QueryClientProvider client={queryClient}>\n      <App />\n    </QueryClientProvider>"
		return imports + "\n" + setup + "ReactDOM.createRoot(document.getElementById('root')).render(\n  <React.StrictMode>\n    " + providerWraps + "\n  </React.StrictMode>,\n)"
	} else {
		// Use Context API from store
		imports += "import { AppProvider } from './store/index.jsx'\n"
		providerWraps = "<AppProvider>\n      <App />\n    </AppProvider>"
	}

	root := "ReactDOM.createRoot(document.getElementById('root')).render(\n  <React.StrictMode>\n    " + providerWraps + "\n  </React.StrictMode>,\n)"

	return imports + "\n" + root
}

func (g *Generator) generateMainTSX(config metadata.ProjectConfig) string {
	hasReactQuery := false
	for _, feature := range config.FrontendConfig.Features {
		if feature.ID == "react-query" {
			hasReactQuery = true
			break
		}
	}

	imports := "import React from 'react'\nimport ReactDOM from 'react-dom/client'\nimport App from './App.tsx'\nimport './styles/index.css'\n"
	
	var providerWraps string
	if hasReactQuery {
		imports = "import React from 'react'\nimport ReactDOM from 'react-dom/client'\nimport { QueryClient, QueryClientProvider } from '@tanstack/react-query'\nimport App from './App.tsx'\nimport './styles/index.css'\n"
		setup := "const queryClient = new QueryClient()\n\n"
		providerWraps = "<QueryClientProvider client={queryClient}>\n      <App />\n    </QueryClientProvider>"
		return imports + "\n" + setup + "ReactDOM.createRoot(document.getElementById('root')!).render(\n  <React.StrictMode>\n    " + providerWraps + "\n  </React.StrictMode>,\n)"
	} else {
		// Use Context API from store
		imports += "import { AppProvider } from './store/index.tsx'\n"
		providerWraps = "<AppProvider>\n      <App />\n    </AppProvider>"
	}

	root := "ReactDOM.createRoot(document.getElementById('root')!).render(\n  <React.StrictMode>\n    " + providerWraps + "\n  </React.StrictMode>,\n)"

	return imports + "\n" + root
}

func (g *Generator) generateAppJSX(config metadata.ProjectConfig) string {
	hasReactRouter := false
	for _, feature := range config.FrontendConfig.Features {
		if feature.ID == "react-router" {
			hasReactRouter = true
			break
		}
	}

	if hasReactRouter {
		return `import { BrowserRouter, Routes, Route } from 'react-router-dom'
import HelloWorld from './components/HelloWorld'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HelloWorld />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
`
	}

	return `import HelloWorld from './components/HelloWorld'

function App() {
  return (
    <div>
      <HelloWorld />
    </div>
  )
}

export default App
`
}

func (g *Generator) generateAppTSX(config metadata.ProjectConfig) string {
	hasReactRouter := false
	for _, feature := range config.FrontendConfig.Features {
		if feature.ID == "react-router" {
			hasReactRouter = true
			break
		}
	}

	if hasReactRouter {
		return `import { BrowserRouter, Routes, Route } from 'react-router-dom'
import HelloWorld from './components/HelloWorld'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HelloWorld />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
`
	}

	return `import HelloWorld from './components/HelloWorld'

function App() {
  return (
    <div>
      <HelloWorld />
    </div>
  )
}

export default App
`
}

func (g *Generator) generateHelloWorldJSX(config metadata.ProjectConfig) string {
	hasTailwind := false
	for _, feature := range config.FrontendConfig.Features {
		if feature.ID == "tailwind" {
			hasTailwind = true
			break
		}
	}

	if hasTailwind {
		return `function HelloWorld() {
  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white p-8 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold text-gray-800 mb-4">Hello, World!</h1>
        <p className="text-gray-600">Welcome to your new React app with Vite and Tailwind CSS.</p>
      </div>
    </div>
  )
}

export default HelloWorld
`
	}

	return `function HelloWorld() {
  return (
    <div style={{ padding: '2rem', textAlign: 'center' }}>
      <h1>Hello, World!</h1>
      <p>Welcome to your new React app with Vite.</p>
    </div>
  )
}

export default HelloWorld
`
}

func (g *Generator) generateHelloWorldTSX(config metadata.ProjectConfig) string {
	hasTailwind := false
	for _, feature := range config.FrontendConfig.Features {
		if feature.ID == "tailwind" {
			hasTailwind = true
			break
		}
	}

	if hasTailwind {
		return `function HelloWorld() {
  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white p-8 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold text-gray-800 mb-4">Hello, World!</h1>
        <p className="text-gray-600">Welcome to your new React app with Vite and Tailwind CSS.</p>
      </div>
    </div>
  )
}

export default HelloWorld
`
	}

	return `function HelloWorld() {
  return (
    <div style={{ padding: '2rem', textAlign: 'center' }}>
      <h1>Hello, World!</h1>
      <p>Welcome to your new React app with Vite.</p>
    </div>
  )
}

export default HelloWorld
`
}

func (g *Generator) generateIndexCSS(config metadata.ProjectConfig, hasTailwind bool) string {
	if hasTailwind {
		return `@tailwind base;
@tailwind components;
@tailwind utilities;
`
	}

	return `:root {
  font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
  line-height: 1.5;
  font-weight: 400;

  color-scheme: light dark;
  color: rgba(255, 255, 255, 0.87);
  background-color: #242424;

  font-synthesis: none;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

body {
  margin: 0;
  display: flex;
  place-items: center;
  min-width: 320px;
  min-height: 100vh;
}

#root {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
  text-align: center;
}
`
}

func (g *Generator) generateFrontendGitignore(config metadata.ProjectConfig) string {
	return `# Logs
logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*
lerna-debug.log*

node_modules
dist
dist-ssr
*.local

# Editor directories and files
.vscode/*
!.vscode/extensions.json
.idea
.DS_Store
*.suo
*.ntvs*
*.njsproj
*.sln
*.sw?
`
}

func (g *Generator) generateReactREADME(config metadata.ProjectConfig) string {
	frontendConfig := config.FrontendConfig
	isTypeScript := frontendConfig.Language.ID == "typescript"
	lang := "JavaScript"
	if isTypeScript {
		lang = "TypeScript"
	}

	features := []string{}
	for _, feature := range frontendConfig.Features {
		features = append(features, feature.Name)
	}

	featuresStr := "None"
	if len(features) > 0 {
		featuresStr = strings.Join(features, ", ")
	}

	return fmt.Sprintf(`# %s

This is a React project generated using [go-ctl](https://github.com/syst3mctl/go-ctl).

## 🚀 Getting Started

### Prerequisites

- Node.js 18+ and npm

### Installation

1. Install dependencies:
   ` + "```bash\n   npm install\n   ```" + `

2. Start the development server:
   ` + "```bash\n   npm run dev\n   ```" + `

The app will be available at http://localhost:5173

## 📚 Project Details

- **Language**: %s
- **Build Tool**: Vite
- **Linter**: ESLint
- **Features**: %s

## 🛠️ Available Scripts

- ` + "`npm run dev`" + ` - Start development server
- ` + "`npm run build`" + ` - Build for production
- ` + "`npm run preview`" + ` - Preview production build
- ` + "`npm run lint`" + ` - Run ESLint%s

## 📦 Dependencies

This project includes the following dependencies:

- React 18
- Vite
%s

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

---

**Generated with ❤️ by go-ctl**
`, config.ProjectName, lang, featuresStr, func() string {
		hasVitest := false
		for _, feature := range frontendConfig.Features {
			if feature.ID == "vitest" {
				hasVitest = true
				break
			}
		}
		if hasVitest {
			return "\n- `npm test` - Run tests with Vitest"
		}
		return ""
	}(), func() string {
		customPkgs := ""
		if len(frontendConfig.CustomPackages) > 0 {
			customPkgs = "\n### Custom Packages:\n"
			for _, pkg := range frontendConfig.CustomPackages {
				customPkgs += fmt.Sprintf("- %s\n", pkg)
			}
		}
		return customPkgs
	}())
}

// generateTypesIndexTS generates TypeScript types index file
func (g *Generator) generateTypesIndexTS(config metadata.ProjectConfig) string {
	return `// Type definitions for the application

export interface User {
  id: string
  name: string
  email: string
}

// Add your type definitions here
`
}

// generateTypesIndexJS generates JavaScript types index file (JSDoc)
func (g *Generator) generateTypesIndexJS(config metadata.ProjectConfig) string {
	return `// Type definitions using JSDoc

/**
 * @typedef {Object} User
 * @property {string} id
 * @property {string} name
 * @property {string} email
 */

// Add your type definitions here
`
}

// generateValidationsIndexTS generates TypeScript validations index file
func (g *Generator) generateValidationsIndexTS(config metadata.ProjectConfig) string {
	return `// Form validation utilities

/**
 * Validates email format
 */
export function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Validates required field
 */
export function validateRequired(value: string | undefined | null): boolean {
  return value !== undefined && value !== null && value.trim() !== ''
}

// Add your validation functions here
`
}

// generateValidationsIndexJS generates JavaScript validations index file
func (g *Generator) generateValidationsIndexJS(config metadata.ProjectConfig) string {
	return `// Form validation utilities

/**
 * Validates email format
 * @param {string} email
 * @returns {boolean}
 */
export function validateEmail(email) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Validates required field
 * @param {string|undefined|null} value
 * @returns {boolean}
 */
export function validateRequired(value) {
  return value !== undefined && value !== null && value.trim() !== ''
}

// Add your validation functions here
`
}

// generateHooksIndexTS generates TypeScript hooks index file
func (g *Generator) generateHooksIndexTS(config metadata.ProjectConfig) string {
	return `// Custom React hooks

import { useState, useEffect } from 'react'

/**
 * Custom hook example
 */
export function useExample() {
  const [value, setValue] = useState<string>('')

  useEffect(() => {
    // Hook logic here
  }, [])

  return { value, setValue }
}

// Add your custom hooks here
`
}

// generateHooksIndexJS generates JavaScript hooks index file
func (g *Generator) generateHooksIndexJS(config metadata.ProjectConfig) string {
	return `// Custom React hooks

import { useState, useEffect } from 'react'

/**
 * Custom hook example
 */
export function useExample() {
  const [value, setValue] = useState('')

  useEffect(() => {
    // Hook logic here
  }, [])

  return { value, setValue }
}

// Add your custom hooks here
`
}

// generateStoreIndexTSX generates TypeScript store index file (Context API)
func (g *Generator) generateStoreIndexTSX(config metadata.ProjectConfig) string {
	hasReactQuery := false
	for _, feature := range config.FrontendConfig.Features {
		if feature.ID == "react-query" {
			hasReactQuery = true
			break
		}
	}

	if hasReactQuery {
		return `// Store/State management using React Query
// You can also use Context API or Redux here

import { QueryClient } from '@tanstack/react-query'

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

// Add your store/state management logic here
`
	}

	return `// Store/State management using Context API
// You can also use Redux or Zustand here

import { createContext, useContext, useState, ReactNode } from 'react'

interface AppState {
  // Add your state properties here
}

interface AppContextType {
  state: AppState
  setState: (state: AppState) => void
}

const AppContext = createContext<AppContextType | undefined>(undefined)

export function AppProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AppState>({})

  return (
    <AppContext.Provider value={{ state, setState }}>
      {children}
    </AppContext.Provider>
  )
}

export function useAppContext() {
  const context = useContext(AppContext)
  if (context === undefined) {
    throw new Error('useAppContext must be used within AppProvider')
  }
  return context
}
`
}

// generateStoreIndexJSX generates JavaScript store index file (Context API)
func (g *Generator) generateStoreIndexJSX(config metadata.ProjectConfig) string {
	hasReactQuery := false
	for _, feature := range config.FrontendConfig.Features {
		if feature.ID == "react-query" {
			hasReactQuery = true
			break
		}
	}

	if hasReactQuery {
		return `// Store/State management using React Query
// You can also use Context API or Redux here

import { QueryClient } from '@tanstack/react-query'

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

// Add your store/state management logic here
`
	}

	return `// Store/State management using Context API
// You can also use Redux or Zustand here

import { createContext, useContext, useState } from 'react'

const AppContext = createContext(undefined)

export function AppProvider({ children }) {
  const [state, setState] = useState({})

  return (
    <AppContext.Provider value={{ state, setState }}>
      {children}
    </AppContext.Provider>
  )
}

export function useAppContext() {
  const context = useContext(AppContext)
  if (context === undefined) {
    throw new Error('useAppContext must be used within AppProvider')
  }
  return context
}
`
}

