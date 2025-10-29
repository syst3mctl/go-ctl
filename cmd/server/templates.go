package main

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>go-ctl - Go Project Initializr</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-core.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js"></script>
    <style>
        .gradient-bg {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .card-hover {
            transition: transform 0.2s, box-shadow 0.2s;
        }
        .card-hover:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
        }
        .loading {
            opacity: 0.6;
            pointer-events: none;
        }
        .file-item {
            transition: all 0.2s ease;
        }
        .file-item:hover {
            background-color: #f3f4f6;
            border-color: #d1d5db;
        }
        .file-item.active {
            background-color: #dbeafe !important;
            border-color: #93c5fd !important;
        }
        .modal-backdrop {
            backdrop-filter: blur(4px);
        }
        pre[class*="language-"] {
            background: #2d3748 !important;
            margin: 0 !important;
        }
        code[class*="language-"] {
            color: #e2e8f0 !important;
            font-size: 0.875rem !important;
            line-height: 1.5 !important;
        }
        .file-tree-container {
            max-height: calc(80vh - 200px);
            overflow-y: auto;
        }
        .file-content-container {
            max-height: calc(80vh - 200px);
            overflow-y: auto;
        }
    </style>
</head>
<body class="bg-gray-50 min-h-screen">
    <!-- Header -->
    <header class="gradient-bg text-white py-6">
        <div class="container mx-auto px-6">
            <h1 class="text-4xl font-bold flex items-center">
                <i class="fas fa-cube mr-3"></i>
                go-ctl
            </h1>
            <p class="text-blue-100 mt-2">Go Project Initializr - Generate production-ready Go projects in seconds</p>
        </div>
    </header>

    <div class="container mx-auto px-6 py-8">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">

            <!-- Left Side: Configuration Form -->
            <div class="bg-white rounded-lg shadow-lg p-6 card-hover">
                <h2 class="text-2xl font-bold mb-6 text-gray-800 flex items-center">
                    <i class="fas fa-cog mr-2 text-blue-600"></i>
                    Project Configuration
                </h2>

                <form id="project-form" action="/generate" method="POST" class="space-y-6">

                    <!-- Project Metadata -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                <i class="fas fa-tag mr-1"></i>Project Name
                            </label>
                            <input type="text"
                                   name="projectName"
                                   value="my-go-app"
                                   class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                   placeholder="my-awesome-app"
                                   required>
                        </div>

                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                <i class="fab fa-golang mr-1"></i>Go Version
                            </label>
                            <select name="goVersion"
                                    class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                                {{range .Options.GoVersions}}
                                <option value="{{.}}"{{if eq . "1.23"}} selected{{end}}>{{.}}</option>
                                {{end}}
                            </select>
                        </div>
                    </div>

                    <!-- HTTP Framework -->
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-3">
                            <i class="fas fa-server mr-1"></i>HTTP Framework
                        </label>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {{range .Options.Http}}
                            <div class="relative">
                                <input type="radio"
                                       name="httpPackage"
                                       value="{{.ID}}"
                                       id="http-{{.ID}}"
                                       {{if eq .ID "gin"}}checked{{end}}
                                       class="sr-only peer">
                                <label for="http-{{.ID}}"
                                       class="flex p-3 bg-gray-50 border border-gray-300 rounded-lg cursor-pointer peer-checked:bg-blue-50 peer-checked:border-blue-500 hover:bg-gray-100">
                                    <div class="flex-1">
                                        <div class="font-semibold text-gray-800">{{.Name}}</div>
                                        <div class="text-sm text-gray-600">{{.Description}}</div>
                                    </div>
                                    <div class="flex-shrink-0 ml-2">
                                        <i class="fas fa-check-circle text-blue-500 opacity-0 peer-checked:opacity-100"></i>
                                    </div>
                                </label>
                            </div>
                            {{end}}
                        </div>
                    </div>

                    <!-- Database -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-3">
                                <i class="fas fa-database mr-1"></i>Database
                            </label>
                            <div class="space-y-2">
                                <div>
                                    <input type="radio" name="database" value="" id="db-none" checked class="mr-2">
                                    <label for="db-none" class="text-sm">None</label>
                                </div>
                                {{range .Options.Databases}}
                                <div>
                                    <input type="radio" name="database" value="{{.ID}}" id="db-{{.ID}}" class="mr-2">
                                    <label for="db-{{.ID}}" class="text-sm font-medium">{{.Name}}</label>
                                    <p class="text-xs text-gray-600 ml-6">{{.Description}}</p>
                                </div>
                                {{end}}
                            </div>
                        </div>

                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-3">
                                <i class="fas fa-plug mr-1"></i>Database Driver
                            </label>
                            <div class="space-y-2">
                                <div>
                                    <input type="radio" name="dbDriver" value="" id="driver-none" checked class="mr-2">
                                    <label for="driver-none" class="text-sm">None</label>
                                </div>
                                {{range .Options.DbDrivers}}
                                <div>
                                    <input type="radio" name="dbDriver" value="{{.ID}}" id="driver-{{.ID}}" class="mr-2">
                                    <label for="driver-{{.ID}}" class="text-sm font-medium">{{.Name}}</label>
                                    <p class="text-xs text-gray-600 ml-6">{{.Description}}</p>
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>

                    <!-- Features -->
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-3">
                            <i class="fas fa-puzzle-piece mr-1"></i>Additional Features
                        </label>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {{range .Options.Features}}
                            <div class="flex items-start">
                                <input type="checkbox"
                                       name="features"
                                       value="{{.ID}}"
                                       id="feature-{{.ID}}"
                                       {{if or (eq .ID "gitignore") (eq .ID "makefile")}}checked{{end}}
                                       class="mt-1 mr-3 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500">
                                <div class="flex-1">
                                    <label for="feature-{{.ID}}" class="text-sm font-medium text-gray-700 cursor-pointer">
                                        {{.Name}}
                                    </label>
                                    <p class="text-xs text-gray-600">{{.Description}}</p>
                                </div>
                            </div>
                            {{end}}
                        </div>
                    </div>

                    <!-- Dynamic Package Search -->
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">
                            <i class="fas fa-search mr-1"></i>Add Dependencies
                        </label>
                        <input type="search"
                               name="q"
                               class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent mb-3"
                               placeholder="Search pkg.go.dev for packages..."
                               hx-get="/search-packages"
                               hx-trigger="keyup changed delay:500ms"
                               hx-target="#search-results"
                               hx-swap="innerHTML"
                               hx-indicator="#search-loading">

                        <!-- Loading indicator -->
                        <div id="search-loading" class="htmx-indicator text-center py-2">
                            <i class="fas fa-spinner fa-spin text-blue-500"></i>
                            <span class="ml-2 text-sm text-gray-600">Searching...</span>
                        </div>

                        <!-- Search results container -->
                        <div id="search-results" class="max-h-48 overflow-y-auto mb-4 space-y-2"></div>

                        <!-- Selected packages -->
                        <div>
                            <h4 class="text-sm font-medium text-gray-700 mb-2">Selected Packages:</h4>
                            <div id="selected-packages" class="space-y-2 min-h-[2rem] p-2 border border-gray-200 rounded-lg bg-gray-50">
                                <p class="text-sm text-gray-500 italic">No packages selected</p>
                            </div>
                        </div>
                    </div>

                    <!-- Action Buttons -->
                    <div class="flex space-x-4 pt-6">
                        <button type="button"
                                hx-post="/explore"
                                hx-include="#project-form"
                                hx-target="#file-tree-content"
                                hx-swap="innerHTML"
                                hx-indicator="#explore-loading"
                                onclick="openExploreModal()"
                                class="flex-1 bg-gray-600 hover:bg-gray-700 text-white font-semibold py-3 px-6 rounded-lg transition duration-200 flex items-center justify-center">
                            <i class="fas fa-eye mr-2"></i>
                            Preview Structure
                        </button>

                        <button type="submit"
                                class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-6 rounded-lg transition duration-200 flex items-center justify-center">
                            <i class="fas fa-download mr-2"></i>
                            Generate Project
                        </button>
                    </div>
                </form>
            </div>

            <!-- Right Side: Preview -->
            <div class="bg-white rounded-lg shadow-lg p-6 card-hover">
                <h2 class="text-2xl font-bold mb-6 text-gray-800 flex items-center">
                    <i class="fas fa-folder-tree mr-2 text-green-600"></i>
                    Project Preview
                </h2>

                <div class="text-center py-12 text-gray-500">
                    <i class="fas fa-folder-open text-4xl mb-4"></i>
                    <p class="text-lg">Click "Preview Structure" to explore your project</p>
                    <p class="text-sm mt-2">Browse through the generated files and see their contents in an interactive modal.</p>
                </div>
            </div>
        </div>

        <!-- Explore Modal -->
        <div id="explore-modal" class="fixed inset-0 bg-black bg-opacity-50 modal-backdrop hidden z-50" onclick="closeExploreModal(event)">
            <div class="flex items-center justify-center min-h-screen p-4">
                <div class="bg-white rounded-lg shadow-2xl w-full max-w-7xl h-[80vh] flex flex-col" onclick="event.stopPropagation()">
                    <!-- Modal Header -->
                    <div class="flex items-center justify-between p-6 border-b border-gray-200">
                        <div class="flex items-center">
                            <i class="fas fa-folder-open text-2xl text-blue-600 mr-3"></i>
                            <h2 class="text-2xl font-bold text-gray-800">Project Explorer</h2>
                        </div>
                        <button onclick="closeExploreModal()" class="text-gray-400 hover:text-gray-600 transition duration-200">
                            <i class="fas fa-times text-2xl"></i>
                        </button>
                    </div>

                    <!-- Modal Body -->
                    <div class="flex flex-1 overflow-hidden">
                        <!-- File Tree Sidebar -->
                        <div class="w-1/3 border-r border-gray-200 bg-gray-50">
                            <div class="p-4 h-full">
                                <h3 class="text-lg font-semibold text-gray-700 mb-3 flex items-center">
                                    <i class="fas fa-sitemap mr-2 text-gray-500"></i>
                                    File Structure
                                </h3>
                                <!-- Loading indicator -->
                                <div id="explore-loading" class="htmx-indicator text-center py-8">
                                    <i class="fas fa-spinner fa-spin text-2xl text-blue-500"></i>
                                    <p class="mt-2 text-gray-600">Generating structure...</p>
                                </div>
                                <!-- File tree content -->
                                <div id="file-tree-content" class="space-y-1 file-tree-container">
                                    <!-- Files will be loaded here -->
                                </div>
                            </div>
                        </div>

                        <!-- File Content Area -->
                        <div class="flex-1 flex flex-col">
                            <!-- File Header -->
                            <div class="p-4 border-b border-gray-200 bg-gray-50">
                                <div class="flex items-center justify-between">
                                    <div id="current-file-header" class="flex items-center text-gray-600">
                                        <i class="fas fa-file-code mr-2"></i>
                                        <span>Select a file to preview</span>
                                    </div>
                                    <div class="flex items-center space-x-2">
                                        <button onclick="copyFileContent()" id="copy-btn" class="hidden bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-sm transition duration-200">
                                            <i class="fas fa-copy mr-1"></i>Copy
                                        </button>
                                    </div>
                                </div>
                            </div>

                            <!-- File Content -->
                            <div class="flex-1 overflow-hidden">
                                <div id="file-content" class="h-full file-content-container bg-gray-800">
                                    <div class="flex items-center justify-center h-full text-gray-400">
                                        <div class="text-center p-8">
                                            <i class="fas fa-mouse-pointer text-4xl mb-4 text-gray-500"></i>
                                            <p class="text-lg text-gray-300">Click on a file to view its content</p>
                                            <p class="text-sm mt-2 text-gray-500">Browse the file tree on the left to explore your generated project</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Modal Footer -->
                    <div class="p-4 border-t border-gray-200 bg-gray-50 flex justify-between items-center">
                        <div class="text-sm text-gray-600">
                            <i class="fas fa-info-circle mr-1"></i>
                            Click files to preview content • All files will be included in the download
                        </div>
                        <div class="flex space-x-3">
                            <button onclick="closeExploreModal()" class="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded transition duration-200">
                                Close
                            </button>
                            <button onclick="downloadProject()" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded transition duration-200">
                                <i class="fas fa-download mr-2"></i>Download Project
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Footer -->
        <footer class="mt-12 text-center text-gray-600 border-t pt-8">
            <div class="flex justify-center items-center space-x-6">
                <a href="https://github.com/syst3mctl/go-ctl" class="hover:text-blue-600 transition duration-200">
                    <i class="fab fa-github mr-2"></i>GitHub
                </a>
                <span class="text-gray-400">•</span>
                <span>Built with ❤️ by systemctl</span>
                <span class="text-gray-400">•</span>
                <span>Powered by Go + HTMX</span>
            </div>
            <p class="mt-2 text-sm">Generate production-ready Go projects with clean architecture</p>
        </footer>
    </div>

    <script>
        // Modal functions
        function openExploreModal() {
            document.getElementById('explore-modal').classList.remove('hidden');
            document.body.style.overflow = 'hidden';
        }

        function closeExploreModal(event) {
            if (!event || event.target.id === 'explore-modal') {
                document.getElementById('explore-modal').classList.add('hidden');
                document.body.style.overflow = 'auto';
            }
        }

        // File selection and content loading
        function selectFile(filepath, filename, element) {
            // Remove active state from all file items
            document.querySelectorAll('.file-item').forEach(item => {
                item.classList.remove('active');
            });

            // Add active state to selected file
            element.classList.add('active');

            // Update header
            const icon = getFileIcon(filename);
            document.getElementById('current-file-header').innerHTML =
                '<i class="' + icon + ' mr-2"></i>' + filepath;

            // Show copy button
            document.getElementById('copy-btn').classList.remove('hidden');

            // Load file content via HTMX
            htmx.ajax('GET', '/file-content?path=' + encodeURIComponent(filepath), {
                target: '#file-content',
                swap: 'innerHTML'
            });

            // Trigger syntax highlighting after content loads
            setTimeout(() => {
                if (window.Prism) {
                    Prism.highlightAll();
                }
            }, 100);
        }

        function getFileIcon(filename) {
            const ext = filename.split('.').pop().toLowerCase();
            switch(ext) {
                case 'go': return 'fab fa-golang text-blue-500';
                case 'mod': return 'fas fa-cube text-green-500';
                case 'json': return 'fas fa-brackets-curly text-yellow-500';
                case 'yml':
                case 'yaml': return 'fas fa-file-code text-orange-500';
                case 'md': return 'fab fa-markdown text-blue-600';
                case 'toml': return 'fas fa-cog text-gray-500';
                case 'env': return 'fas fa-key text-green-600';
                default: return 'fas fa-file-code text-gray-500';
            }
        }

        // Copy file content
        function copyFileContent() {
            const content = document.getElementById('file-content').textContent;
            navigator.clipboard.writeText(content).then(() => {
                const btn = document.getElementById('copy-btn');
                const originalText = btn.innerHTML;
                btn.innerHTML = '<i class="fas fa-check mr-1"></i>Copied!';
                btn.classList.add('bg-green-500', 'hover:bg-green-600');
                btn.classList.remove('bg-blue-500', 'hover:bg-blue-600');

                setTimeout(() => {
                    btn.innerHTML = originalText;
                    btn.classList.remove('bg-green-500', 'hover:bg-green-600');
                    btn.classList.add('bg-blue-500', 'hover:bg-blue-600');
                }, 2000);
            });
        }

        // Download project from modal
        function downloadProject() {
            const form = document.getElementById('project-form');
            form.submit();
            closeExploreModal();
        }

        // Remove placeholder text when packages are selected
        document.addEventListener('htmx:afterRequest', function(event) {
            const selectedPackages = document.getElementById('selected-packages');
            const placeholder = selectedPackages.querySelector('p.italic');
            if (placeholder && selectedPackages.children.length > 1) {
                placeholder.remove();
            }
        });

        // Form validation
        document.getElementById('project-form').addEventListener('submit', function(e) {
            const projectName = document.querySelector('input[name="projectName"]').value;
            if (!projectName.trim()) {
                e.preventDefault();
                alert('Please enter a project name');
                return;
            }

            // Show loading state
            const submitBtn = e.target.querySelector('button[type="submit"]');
            submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin mr-2"></i>Generating...';
            submitBtn.disabled = true;
        });

        // ESC key to close modal
        document.addEventListener('keydown', function(event) {
            if (event.key === 'Escape') {
                closeExploreModal();
            }
        });
    </script>
</body>
</html>`

const exploreTemplate = `
{{range .Files}}
<div class="file-item border border-gray-200 rounded p-2 cursor-pointer hover:bg-gray-100 transition duration-150 mb-1"
     onclick="selectFile('{{.Path}}', '{{.Name}}', this)">
    <div class="flex items-center">
        <i class="{{.Icon}} mr-2"></i>
        <span class="text-sm font-medium text-gray-700">{{.Name}}</span>
    </div>
    <div class="text-xs text-gray-500 ml-6 mt-1">{{.Path}}</div>
</div>
{{end}}

{{if not .Files}}
<div class="text-center py-8 text-gray-500">
    <i class="fas fa-exclamation-triangle text-2xl mb-2"></i>
    <p>No files generated</p>
    <p class="text-xs mt-1">Please check your configuration</p>
</div>
{{end}}
`

const fileContentTemplate = `
<div class="h-full">
    <pre class="language-{{.Language}} h-full m-0"><code class="language-{{.Language}}">{{.Content}}</code></pre>
</div>
`

const searchResultsTemplate = `
{{if .Results}}
    {{range .Results}}
    <div class="flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 transition duration-150">
        <div class="flex-1 min-w-0">
            <div class="font-semibold text-gray-900 text-sm">{{.Path}}</div>
            <div class="text-xs text-gray-600 mt-1 truncate">{{.Synopsis}}</div>
        </div>
        <button type="button"
                hx-post="/add-package"
                hx-vals='{"pkgPath": "{{.Path}}"}'
                hx-target="#selected-packages"
                hx-swap="beforeend"
                class="ml-3 bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-sm font-medium transition duration-150 flex items-center">
            <i class="fas fa-plus mr-1"></i>Add
        </button>
    </div>
    {{end}}
{{else}}
    <div class="text-center py-4 text-gray-500">
        <i class="fas fa-search text-2xl mb-2"></i>
        <p class="text-sm">No packages found for "{{.Query}}"</p>
        <p class="text-xs text-gray-400 mt-1">Try a different search term</p>
    </div>
{{end}}
`

const selectedPackageTemplate = `
<div id="pkg-{{.ID}}" class="flex items-center justify-between bg-blue-100 border border-blue-200 rounded-lg p-2">
    <div class="flex items-center">
        <i class="fas fa-cube text-blue-600 mr-2"></i>
        <span class="text-sm font-medium text-blue-800">{{.PkgPath}}</span>
    </div>

    <!-- Hidden input for form submission -->
    <input type="hidden" name="customPackages" value="{{.PkgPath}}">

    <button type="button"
            hx-target="#pkg-{{.ID}}"
            hx-swap="delete"
            class="text-red-500 hover:text-red-700 font-bold text-sm ml-2 transition duration-150">
        <i class="fas fa-times"></i>
    </button>
</div>
`
