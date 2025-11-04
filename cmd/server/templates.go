package main

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>go-ctl - Go Project Initializr</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
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

        /* File Explorer Styles */
        .tree-item {
            user-select: none;
        }

        .tree-item {
            display: none;
        }

        .tree-item.show {
            display: block;
        }

        .tree-item[data-level="0"] {
            display: block;
        }

        .folder-item:hover {
            background-color: rgba(59, 130, 246, 0.1);
        }

        .file-item:hover {
            background-color: rgba(0, 0, 0, 0.05);
        }

        .file-item.selected {
            background-color: rgba(59, 130, 246, 0.2);
            border-left: 3px solid #3b82f6;
        }

        .folder-chevron.expanded {
            transform: rotate(90deg);
        }

        .file-tree-container {
            font-size: 14px;
            line-height: 1.4;
        }

        .file-content-container {
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        }

        /* Scrollbar styles */
        .file-tree-container::-webkit-scrollbar,
        .file-content-container::-webkit-scrollbar {
            width: 8px;
        }

        .file-tree-container::-webkit-scrollbar-track,
        .file-content-container::-webkit-scrollbar-track {
            background: #f1f1f1;
            border-radius: 4px;
        }

        .file-tree-container::-webkit-scrollbar-thumb,
        .file-content-container::-webkit-scrollbar-thumb {
            background: #c1c1c1;
            border-radius: 4px;
        }

        .file-tree-container::-webkit-scrollbar-thumb:hover,
        .file-content-container::-webkit-scrollbar-thumb:hover {
            background: #a1a1a1;
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
    <header class="bg-[#11A32B] text-white py-6">
        <div class="container mx-auto px-6">
            <h1 class="text-4xl font-bold flex items-center">
                <i class="fas fa-cube mr-3"></i>
                go-ctl
            </h1>
            <p class="text-blue-100 mt-2">Go Project Initializer - Generate production-ready Go projects in seconds</p>
        </div>
    </header>

    <div class="container mx-auto px-6 py-8">
        <div class="grid grid-cols-1 lg:grid-cols-1 gap-8">

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
                                {{range .Options.Databases}}
                                <div>
                                    <input type="checkbox" name="databases" value="{{.ID}}" id="db-{{.ID}}" class="mr-2">
                                    <label for="db-{{.ID}}" class="text-sm font-medium">{{.Name}}</label>
                                    <p class="text-xs text-gray-600 ml-6">{{.Description}}</p>
                                </div>
                                {{end}}
                            </div>
                        </div>

                        <div id="driver-section" style="display: none;">
                            <label class="block text-sm font-medium text-gray-700 mb-3">
                                <i class="fas fa-plug mr-1"></i>Database Drivers
                            </label>
                            <div id="driver-container" class="space-y-4">
                                <!-- Dynamic driver options will be inserted here -->
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
                               hx-get="/fetch-packages"
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
        </div>

        <!-- Explore Modal -->
        <div id="explore-modal" class="fixed inset-0 bg-black bg-opacity-50 modal-backdrop hidden z-50" onclick="closeExploreModal(event)">
            <div class="flex items-center justify-center min-h-screen p-4">
                <div class="bg-white rounded-lg shadow-2xl w-full max-w-7xl h-[85vh] flex flex-col" onclick="event.stopPropagation()">
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
                            <div class="p-4 h-full flex flex-col">
                                <h3 class="text-lg font-semibold text-gray-700 mb-3 flex items-center">
                                    <i class="fas fa-folder-tree mr-2 text-blue-600"></i>
                                    Explorer
                                </h3>
                                <!-- Loading indicator -->

                                <!-- File tree content -->
                                <div id="file-tree-content" class="file-tree-container overflow-auto flex-1">
                                    <!-- Files will be loaded here -->
                                </div>
                            </div>
                        </div>

                        <!-- File Content Area -->
                        <div class="flex-1 flex flex-col bg-gray-900">
                            <!-- File Header -->
                            <div class="p-3 border-b border-gray-700 bg-gray-800">
                                <div class="flex items-center justify-between">
                                    <div id="current-file-header" class="flex items-center text-gray-300">
                                        <i class="fas fa-file-code mr-2"></i>
                                        <span class="text-sm">Select a file to preview</span>
                                    </div>
                                    <div class="flex items-center space-x-2">
                                        <button onclick="copyFileContent()" id="copy-btn" class="hidden bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-xs transition duration-200">
                                            <i class="fas fa-copy mr-1"></i>Copy
                                        </button>
                                    </div>
                                </div>
                            </div>

                            <!-- File Content -->
                            <div class="flex-1 overflow-hidden">
                                <div id="file-content" class="h-full file-content-container bg-gray-900 overflow-auto">
                                    <div class="flex items-center justify-center h-full text-gray-400">
                                        <div class="text-center p-8">
                                            <i class="fas fa-code text-4xl mb-4 text-gray-600"></i>
                                            <p class="text-lg text-gray-300 mb-2">Select a file to view its content</p>
                                            <p class="text-sm text-gray-500">Click on folders to expand and files to preview</p>
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

        // Database driver selection functions
        function updateDriverOptions() {
            const databaseCheckboxes = document.querySelectorAll('input[name="databases"]:checked');
            const driverSection = document.getElementById('driver-section');
            const driverContainer = document.getElementById('driver-container');

            // Save current selected driver values before clearing
            const savedSelections = {};
            databaseCheckboxes.forEach(checkbox => {
                const database = checkbox.value;
                const driverRadio = document.querySelector('input[name="driver_' + database + '"]:checked');
                if (driverRadio) {
                    savedSelections[database] = driverRadio.value;
                }
            });

            // Clear existing driver options
            driverContainer.innerHTML = '';

            if (databaseCheckboxes.length === 0) {
                driverSection.style.display = 'none';
                return;
            }

            driverSection.style.display = 'block';

            // Define driver options for each database
            const driverOptions = {
                'postgres': [
                    { id: 'gorm', name: 'GORM', description: 'The fantastic ORM library for Golang' },
                    { id: 'sqlx', name: 'sqlx', description: 'Extensions to database/sql with easier scanning' },
                    { id: 'ent', name: 'Ent', description: 'An entity framework for Go with schema-first code generation' },
                    { id: 'database-sql', name: 'database/sql', description: 'Standard library SQL interface' }
                ],
                'mysql': [
                    { id: 'gorm', name: 'GORM', description: 'The fantastic ORM library for Golang' },
                    { id: 'sqlx', name: 'sqlx', description: 'Extensions to database/sql with easier scanning' },
                    { id: 'ent', name: 'Ent', description: 'An entity framework for Go with schema-first code generation' },
                    { id: 'database-sql', name: 'database/sql', description: 'Standard library SQL interface' }
                ],
                'sqlite': [
                    { id: 'gorm', name: 'GORM', description: 'The fantastic ORM library for Golang' },
                    { id: 'sqlx', name: 'sqlx', description: 'Extensions to database/sql with easier scanning' },
                    { id: 'ent', name: 'Ent', description: 'An entity framework for Go with schema-first code generation' },
                    { id: 'database-sql', name: 'database/sql', description: 'Standard library SQL interface' }
                ],
                'mongodb': [
                    { id: 'mongo-driver', name: 'MongoDB Driver', description: 'Official MongoDB Go driver' }
                ],
                'redis': [
                    { id: 'redis-client', name: 'Redis Client', description: 'Redis client with Cluster, Sentinel support' }
                ],
                'bigquery': [
                    { id: 'database-sql', name: 'database/sql', description: 'Standard library SQL interface' }
                ]
            };

            // Create driver sections for selected databases
            databaseCheckboxes.forEach(checkbox => {
                const database = checkbox.value;
                if (driverOptions[database]) {
                    const dbSection = document.createElement('div');
                    dbSection.className = 'border border-gray-200 rounded-lg p-4';
                    const savedSelection = savedSelections[database];
                    const defaultDriver = driverOptions[database][0].id;
                    const selectedDriver = savedSelection || defaultDriver;
                    
                    dbSection.innerHTML =
                        '<h4 class="font-medium text-gray-900 mb-3 capitalize">' + database + ' Driver</h4>' +
                        '<div class="space-y-2">' +
                            driverOptions[database].map(function(driver) {
                                const isChecked = driver.id === selectedDriver ? 'checked' : '';
                                return '<label class="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 border border-gray-200 cursor-pointer transition-colors duration-150">' +
                                    '<input type="radio" name="driver_' + database + '" value="' + driver.id + '" ' + isChecked + ' class="mt-1 text-blue-600 border-gray-300 focus:ring-blue-500" required>' +
                                    '<div class="flex-1 min-w-0">' +
                                        '<div class="text-sm font-medium text-gray-900">' + driver.name + '</div>' +
                                        '<div class="text-sm text-gray-500">' + driver.description + '</div>' +
                                    '</div>' +
                                '</label>';
                            }).join('') +
                        '</div>';
                    driverContainer.appendChild(dbSection);
                }
            });
        }

        // Initialize driver options on page load
        document.addEventListener('DOMContentLoaded', function() {
            // Add event listeners to database checkboxes
            document.querySelectorAll('input[name="databases"]').forEach(checkbox => {
                checkbox.addEventListener('change', updateDriverOptions);
            });

            // Initial call to set up driver options
            updateDriverOptions();
        });
        // File tree and explorer functions
        // Initialize file tree - expand root folders automatically
        function initializeFileTree() {
            // Show all root level items (level 0)
            const rootItems = document.querySelectorAll('.tree-item[data-level="0"]');
            rootItems.forEach(item => {
                if (item.dataset.isFolder === 'true') {
                    const chevron = item.querySelector('.folder-chevron');
                    if (chevron) {
                        chevron.classList.add('expanded');
                        showDirectChildren(item.dataset.path);
                    }
                }
                item.classList.add('show');
            });
        }

        // Tree functionality
        function toggleFolder(folderElement) {
            console.log('toggleFolder called with:', folderElement);
            const chevron = folderElement.querySelector('.folder-chevron');
            const isExpanded = chevron.classList.contains('expanded');
            const folderPath = folderElement.closest('.tree-item').dataset.path;

            chevron.classList.toggle('expanded');

            if (isExpanded) {
                // Collapse: hide all descendants
                hideDescendants(folderPath);
            } else {
                // Expand: show direct children only
                showDirectChildren(folderPath);
            }
        }

        function showDirectChildren(parentPath) {
            const allItems = document.querySelectorAll('.tree-item');
            const parentElement = document.querySelector('[data-path="' + parentPath + '"]');
            if (!parentElement) return;
            const parentLevel = parseInt(parentElement.dataset.level);

            allItems.forEach(item => {
                const itemPath = item.dataset.path;
                const itemLevel = parseInt(item.dataset.level);
                const isFolder = item.dataset.isFolder === 'true';

                // Show direct children (one level deeper)
                if (itemLevel === parentLevel + 1 && itemPath.startsWith(parentPath + '/')) {
                    if (!isFolder || isDirectChild(parentPath, itemPath)) {
                        item.classList.add('show');
                    }
                }
            });
        }

        function hideDescendants(parentPath) {
            const allItems = document.querySelectorAll('.tree-item');
            const parentElement = document.querySelector('[data-path="' + parentPath + '"]');
            if (!parentElement) return;
            const parentLevel = parseInt(parentElement.dataset.level);

            allItems.forEach(item => {
                const itemPath = item.dataset.path;
                const itemLevel = parseInt(item.dataset.level);

                // Hide all descendants (any level deeper)
                if (itemLevel > parentLevel && itemPath.startsWith(parentPath + '/')) {
                    item.classList.remove('show');
                    // Also collapse any expanded folders in descendants
                    const chevron = item.querySelector('.folder-chevron');
                    if (chevron) {
                        chevron.classList.remove('expanded');
                    }
                }
            });
        }

        function isDirectChild(parentPath, childPath) {
            if (!childPath.startsWith(parentPath + '/')) return false;
            const remainingPath = childPath.substring(parentPath.length + 1);
            return !remainingPath.includes('/');
        }

        function selectFile(path, name, element) {
            // Remove previous selection
            document.querySelectorAll('.file-item').forEach(item => {
                item.classList.remove('selected');
            });

            // Add selection to clicked file
            element.classList.add('selected');

            // Update header
            const header = document.getElementById('current-file-header');
            const icon = getFileIcon(name);
            header.innerHTML = '<i class="' + icon + ' mr-2"></i><span class="font-medium">' + name + '</span>';

            // Show copy button
            document.getElementById('copy-btn').classList.remove('hidden');

            // Load file content
            loadFileContent(path, name);
        }

        function getFileIcon(filename) {
            const ext = filename.split('.').pop().toLowerCase();
            switch(ext) {
                case 'go': return 'fab fa-golang text-blue-500';
                case 'md': return 'fab fa-markdown text-blue-600';
                case 'json': return 'fas fa-code text-yellow-600';
                case 'yml':
                case 'yaml': return 'fas fa-code text-red-600';
                case 'toml': return 'fas fa-code text-purple-600';
                case 'mod': return 'fas fa-cube text-green-500';
                default: return 'fas fa-file-code text-gray-500';
            }
        }

        // Load file content function
        function loadFileContent(path, name) {
            const contentDiv = document.getElementById('file-content');

            // Show loading state
            contentDiv.innerHTML = '<div class="flex items-center justify-center h-full text-gray-400">' +
                '<div class="text-center">' +
                '<i class="fas fa-spinner fa-spin text-2xl mb-4 text-blue-400"></i>' +
                '<p class="text-gray-300">Loading ' + name + '...</p>' +
                '</div>' +
                '</div>';

            // Get project configuration from form
            const form = document.getElementById('project-form');
            const formData = new FormData(form);

            // Build URL with project configuration parameters
            const params = new URLSearchParams();
            params.append('path', path);
            params.append('projectName', formData.get('projectName') || 'my-go-app');
            params.append('goVersion', formData.get('goVersion') || '1.23');
            params.append('httpPackage', formData.get('httpPackage') || 'gin');

            // Handle multiple database selections
            const selectedDatabases = formData.getAll('databases');
            if (selectedDatabases.length > 0) {
                params.append('databases', selectedDatabases.join(','));
                // Add driver selections for each database
                selectedDatabases.forEach(dbId => {
                    const driverValue = formData.get('driver_' + dbId);
                    if (driverValue) {
                        params.append('driver_' + dbId, driverValue);
                    }
                });
            }

            // Fetch file content with configuration
            fetch('/file-content?' + params.toString())
                .then(response => {
                    if (!response.ok) {
                        throw new Error('HTTP ' + response.status + ': ' + response.statusText);
                    }
                    return response.text();
                })
                .then(content => {
                    const language = getLanguageFromPath(path);
                    const escapedContent = escapeHtml(content);

                    contentDiv.innerHTML = '<div class="h-full overflow-auto">' +
                        '<div class="text-xs text-gray-400 px-4 py-2 border-b border-gray-700 font-mono">' + path + '</div>' +
                        '<pre class="m-0 bg-gray-900 p-4" style="font-size: 13px; line-height: 1.5; min-height: calc(100% - 32px);">' +
                        '<code class="language-' + language + '" style="color: #e2e8f0;">' + escapedContent + '</code></pre>' +
                        '</div>';

                    // Apply syntax highlighting if Prism is available
                    if (typeof Prism !== 'undefined') {
                        Prism.highlightAll();
                    }
                })
                .catch(error => {
                    contentDiv.innerHTML = '<div class="flex items-center justify-center h-full text-red-400">' +
                        '<div class="text-center p-8">' +
                        '<i class="fas fa-exclamation-triangle text-3xl mb-4"></i>' +
                        '<p class="text-lg mb-2">Error loading file</p>' +
                        '<p class="text-sm text-gray-500">' + error.message + '</p>' +
                        '</div>' +
                        '</div>';
                });
        }

        function getLanguageFromPath(path) {
            const ext = path.split('.').pop().toLowerCase();
            switch(ext) {
                case 'go': return 'go';
                case 'js': return 'javascript';
                case 'json': return 'json';
                case 'yaml':
                case 'yml': return 'yaml';
                case 'toml': return 'toml';
                case 'md': return 'markdown';
                case 'html': return 'html';
                case 'css': return 'css';
                case 'sql': return 'sql';
                default: return 'text';
            }
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        // Download project from modal
        function downloadProject() {
            const form = document.getElementById('project-form');
            form.submit();
            closeExploreModal();
        }

        // Copy file content to clipboard
        function copyFileContent() {
            const codeElement = document.querySelector('#file-content code');
            if (codeElement) {
                const text = codeElement.textContent;
                navigator.clipboard.writeText(text).then(() => {
                    const btn = document.getElementById('copy-btn');
                    const originalText = btn.innerHTML;
                    btn.innerHTML = '<i class="fas fa-check mr-1"></i>Copied!';
                    btn.classList.remove('bg-blue-500', 'hover:bg-blue-600');
                    btn.classList.add('bg-green-500');

                    setTimeout(() => {
                        btn.innerHTML = originalText;
                        btn.classList.remove('bg-green-500');
                        btn.classList.add('bg-blue-500', 'hover:bg-blue-600');
                    }, 2000);
                }).catch(err => {
                    console.error('Failed to copy: ', err);
                });
            }
        }

        // Initialize tree when content is loaded via HTMX
        document.addEventListener('htmx:afterSettle', function(event) {
            if (event.detail.target.id === 'file-tree-content') {
                console.log('Initializing file tree after HTMX load');
                initializeFileTree();
                // Ensure functions are still globally accessible
                window.toggleFolder = toggleFolder;
                window.selectFile = selectFile;
            }
        });

        // Package search functionality
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

        // Ensure functions are globally accessible
        window.toggleFolder = toggleFolder;
        window.selectFile = selectFile;
    </script>
</body>
</html>`

const exploreTemplate = `
{{range .Files}}
<div class="tree-item" data-level="{{.Level}}" data-path="{{.Path}}" data-is-folder="{{.IsFolder}}">
    {{if .IsFolder}}
    <div class="folder-item flex items-center py-1 px-2 cursor-pointer hover:bg-blue-50 rounded transition duration-150"
         onclick="toggleFolder(this)" style="padding-left: {{mul .Level 20}}px;">
        <i class="fas fa-chevron-right folder-chevron text-gray-400 text-xs mr-2 transition-transform duration-200"></i>
        <i class="{{.Icon}} mr-2"></i>
        <span class="text-sm font-medium text-gray-700">{{.Name}}</span>
    </div>
    {{else}}
    <div class="file-item flex items-center py-1 px-2 cursor-pointer hover:bg-gray-100 rounded transition duration-150"
         onclick="selectFile('{{.Path}}', '{{.Name}}', this)" style="padding-left: {{add (mul .Level 20) 20}}px;">
        <i class="{{.Icon}} mr-2 text-sm"></i>
        <span class="text-sm text-gray-700">{{.Name}}</span>
    </div>
    {{end}}
</div>
{{end}}

{{if not .Files}}
<div class="text-center py-8 text-gray-500">
    <i class="fas fa-folder-open text-3xl mb-4"></i>
    <p>No files to preview</p>
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
