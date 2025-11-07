package generator

import (
	"fmt"
	"strings"

	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// GenerateReactProject generates a frontend project structure (supports React, Angular, Solid JS, Vue, Svelte)
func (g *Generator) GenerateReactProject(config metadata.ProjectConfig) map[string]string {
	if config.FrontendConfig == nil {
		return make(map[string]string)
	}

	files := make(map[string]string)
	frontendConfig := config.FrontendConfig
	frameworkID := frontendConfig.Framework.ID
	if frameworkID == "" {
		frameworkID = "react" // Default to React
	}
	isTypeScript := frontendConfig.Language.ID == "typescript"

	// Generate framework-specific project structure
	switch frameworkID {
	case "react":
		return g.generateReactProject(files, config, isTypeScript)
	case "angular":
		return g.generateAngularProject(files, config, isTypeScript)
	case "solid":
		return g.generateSolidProject(files, config, isTypeScript)
	case "vue":
		return g.generateVueProject(files, config, isTypeScript)
	case "svelte":
		return g.generateSvelteProject(files, config, isTypeScript)
	default:
		// Default to React if framework is unknown
		return g.generateReactProject(files, config, isTypeScript)
	}
}

// generateReactProject generates a React project structure
func (g *Generator) generateReactProject(files map[string]string, config metadata.ProjectConfig, isTypeScript bool) map[string]string {
	frontendConfig := config.FrontendConfig

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


// generateAngularProject generates an Angular project structure
func (g *Generator) generateAngularProject(files map[string]string, config metadata.ProjectConfig, isTypeScript bool) map[string]string {
	frontendConfig := config.FrontendConfig

	// package.json
	files["package.json"] = g.generateAngularPackageJson(config)

	// Angular config
	files["angular.json"] = g.generateAngularConfig(config)

	// TypeScript configs
	files["tsconfig.json"] = g.generateAngularTSConfig(config)
	files["tsconfig.app.json"] = g.generateAngularTSConfigApp(config)
	files["tsconfig.spec.json"] = g.generateAngularTSConfigSpec(config)

	// ESLint config
	if frontendConfig.Linter.ID == "eslint" {
		files[".eslintrc.json"] = g.generateAngularESLintConfig(config)
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

	// Main files
	files["src/index.html"] = g.generateAngularIndexHTML(config)
	files["src/main.ts"] = g.generateAngularMainTS(config)
	files["src/app/app.config.ts"] = g.generateAngularAppConfig(config)
	files["src/app/app.component.ts"] = g.generateAngularAppComponent(config)
	files["src/app/app.component.html"] = g.generateAngularAppTemplate(config)
	files["src/app/app.component.css"] = g.generateAngularAppStyles(config, hasTailwind)

	// Example component
	files["src/app/components/hello-world.component.ts"] = g.generateAngularHelloWorldComponent(config)
	files["src/app/components/hello-world.component.html"] = g.generateAngularHelloWorldTemplate(config)
	files["src/app/components/hello-world.component.css"] = g.generateAngularHelloWorldStyles(config)

	// Styles
	files["src/styles.css"] = g.generateAngularStyles(config, hasTailwind)

	// .gitignore
	files[".gitignore"] = g.generateFrontendGitignore(config)

	// README
	files["README.md"] = g.generateAngularREADME(config)

	return files
}

// generateAngularPackageJson generates package.json for Angular
func (g *Generator) generateAngularPackageJson(config metadata.ProjectConfig) string {
	frontendConfig := config.FrontendConfig

	dependencies := map[string]string{
		"@angular/core":                  "^17.0.0",
		"@angular/platform-browser":     "^17.0.0",
		"@angular/common":                "^17.0.0",
		"@angular/forms":                "^17.0.0",
		"rxjs":                           "^7.8.1",
		"tslib":                          "^2.6.2",
		"zone.js":                        "^0.14.2",
	}

	devDependencies := map[string]string{
		"@angular-devkit/build-angular": "^17.0.0",
		"@angular/cli":                  "^17.0.0",
		"@angular/compiler-cli":          "^17.0.0",
		"@types/node":                    "^20.10.0",
		"typescript":                     "~5.2.2",
	}

	// Add feature dependencies
	for _, feature := range frontendConfig.Features {
		switch feature.ID {
		case "axios":
			dependencies["axios"] = "^1.6.0"
		}
	}

	// Add linter
	if frontendConfig.Linter.ID == "eslint" {
		devDependencies["@angular-eslint/builder"] = "^17.0.0"
		devDependencies["@angular-eslint/eslint-plugin"] = "^17.0.0"
		devDependencies["@angular-eslint/template-parser"] = "^17.0.0"
		devDependencies["@typescript-eslint/eslint-plugin"] = "^6.12.0"
		devDependencies["@typescript-eslint/parser"] = "^6.12.0"
		devDependencies["eslint"] = "^8.54.0"
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
  "version": "0.0.0",
  "scripts": {
    "ng": "ng",
    "start": "ng serve",
    "build": "ng build",
    "watch": "ng build --watch --configuration development",
    "test": "ng test"%s
  },
  "private": true,
  "dependencies": {
%s  },
  "devDependencies": {
%s  }
}
`, config.ProjectName, func() string {
		if frontendConfig.Linter.ID == "eslint" {
			return `,
    "lint": "ng lint"`
		}
		return ""
	}(), depsStr, devDepsStr)
}

// generateAngularConfig generates angular.json
func (g *Generator) generateAngularConfig(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`{
  "$schema": "./node_modules/@angular/cli/lib/config/schema.json",
  "version": 1,
  "newProjectRoot": "projects",
  "projects": {
    "%s": {
      "projectType": "application",
      "schematics": {
        "@schematics/angular:component": {
          "style": "css",
          "skipTests": false
        }
      },
      "root": "",
      "sourceRoot": "src",
      "prefix": "app",
      "architect": {
        "build": {
          "builder": "@angular-devkit/build-angular:browser",
          "options": {
            "outputPath": "dist/%s",
            "index": "src/index.html",
            "main": "src/main.ts",
            "polyfills": [
              "zone.js"
            ],
            "tsConfig": "tsconfig.app.json",
            "assets": [
              "src/favicon.ico",
              "src/assets"
            ],
            "styles": [
              "src/styles.css"
            ],
            "scripts": []
          },
          "configurations": {
            "production": {
              "budgets": [
                {
                  "type": "initial",
                  "maximumWarning": "500kb",
                  "maximumError": "1mb"
                },
                {
                  "type": "anyComponentStyle",
                  "maximumWarning": "2kb",
                  "maximumError": "4kb"
                }
              ],
              "outputHashing": "all"
            },
            "development": {
              "buildOptimizer": false,
              "optimization": false,
              "vendorChunk": true,
              "extractLicenses": false,
              "sourceMap": true,
              "namedChunks": true
            }
          },
          "defaultConfiguration": "production"
        },
        "serve": {
          "builder": "@angular-devkit/build-angular:dev-server",
          "configurations": {
            "production": {
              "buildTarget": "%s:build:production"
            },
            "development": {
              "buildTarget": "%s:build:development"
            }
          },
          "defaultConfiguration": "development"
        }
      }
    }
  }
}
`, config.ProjectName, config.ProjectName, config.ProjectName, config.ProjectName)
}

// generateAngularTSConfig generates tsconfig.json for Angular
func (g *Generator) generateAngularTSConfig(config metadata.ProjectConfig) string {
	return `{
  "compileOnSave": false,
  "compilerOptions": {
    "outDir": "./dist/out-tsc",
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "noImplicitOverride": true,
    "noPropertyAccessFromIndexSignature": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "esModuleInterop": true,
    "sourceMap": true,
    "declaration": false,
    "experimentalDecorators": true,
    "moduleResolution": "bundler",
    "importHelpers": true,
    "target": "ES2022",
    "module": "ES2022",
    "lib": [
      "ES2022",
      "dom"
    ],
    "skipLibCheck": true
  },
  "angularCompilerOptions": {
    "enableI18nLegacyMessageIdFormat": false,
    "strictInjectionParameters": true,
    "strictInputAccessModifiers": true,
    "strictTemplates": true
  }
}
`
}

// generateAngularTSConfigApp generates tsconfig.app.json
func (g *Generator) generateAngularTSConfigApp(config metadata.ProjectConfig) string {
	return `{
  "extends": "./tsconfig.json",
  "compilerOptions": {
    "outDir": "./out-tsc/app",
    "types": []
  },
  "files": [
    "src/main.ts"
  ],
  "include": [
    "src/**/*.d.ts"
  ]
}
`
}

// generateAngularTSConfigSpec generates tsconfig.spec.json
func (g *Generator) generateAngularTSConfigSpec(config metadata.ProjectConfig) string {
	return `{
  "extends": "./tsconfig.json",
  "compilerOptions": {
    "outDir": "./out-tsc/spec",
    "types": [
      "jasmine"
    ]
  },
  "include": [
    "src/**/*.spec.ts",
    "src/**/*.d.ts"
  ]
}
`
}

// generateAngularESLintConfig generates .eslintrc.json for Angular
func (g *Generator) generateAngularESLintConfig(config metadata.ProjectConfig) string {
	return `{
  "root": true,
  "ignorePatterns": [
    "projects/**/*"
  ],
  "overrides": [
    {
      "files": [
        "*.ts"
      ],
      "extends": [
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended",
        "plugin:@angular-eslint/recommended",
        "plugin:@angular-eslint/template/process-inline-templates"
      ],
      "rules": {}
    },
    {
      "files": [
        "*.html"
      ],
      "extends": [
        "plugin:@angular-eslint/template/recommended",
        "plugin:@angular-eslint/template/accessibility"
      ],
      "rules": {}
    }
  ]
}
`
}

// generateAngularMainTS generates src/main.ts
func (g *Generator) generateAngularMainTS(config metadata.ProjectConfig) string {
	return `import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

bootstrapApplication(AppComponent, appConfig)
  .catch((err) => console.error(err));
`
}

// generateAngularAppConfig generates src/app/app.config.ts
func (g *Generator) generateAngularAppConfig(config metadata.ProjectConfig) string {
	return `import { ApplicationConfig } from '@angular/core';

export const appConfig: ApplicationConfig = {
  providers: []
};
`
}

// generateAngularAppComponent generates src/app/app.component.ts
func (g *Generator) generateAngularAppComponent(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HelloWorldComponent } from './components/hello-world.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, HelloWorldComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = '%s';
}
`, config.ProjectName)
}

// generateAngularAppTemplate generates src/app/app.component.html
func (g *Generator) generateAngularAppTemplate(config metadata.ProjectConfig) string {
	return `<div class="app-container">
  <header>
    <h1>{{ title }}</h1>
    <p>Welcome to your Angular application!</p>
  </header>
  <main>
    <app-hello-world></app-hello-world>
  </main>
</div>
`
}

// generateAngularAppStyles generates src/app/app.component.css
func (g *Generator) generateAngularAppStyles(config metadata.ProjectConfig, hasTailwind bool) string {
	if hasTailwind {
		return `.app-container {
  @apply min-h-screen p-8;
}

header {
  @apply text-center mb-8;
}

h1 {
  @apply text-4xl font-bold mb-4;
}
`
	}
	return `.app-container {
  min-height: 100vh;
  padding: 2rem;
}

header {
  text-align: center;
  margin-bottom: 2rem;
}

h1 {
  font-size: 2.5rem;
  font-weight: bold;
  margin-bottom: 1rem;
}
`
}

// generateAngularHelloWorldComponent generates hello-world.component.ts
func (g *Generator) generateAngularHelloWorldComponent(config metadata.ProjectConfig) string {
	return `import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-hello-world',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './hello-world.component.html',
  styleUrl: './hello-world.component.css'
})
export class HelloWorldComponent {
  message = 'Hello, World!';
}
`
}

// generateAngularHelloWorldTemplate generates hello-world.component.html
func (g *Generator) generateAngularHelloWorldTemplate(config metadata.ProjectConfig) string {
	return `<div class="hello-world">
  <p>{{ message }}</p>
</div>
`
}

// generateAngularHelloWorldStyles generates hello-world.component.css
func (g *Generator) generateAngularHelloWorldStyles(config metadata.ProjectConfig) string {
	return `.hello-world {
  padding: 1rem;
  text-align: center;
}

.hello-world p {
  font-size: 1.25rem;
  color: #333;
}
`
}

// generateAngularStyles generates src/styles.css
func (g *Generator) generateAngularStyles(config metadata.ProjectConfig, hasTailwind bool) string {
	if hasTailwind {
		return `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}
`
	}
	return `/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  line-height: 1.6;
  color: #333;
}
`
}

// generateAngularIndexHTML generates src/index.html
func (g *Generator) generateAngularIndexHTML(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>%s</title>
  <base href="/">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="icon" type="image/x-icon" href="favicon.ico">
</head>
<body>
  <app-root></app-root>
</body>
</html>
`, config.ProjectName)
}

// generateAngularREADME generates README.md for Angular
func (g *Generator) generateAngularREADME(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`# %s

This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 17.0.0.

## Development server

Run `+"`ng serve`"+` for a dev server. Navigate to `+"`http://localhost:4200/`"+`. The application will automatically reload if you change any of the source files.

## Code scaffolding

Run `+"`ng generate component component-name`"+` to generate a new component. You can also use `+"`ng generate directive|pipe|service|class|guard|interface|enum|module`"+`.

## Build

Run `+"`ng build`"+` to build the project. The build artifacts will be stored in the `+"`dist/`"+` directory.

## Running unit tests

Run `+"`ng test`"+` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Further help

To get more help on the Angular CLI use `+"`ng help`"+` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.

## Generated with ❤️ by go-ctl
`, config.ProjectName)
}

// generateSolidProject generates a Solid JS project structure
func (g *Generator) generateSolidProject(files map[string]string, config metadata.ProjectConfig, isTypeScript bool) map[string]string {
	// Use React generator as base for now
	// TODO: Implement full Solid JS project structure
	return g.generateReactProject(files, config, isTypeScript)
}

// generateVueProject generates a Vue project structure
func (g *Generator) generateVueProject(files map[string]string, config metadata.ProjectConfig, isTypeScript bool) map[string]string {
	frontendConfig := config.FrontendConfig

	// package.json
	files["package.json"] = g.generateVuePackageJson(config, isTypeScript)

	// Vite config
	if isTypeScript {
		files["vite.config.ts"] = g.generateVueViteConfigTS(config)
	} else {
		files["vite.config.js"] = g.generateVueViteConfigJS(config)
	}

	// TypeScript config
	if isTypeScript {
		files["tsconfig.json"] = g.generateTSConfig(config)
		files["tsconfig.node.json"] = g.generateTSConfigNode(config)
		files["src/env.d.ts"] = g.generateVueEnvDTS(config)
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
		files["src/main.ts"] = g.generateVueMainTS(config)
	} else {
		files["src/main.js"] = g.generateVueMainJS(config)
	}

	// App component
	if isTypeScript {
		files["src/App.vue"] = g.generateVueAppTS(config)
	} else {
		files["src/App.vue"] = g.generateVueAppJS(config)
	}

	// Example component
	if isTypeScript {
		files["src/components/HelloWorld.vue"] = g.generateVueHelloWorldTS(config)
	} else {
		files["src/components/HelloWorld.vue"] = g.generateVueHelloWorldJS(config)
	}

	// Assets directory
	files["src/assets/.gitkeep"] = ""

	// Styles
	files["src/style.css"] = g.generateVueStyles(config, hasTailwind)

	// .gitignore
	files[".gitignore"] = g.generateFrontendGitignore(config)

	// README
	files["README.md"] = g.generateVueREADME(config)

	return files
}

// generateVuePackageJson generates package.json for Vue
func (g *Generator) generateVuePackageJson(config metadata.ProjectConfig, isTypeScript bool) string {
	frontendConfig := config.FrontendConfig

	dependencies := map[string]string{
		"vue":  "^3.4.0",
		"vite": "^5.0.0",
	}

	devDependencies := map[string]string{
		"@vitejs/plugin-vue": "^5.0.0",
	}

	if isTypeScript {
		devDependencies["@vue/tsconfig"] = "^0.5.0"
		devDependencies["typescript"] = "~5.2.2"
	}

	// Add feature dependencies
	for _, feature := range frontendConfig.Features {
		switch feature.ID {
		case "axios":
			dependencies["axios"] = "^1.6.0"
		}
	}

	// Add linter
	if frontendConfig.Linter.ID == "eslint" {
		devDependencies["eslint"] = "^8.54.0"
		devDependencies["eslint-plugin-vue"] = "^9.19.0"
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
    "preview": "vite preview"%s
  },
  "dependencies": {
%s  },
  "devDependencies": {
%s  }
}
`, config.ProjectName, func() string {
		if frontendConfig.Linter.ID == "eslint" {
			return `,
    "lint": "eslint . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts --fix --ignore-path .gitignore"`
		}
		return ""
	}(), depsStr, devDepsStr)
}

// generateVueViteConfigTS generates vite.config.ts for Vue
func (g *Generator) generateVueViteConfigTS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
})
`, config.ProjectName)
}

// generateVueViteConfigJS generates vite.config.js for Vue
func (g *Generator) generateVueViteConfigJS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
})
`, config.ProjectName)
}

// generateVueEnvDTS generates src/env.d.ts for Vue
func (g *Generator) generateVueEnvDTS(config metadata.ProjectConfig) string {
	return `/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
`
}

// generateVueMainTS generates src/main.ts for Vue
func (g *Generator) generateVueMainTS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

createApp(App).mount('#app')
`, config.ProjectName)
}

// generateVueMainJS generates src/main.js for Vue
func (g *Generator) generateVueMainJS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

createApp(App).mount('#app')
`, config.ProjectName)
}

// generateVueAppTS generates src/App.vue for Vue (TypeScript)
func (g *Generator) generateVueAppTS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`<script setup lang="ts">
import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <div class="app-container">
    <header>
      <h1>%s</h1>
      <p>Welcome to your Vue application!</p>
    </header>
    <main>
      <HelloWorld />
    </main>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  padding: 2rem;
}

header {
  text-align: center;
  margin-bottom: 2rem;
}

h1 {
  font-size: 2.5rem;
  font-weight: bold;
  margin-bottom: 1rem;
}
</style>
`, config.ProjectName)
}

// generateVueAppJS generates src/App.vue for Vue (JavaScript)
func (g *Generator) generateVueAppJS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`<script setup>
import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <div class="app-container">
    <header>
      <h1>%s</h1>
      <p>Welcome to your Vue application!</p>
    </header>
    <main>
      <HelloWorld />
    </main>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  padding: 2rem;
}

header {
  text-align: center;
  margin-bottom: 2rem;
}

h1 {
  font-size: 2.5rem;
  font-weight: bold;
  margin-bottom: 1rem;
}
</style>
`, config.ProjectName)
}

// generateVueHelloWorldTS generates src/components/HelloWorld.vue (TypeScript)
func (g *Generator) generateVueHelloWorldTS(config metadata.ProjectConfig) string {
	return `<script setup lang="ts">
const message = 'Hello, World!'
</script>

<template>
  <div class="hello-world">
    <p>{{ message }}</p>
  </div>
</template>

<style scoped>
.hello-world {
  padding: 1rem;
  text-align: center;
}

.hello-world p {
  font-size: 1.25rem;
  color: #333;
}
</style>
`
}

// generateVueHelloWorldJS generates src/components/HelloWorld.vue (JavaScript)
func (g *Generator) generateVueHelloWorldJS(config metadata.ProjectConfig) string {
	return `<script setup>
const message = 'Hello, World!'
</script>

<template>
  <div class="hello-world">
    <p>{{ message }}</p>
  </div>
</template>

<style scoped>
.hello-world {
  padding: 1rem;
  text-align: center;
}

.hello-world p {
  font-size: 1.25rem;
  color: #333;
}
</style>
`
}

// generateVueStyles generates src/style.css for Vue
func (g *Generator) generateVueStyles(config metadata.ProjectConfig, hasTailwind bool) string {
	if hasTailwind {
		return `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}
`
	}
	return `/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  line-height: 1.6;
  color: #333;
}
`
}

// generateVueREADME generates README.md for Vue
func (g *Generator) generateVueREADME(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`# %s

This is a Vue 3 project generated with [go-ctl](https://github.com/syst3mctl/go-ctl).

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Project Setup

`+"```sh"+`
npm install
`+"```"+`

### Compile and Hot-Reload for Development

`+"```sh"+`
npm run dev
`+"```"+`

### Compile and Minify for Production

`+"```sh"+`
npm run build
`+"```"+`

### Preview Production Build

`+"```sh"+`
npm run preview
`+"```"+`

## Generated with ❤️ by go-ctl
`, config.ProjectName)
}

// generateSvelteProject generates a Svelte project structure
func (g *Generator) generateSvelteProject(files map[string]string, config metadata.ProjectConfig, isTypeScript bool) map[string]string {
	frontendConfig := config.FrontendConfig

	// package.json
	files["package.json"] = g.generateSveltePackageJson(config, isTypeScript)

	// Vite config
	if isTypeScript {
		files["vite.config.ts"] = g.generateSvelteViteConfigTS(config)
	} else {
		files["vite.config.js"] = g.generateSvelteViteConfigJS(config)
	}

	// TypeScript config
	if isTypeScript {
		files["tsconfig.json"] = g.generateTSConfig(config)
		files["tsconfig.node.json"] = g.generateTSConfigNode(config)
		files["src/app.d.ts"] = g.generateSvelteAppDTS(config)
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
		files["src/main.ts"] = g.generateSvelteMainTS(config)
	} else {
		files["src/main.js"] = g.generateSvelteMainJS(config)
	}

	// App component
	if isTypeScript {
		files["src/App.svelte"] = g.generateSvelteAppTS(config)
	} else {
		files["src/App.svelte"] = g.generateSvelteAppJS(config)
	}

	// Example component
	if isTypeScript {
		files["src/components/HelloWorld.svelte"] = g.generateSvelteHelloWorldTS(config)
	} else {
		files["src/components/HelloWorld.svelte"] = g.generateSvelteHelloWorldJS(config)
	}

	// Assets directory
	files["src/assets/.gitkeep"] = ""

	// Styles
	files["src/app.css"] = g.generateSvelteStyles(config, hasTailwind)

	// .gitignore
	files[".gitignore"] = g.generateFrontendGitignore(config)

	// README
	files["README.md"] = g.generateSvelteREADME(config)

	return files
}

// generateSveltePackageJson generates package.json for Svelte
func (g *Generator) generateSveltePackageJson(config metadata.ProjectConfig, isTypeScript bool) string {
	frontendConfig := config.FrontendConfig

	dependencies := map[string]string{
		"svelte": "^4.2.0",
		"vite":   "^5.0.0",
	}

	devDependencies := map[string]string{
		"@sveltejs/vite-plugin-svelte": "^3.0.0",
	}

	if isTypeScript {
		devDependencies["typescript"] = "~5.2.2"
		devDependencies["@tsconfig/svelte"] = "^5.0.0"
		devDependencies["svelte-check"] = "^3.6.0"
	}

	// Add feature dependencies
	for _, feature := range frontendConfig.Features {
		switch feature.ID {
		case "axios":
			dependencies["axios"] = "^1.6.0"
		}
	}

	// Add linter
	if frontendConfig.Linter.ID == "eslint" {
		devDependencies["eslint"] = "^8.54.0"
		devDependencies["eslint-plugin-svelte"] = "^2.35.0"
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
		devDependencies["prettier-plugin-svelte"] = "^3.1.0"
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
    "dev": "vite dev",
    "build": "vite build",
    "preview": "vite preview"%s%s
  },
  "dependencies": {
%s  },
  "devDependencies": {
%s  }
}
`, config.ProjectName, func() string {
		if frontendConfig.Linter.ID == "eslint" {
			return `,
    "lint": "eslint . --fix"`
		}
		return ""
	}(), func() string {
		if isTypeScript {
			return `,
    "check": "svelte-check --tsconfig ./tsconfig.json"`
		}
		return ""
	}(), depsStr, devDepsStr)
}

// generateSvelteViteConfigTS generates vite.config.ts for Svelte
func (g *Generator) generateSvelteViteConfigTS(config metadata.ProjectConfig) string {
	return `import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
})
`
}

// generateSvelteViteConfigJS generates vite.config.js for Svelte
func (g *Generator) generateSvelteViteConfigJS(config metadata.ProjectConfig) string {
	return `import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
})
`
}

// generateSvelteAppDTS generates src/app.d.ts for Svelte
func (g *Generator) generateSvelteAppDTS(config metadata.ProjectConfig) string {
	return `// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface Platform {}
	}
}

export {};
`
}

// generateSvelteMainTS generates src/main.ts for Svelte
func (g *Generator) generateSvelteMainTS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`import './app.css'
import App from './App.svelte'

const app = new App({
  target: document.getElementById('app')!,
})

export default app
`, config.ProjectName)
}

// generateSvelteMainJS generates src/main.js for Svelte
func (g *Generator) generateSvelteMainJS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`import './app.css'
import App from './App.svelte'

const app = new App({
  target: document.getElementById('app'),
})

export default app
`, config.ProjectName)
}

// generateSvelteAppTS generates src/App.svelte for Svelte (TypeScript)
func (g *Generator) generateSvelteAppTS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`<script lang="ts">
  import HelloWorld from './components/HelloWorld.svelte'
</script>

<main>
  <div class="app-container">
    <header>
      <h1>%s</h1>
      <p>Welcome to your Svelte application!</p>
    </header>
    <HelloWorld />
  </div>
</main>

<style>
  .app-container {
    min-height: 100vh;
    padding: 2rem;
  }

  header {
    text-align: center;
    margin-bottom: 2rem;
  }

  h1 {
    font-size: 2.5rem;
    font-weight: bold;
    margin-bottom: 1rem;
  }
</style>
`, config.ProjectName)
}

// generateSvelteAppJS generates src/App.svelte for Svelte (JavaScript)
func (g *Generator) generateSvelteAppJS(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`<script>
  import HelloWorld from './components/HelloWorld.svelte'
</script>

<main>
  <div class="app-container">
    <header>
      <h1>%s</h1>
      <p>Welcome to your Svelte application!</p>
    </header>
    <HelloWorld />
  </div>
</main>

<style>
  .app-container {
    min-height: 100vh;
    padding: 2rem;
  }

  header {
    text-align: center;
    margin-bottom: 2rem;
  }

  h1 {
    font-size: 2.5rem;
    font-weight: bold;
    margin-bottom: 1rem;
  }
</style>
`, config.ProjectName)
}

// generateSvelteHelloWorldTS generates src/components/HelloWorld.svelte (TypeScript)
func (g *Generator) generateSvelteHelloWorldTS(config metadata.ProjectConfig) string {
	return `<script lang="ts">
  let message = 'Hello, World!'
</script>

<div class="hello-world">
  <p>{message}</p>
</div>

<style>
  .hello-world {
    padding: 1rem;
    text-align: center;
  }

  .hello-world p {
    font-size: 1.25rem;
    color: #333;
  }
</style>
`
}

// generateSvelteHelloWorldJS generates src/components/HelloWorld.svelte (JavaScript)
func (g *Generator) generateSvelteHelloWorldJS(config metadata.ProjectConfig) string {
	return `<script>
  let message = 'Hello, World!'
</script>

<div class="hello-world">
  <p>{message}</p>
</div>

<style>
  .hello-world {
    padding: 1rem;
    text-align: center;
  }

  .hello-world p {
    font-size: 1.25rem;
    color: #333;
  }
</style>
`
}

// generateSvelteStyles generates src/app.css for Svelte
func (g *Generator) generateSvelteStyles(config metadata.ProjectConfig, hasTailwind bool) string {
	if hasTailwind {
		return `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}
`
	}
	return `/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  line-height: 1.6;
  color: #333;
}
`
}

// generateSvelteREADME generates README.md for Svelte
func (g *Generator) generateSvelteREADME(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`# %s

This is a Svelte project generated with [go-ctl](https://github.com/syst3mctl/go-ctl).

## Get started

Install the dependencies...

`+"```bash"+`
npm install
`+"```"+`

...then start [Vite](https://vitejs.dev):

`+"```bash"+`
npm run dev
`+"```"+`

Navigate to [localhost:5173](http://localhost:5173). You should see your app running. Edit a component file in `+"`src`"+`, save it, and the page will reload.

## Building and previewing the production build

If you've already installed the dependencies with `+"`npm install`"+`, skip ahead. If not, run:

`+"```bash"+`
npm install
`+"```"+`

To create a production version of your app:

`+"```bash"+`
npm run build
`+"```"+`

You can preview the production build with:

`+"```bash"+`
npm run preview
`+"```"+`

## Generated with ❤️ by go-ctl
`, config.ProjectName)
}
