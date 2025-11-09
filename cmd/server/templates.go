package main

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=5.0, user-scalable=yes">
    <title>SYSCTL - Go Project Initializr | Generate Production-Ready Go Projects</title>
    <link rel="icon" href="/static/Group110.svg" type="image/svg+xml">
    
    <!-- SEO Meta Tags -->
    <meta name="description" content="Generate production-ready Go projects with clean architecture. Free Go project initializer supporting Gin, Echo, Fiber, Chi, and net/http. Create Go web applications with PostgreSQL, MySQL, SQLite, MongoDB, Redis. Spring Boot Initializr for Go.">
    <meta name="keywords" content="go project generator, go initializer, go project scaffold, go boilerplate, go web framework, clean architecture go, go project setup, golang generator, go project template, spring boot initializr go, go web application generator, scaffold go projects">
    <meta name="author" content="systemctl">
    <meta name="robots" content="index, follow">
    <meta name="language" content="English">
    <link rel="canonical" href="https://go-ctl.systemctl.dev/generator">
    
    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="website">
    <meta property="og:url" content="https://go-ctl.systemctl.dev/generator">
    <meta property="og:title" content="SYSCTL - Go Project Initializr | Generate Production-Ready Go Projects">
    <meta property="og:description" content="Generate production-ready Go projects with clean architecture. Free Go project initializer supporting Gin, Echo, Fiber, Chi, and net/http. Create Go web applications with PostgreSQL, MySQL, SQLite, MongoDB, Redis.">
    <meta property="og:image" content="https://go-ctl.systemctl.dev/static/Group110.svg">
    <meta property="og:site_name" content="SYSCTL Go Project Initializr">
    
    <!-- Twitter -->
    <meta property="twitter:card" content="summary_large_image">
    <meta property="twitter:url" content="https://go-ctl.systemctl.dev/generator">
    <meta property="twitter:title" content="SYSCTL - Go Project Initializr | Generate Production-Ready Go Projects">
    <meta property="twitter:description" content="Generate production-ready Go projects with clean architecture. Free Go project initializer supporting Gin, Echo, Fiber, Chi, and net/http.">
    <meta property="twitter:image" content="https://go-ctl.systemctl.dev/static/Group110.svg">
    
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Manrope:wght@200..800&display=swap" rel="stylesheet">
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-core.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js"></script>
    <style>
        * {
            font-family: 'Manrope', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
        }
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
            overflow-x: hidden;
        }
        .file-content-container {
            max-height: calc(80vh - 200px);
            overflow-y: auto;
            overflow-x: auto;
        }
        
        @media (max-width: 768px) {
            .file-tree-container {
                max-height: calc(50vh - 100px);
                font-size: 0.875rem;
            }
            .file-content-container {
                max-height: calc(50vh - 100px);
                font-size: 0.875rem;
            }
        }
        
        /* Tab Styles */
        .tab-button {
            transition: all 0.2s ease;
        }
        .tab-button.active {
            background-color: #3b82f6;
            color: white;
        }
        .tab-button:not(.active) {
            background-color: #f3f4f6;
            color: #6b7280;
        }
        .tab-button:not(.active):hover {
            background-color: #e5e7eb;
        }
        .tab-content {
            display: none;
        }
        .tab-content.active {
            display: block;
        }
        
        /* Dark Theme Overrides */
        body {
            background-color: #1a1a1a;
            color: #ffffff;
            margin: 0;
            padding: 0;
        }
        
        /* Container spacing fix */
        .container {
            max-width: 1400px;
            margin-left: auto;
            margin-right: auto;
            padding-left: clamp(0.5rem, 2vw, 1rem);
            padding-right: clamp(0.5rem, 2vw, 1rem);
            width: 100%;
        }
        
        /* Responsive Typography */
        body {
            font-size: clamp(14px, 2vw, 16px);
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
            overflow-x: hidden;
        }
        
        /* Touch-friendly interactive elements */
        button,
        input[type="button"],
        input[type="submit"],
        a.button,
        label {
            min-height: 44px;
            min-width: 44px;
            display: inline-flex;
            align-items: center;
            justify-content: center;
        }
        
        input[type="text"],
        input[type="search"],
        input[type="email"],
        select,
        textarea {
            min-height: 44px;
            font-size: 16px; /* Prevents zoom on iOS */
            padding: 0.75rem;
        }
        
        /* Responsive Tab Buttons */
        .tab-button {
            min-height: 44px;
            padding: 0.75rem 1rem;
            font-size: clamp(0.875rem, 2vw, 1rem);
            white-space: nowrap;
        }
        
        @media (max-width: 640px) {
            .tab-button {
                flex: 1;
                min-width: 0;
                padding: 0.75rem 0.5rem;
                font-size: 0.875rem;
            }
        }
        
        /* Modal Responsive Styles */
        #explore-modal {
            padding: 0;
        }
        
        @media (max-width: 768px) {
            #explore-modal {
                padding: 0;
                margin: 0;
                max-width: 100vw;
                max-height: 100vh;
                width: 100vw;
                height: 100vh;
                border-radius: 0;
            }
            
            #explore-modal > div {
                width: 100%;
                height: 100vh;
                max-width: 100vw;
                max-height: 100vh;
                border-radius: 0;
                flex-direction: column;
            }
            
            #explore-modal .file-tree-container,
            #explore-modal .file-content-container {
                max-height: calc(50vh - 100px);
            }
        }
        
        /* File Explorer Modal Layout */
        @media (max-width: 768px) {
            #explore-modal .flex.flex-1.overflow-hidden {
                flex-direction: column;
            }
            
            #explore-modal .w-1\/3 {
                width: 100%;
                border-right: none;
                border-bottom: 1px solid #404040;
                max-height: 50vh;
            }
            
            #explore-modal .flex-1.overflow-hidden:last-child {
                width: 100%;
                max-height: 50vh;
            }
        }
        
        /* Form Grid Responsive */
        @media (max-width: 640px) {
            .grid.grid-cols-1.md\:grid-cols-2,
            .grid.grid-cols-1.md\:grid-cols-2.lg\:grid-cols-3 {
                grid-template-columns: 1fr;
                gap: 0.75rem;
            }
        }
        
        /* Card Responsive Padding */
        @media (max-width: 640px) {
            .bg-white.rounded-lg.shadow-lg {
                padding: 1rem !important;
            }
        }
        
        /* Search Results Responsive */
        @media (max-width: 640px) {
            #search-results,
            #npm-search-results {
                max-height: 200px;
                font-size: 0.875rem;
            }
            
            #selected-packages,
            #selected-npm-packages {
                max-height: 150px;
                font-size: 0.875rem;
            }
        }
        
        /* Button Groups Responsive */
        @media (max-width: 640px) {
            .flex.gap-4 button,
            .flex.gap-4 a {
                width: 100%;
                margin: 0;
            }
        }
        
        /* Framework Options Responsive */
        @media (max-width: 640px) {
            .http-framework-option,
            .frontend-framework-option,
            .language-option {
                padding: 1rem 0.75rem;
                min-height: 60px;
            }
        }
        
        /* Large Screen Optimizations */
        @media (min-width: 1920px) {
            .container {
                max-width: 1600px;
            }
        }
        
        /* Landscape Mobile Adjustments */
        @media (orientation: landscape) and (max-height: 500px) {
            #explore-modal .file-tree-container,
            #explore-modal .file-content-container {
                max-height: calc(100vh - 150px);
            }
        }
        
        /* Back arrow link animation */
        .back-arrow-link {
            color: #ffffff;
            transition: all 0.3s ease;
        }
        
        .back-arrow-link:hover {
            color: #11A32B;
        }
        
        .back-arrow-svg {
            transition: transform 0.3s ease;
        }
        
        .back-arrow-link:hover .back-arrow-svg {
            transform: translateX(-4px);
        }
        
        .back-arrow-line {
            transition: stroke-dashoffset 0.3s ease;
        }
        
        .back-arrow-link:hover .back-arrow-line {
            stroke-dashoffset: 0;
        }
        
        .bg-white {
            background-color: #2d2d2d !important;
        }
        
        .text-gray-800, .text-gray-700 {
            color: #e5e5e5 !important;
        }
        
        .text-gray-600 {
            color: #b0b0b0 !important;
        }
        
        .border-gray-200, .border-gray-300 {
            border-color: #404040 !important;
        }
        
        .bg-gray-50 {
            background-color: #252525 !important;
        }
        
        input[type="text"],
        input[type="search"],
        select {
            background-color: #1a1a1a !important;
            color: #ffffff !important;
            border-color: #404040 !important;
        }
        
        input[type="text"]:focus,
        input[type="search"]:focus,
        select:focus {
            border-color: #11A32B !important;
            outline: none;
        }
        
        .tab-button:not(.active) {
            background-color: #2d2d2d !important;
            color: #b0b0b0 !important;
        }
        
        .tab-button.active {
            background-color: #11A32B !important;
            color: white !important;
        }
        
        .card-hover:hover {
            box-shadow: 0 10px 25px rgba(17, 163, 43, 0.2) !important;
        }
        
        /* Button Styles */
        button[type="submit"],
        .bg-blue-600 {
            background-color: #11A32B !important;
        }
        
        button[type="submit"]:hover,
        .bg-blue-600:hover {
            background-color: #0d8a22 !important;
        }
        
        .bg-gray-600 {
            background-color: #404040 !important;
        }
        
        .bg-gray-600:hover {
            background-color: #505050 !important;
        }
        
        /* Checkbox and Radio Styles */
        input[type="checkbox"],
        input[type="radio"] {
            accent-color: #11A32B;
        }
        
        /* Select dropdown */
        select option {
            background-color: #1a1a1a;
            color: #ffffff;
        }
        
        /* Search results and package items */
        .bg-blue-100 {
            background-color: rgba(17, 163, 43, 0.2) !important;
        }
        
        .border-blue-200 {
            border-color: #11A32B !important;
        }
        
        .text-blue-800 {
            color: #11A32B !important;
        }
        
        /* Radio button labels - unselected state */
        .bg-gray-50 {
            background-color: #252525 !important;
            border-color: #404040 !important;
        }
        
        .hover\:bg-gray-100:hover {
            background-color: #2d2d2d !important;
        }
        
        /* HTTP Framework options - unselected state */
        .http-framework-option {
            background-color: #252525 !important;
            border: 2px solid #11A32B !important;
        }
        
        .http-framework-option:hover {
            background-color: #2d2d2d !important;
        }
        
        /* HTTP Framework options - selected state */
        input[type="radio"]:checked + label.http-framework-option,
        .http-framework-option.selected {
            background-color: #11A32B !important;
            border-color: #11A32B !important;
        }
        
        input[type="radio"]:checked + label.http-framework-option .font-semibold,
        .http-framework-option.selected .font-semibold {
            color: #ffffff !important;
        }
        
        input[type="radio"]:checked + label.http-framework-option .text-sm,
        .http-framework-option.selected .text-sm {
            color: #e0e0e0 !important;
        }
        
        /* Frontend Framework options - unselected state */
        .frontend-framework-option {
            background-color: #252525 !important;
            border: 2px solid #11A32B !important;
        }
        
        .frontend-framework-option:hover {
            background-color: #2d2d2d !important;
        }
        
        /* Frontend Framework options - selected state */
        input[type="radio"]:checked + label.frontend-framework-option,
        .frontend-framework-option.selected {
            background-color: #11A32B !important;
            border-color: #11A32B !important;
        }
        
        input[type="radio"]:checked + label.frontend-framework-option .font-semibold,
        .frontend-framework-option.selected .font-semibold {
            color: #ffffff !important;
        }
        
        input[type="radio"]:checked + label.frontend-framework-option .text-sm,
        .frontend-framework-option.selected .text-sm {
            color: #e0e0e0 !important;
        }
        
        input[type="radio"]:checked + label.frontend-framework-option .fas,
        .frontend-framework-option.selected .fas {
            color: #ffffff !important;
        }
        
        /* Language options - unselected state */
        .language-option {
            background-color: #252525 !important;
            border: 2px solid #11A32B !important;
        }
        
        .language-option:hover {
            background-color: #2d2d2d !important;
        }
        
        /* Language options - selected state */
        input[type="radio"]:checked + label.language-option,
        .language-option.selected {
            background-color: #11A32B !important;
            border-color: #11A32B !important;
        }
        
        input[type="radio"]:checked + label.language-option .font-semibold,
        .language-option.selected .font-semibold {
            color: #ffffff !important;
        }
        
        input[type="radio"]:checked + label.language-option .text-sm,
        .language-option.selected .text-sm {
            color: #e0e0e0 !important;
        }
        
        input[type="radio"]:checked + label.language-option .fas,
        .language-option.selected .fas {
            color: #ffffff !important;
        }
        
        /* Radio button labels - selected state (more visible) */
        .peer-checked\:bg-blue-50 {
            background-color: #11A32B !important;
            border-color: #11A32B !important;
            border-width: 2px !important;
        }
        
        .peer-checked\:border-blue-500 {
            border-color: #11A32B !important;
            border-width: 2px !important;
        }
        
        /* Radio button check icon - make it more visible */
        .text-blue-500 {
            color: #11A32B !important;
        }
        
        /* Make selected text white for HTTP frameworks */
        .peer:checked ~ .http-framework-option .font-semibold,
        .http-framework-option.peer-checked .font-semibold {
            color: #ffffff !important;
        }
        
        .peer:checked ~ .http-framework-option .text-sm,
        .http-framework-option.peer-checked .text-sm {
            color: #e0e0e0 !important;
        }
        
        /* Search results container */
        #search-results,
        #npm-search-results {
            background-color: #2d2d2d !important;
        }
        
        /* Selected packages container */
        #selected-packages,
        #selected-npm-packages {
            background-color: #252525 !important;
            border-color: #404040 !important;
        }
        
        /* Loading indicators */
        .htmx-indicator {
            color: #11A32B !important;
        }
        
        /* Modal styles */
        .modal-backdrop {
            background-color: rgba(0, 0, 0, 0.8) !important;
        }
        
        /* File tree dark theme */
        .file-item:hover {
            background-color: rgba(17, 163, 43, 0.1) !important;
        }
        
        .file-item.selected {
            background-color: rgba(17, 163, 43, 0.2) !important;
            border-left-color: #11A32B !important;
        }
        
        /* Footer */
        footer {
            background-color: #1a1a1a !important;
            color: #b0b0b0 !important;
            border-top-color: #404040 !important;
        }
        
        footer a {
            color: #11A32B !important;
        }
        
        footer a:hover {
            color: #0d8a22 !important;
        }
    </style>
    
    <!-- Structured Data (JSON-LD) -->
    <script type="application/ld+json">
    {
        "@context": "https://schema.org",
        "@type": "WebApplication",
        "name": "SYSCTL Go Project Initializr",
        "description": "Generate production-ready Go projects with clean architecture. Free Go project initializer supporting Gin, Echo, Fiber, Chi, and net/http. Create Go web applications with PostgreSQL, MySQL, SQLite, MongoDB, Redis.",
        "url": "https://go-ctl.systemctl.dev/generator",
        "applicationCategory": "DeveloperApplication",
        "operatingSystem": "Web",
        "offers": {
            "@type": "Offer",
            "price": "0",
            "priceCurrency": "USD"
        },
        "featureList": [
            "Go project generation",
            "Clean architecture templates",
            "Multiple HTTP framework support (Gin, Echo, Fiber, Chi, net/http)",
            "Database integration (PostgreSQL, MySQL, SQLite, MongoDB, Redis)",
            "Production-ready project structure",
            "Interactive project explorer",
            "Package dependency management"
        ],
        "creator": {
            "@type": "Organization",
            "name": "systemctl",
            "url": "https://github.com/syst3mctl/go-ctl"
        }
    }
    </script>
</head>
<body class="min-h-screen">
    <!-- Header -->
    <header style="padding: 1rem 0; background-color: #1a1a1a; position: relative;">
        <div class="container mx-auto px-4">
            <div class="flex items-center justify-between">
                <div style="display: flex; flex-direction: column;">
                    <div style="width: 60px; height: 3px; background-color: #11A32B; margin-bottom: 0.5rem;"></div>
                    <div class="flex items-center gap-3">
                        <div style="width: 48px; height: 48px; display: flex; align-items: center; justify-content: center;">
                            <img src="/static/Group110.svg" alt="SYSCTL Logo" style="width: 100%; height: 100%; object-fit: contain;">
                        </div>
                    </div>
                </div>
                <a href="/" class="back-arrow-link flex items-center">
                    <svg class="back-arrow-svg" width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path class="back-arrow-line" d="M20 12H4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-dasharray="20" stroke-dashoffset="20"/>
                        <path class="back-arrow-head" d="M10 18L4 12L10 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                </a>
            </div>
        </div>
    </header>

    <div class="container mx-auto px-4 py-4">
        <div class="grid grid-cols-1 lg:grid-cols-1 gap-8">

            <!-- Left Side: Configuration Form -->
            <div class="bg-white rounded-lg shadow-lg p-6 card-hover" style="border: 1px solid #404040;">
                <h2 class="text-2xl font-bold mb-6 flex items-center" style="color: #ffffff;">
                    <i class="fas fa-cog mr-2" style="color: #11A32B;"></i>
                    Project Configuration
                </h2>

                <!-- Tab Navigation -->
                <div class="flex space-x-2 mb-6" style="border-bottom: 1px solid #404040;">
                    <button type="button" 
                            onclick="switchTab('backend')"
                            id="tab-backend"
                            class="tab-button active px-6 py-3 font-semibold rounded-t-lg">
                        <i class="fab fa-golang mr-2"></i>Back-end
                    </button>
                    <button type="button"
                            onclick="switchTab('frontend')"
                            id="tab-frontend"
                            class="tab-button px-6 py-3 font-semibold rounded-t-lg">
                        <i class="fab fa-react mr-2"></i>Front-end
                    </button>
                </div>

                <form id="project-form" action="/generate" method="POST" class="space-y-6">
                    <!-- Hidden input for project type -->
                    <input type="hidden" name="projectType" id="projectType" value="backend">

                    <!-- Backend Tab Content -->
                    <div id="tab-content-backend" class="tab-content active">
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
                    <div class="mt-6">
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
                                       class="flex p-3 rounded-lg cursor-pointer transition-all duration-200 http-framework-option"
                                       style="background-color: #252525; border: 2px solid #11A32B;">
                                    <div class="flex-1">
                                        <div class="font-semibold" style="color: #ffffff;">{{.Name}}</div>
                                        <div class="text-sm" style="color: #b0b0b0;">{{.Description}}</div>
                                    </div>
                                    <div class="flex-shrink-0 ml-2">
                                        <i class="fas fa-check-circle opacity-0 peer-checked:opacity-100" style="font-size: 1.25rem; color: #11A32B;"></i>
                                    </div>
                                </label>
                            </div>
                            {{end}}
                        </div>
                    </div>

                    <!-- Database -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
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
                    <div class="mt-6">
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
                    <div class="mt-6">
                        <label class="block text-sm font-medium text-gray-700 mb-2">
                            <i class="fas fa-search mr-1"></i>Add Dependencies
                        </label>
                        <input type="search"
                               name="q"
                               id="package-search-input"
                               class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent mb-3"
                               placeholder="Search pkg.go.dev for packages..."
                               hx-get="/fetch-packages"
                               hx-trigger="keyup changed delay:500ms, input delay:500ms"
                               hx-target="#search-results"
                               hx-swap="innerHTML"
                               hx-indicator="#search-loading"
                               oninput="handlePackageSearchInput(this)">

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
                    </div>

                    <!-- Frontend Tab Content -->
                    <div id="tab-content-frontend" class="tab-content">
                        <!-- Project Name (shared) -->
                        <div class="mb-4">
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                <i class="fas fa-tag mr-1"></i>Project Name
                            </label>
                            <input type="text"
                                   name="frontendProjectName"
                                   id="frontendProjectName"
                                   value="my-react-app"
                                   class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                   placeholder="my-react-app"
                                   required>
                        </div>

                        <!-- Framework Selection -->
                        <div class="mb-6">
                            <label class="block text-sm font-medium text-gray-700 mb-3">
                                <i class="fab fa-react mr-1"></i>Framework
                            </label>
                            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                                {{if .Options.Frontend}}
                                {{range .Options.Frontend.Frameworks}}
                                <div class="relative">
                                    <input type="radio"
                                           name="frontendFramework"
                                           value="{{.ID}}"
                                           id="framework-{{.ID}}"
                                           {{if eq .ID "react"}}checked{{end}}
                                           class="sr-only peer">
                                    <label for="framework-{{.ID}}"
                                           class="flex p-3 rounded-lg cursor-pointer transition-all duration-200 frontend-framework-option"
                                           style="background-color: #252525; border: 2px solid #11A32B;">
                                        <div class="flex-1">
                                            <div class="font-semibold" style="color: #ffffff;">{{.Name}}</div>
                                            <div class="text-sm" style="color: #b0b0b0;">{{.Description}}</div>
                                        </div>
                                        <div class="flex-shrink-0 ml-2">
                                            <i class="fas fa-check-circle opacity-0 peer-checked:opacity-100" style="color: #ffffff;"></i>
                                        </div>
                                    </label>
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                        </div>

                        <!-- Language Selection (conditional) -->
                        <div id="language-selection" class="mb-6">
                            <label class="block text-sm font-medium text-gray-700 mb-3">
                                <i class="fas fa-code mr-1"></i>Language
                            </label>
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                                {{if .Options.Frontend}}
                                {{range .Options.Frontend.Languages}}
                                <div class="relative">
                                    <input type="radio"
                                           name="frontendLanguage"
                                           value="{{.ID}}"
                                           id="lang-{{.ID}}"
                                           {{if eq .ID "typescript"}}checked{{end}}
                                           class="sr-only peer">
                                    <label for="lang-{{.ID}}"
                                           class="flex p-3 rounded-lg cursor-pointer transition-all duration-200 language-option"
                                           style="background-color: #252525; border: 2px solid #11A32B;">
                                        <div class="flex-1">
                                            <div class="font-semibold" style="color: #ffffff;">{{.Name}}</div>
                                            <div class="text-sm" style="color: #b0b0b0;">{{.Description}}</div>
                                        </div>
                                        <div class="flex-shrink-0 ml-2">
                                            <i class="fas fa-check-circle opacity-0 peer-checked:opacity-100" style="font-size: 1.25rem; color: #11A32B;"></i>
                                        </div>
                                    </label>
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                        </div>

                        <!-- Build Tool (Vite - default) -->
                        <div class="mb-6">
                            <label class="block text-sm font-medium text-gray-700 mb-3">
                                <i class="fas fa-tools mr-1"></i>Build Tool
                            </label>
                            <div class="p-3 rounded-lg" style="background-color: rgba(17, 163, 43, 0.1); border: 1px solid #11A32B;">
                                <div class="font-semibold" style="color: #ffffff;">Vite</div>
                                <div class="text-sm" style="color: #b0b0b0;">Next generation frontend tooling with fast HMR</div>
                            </div>
                            <input type="hidden" name="frontendBuildTool" value="vite">
                        </div>

                        <!-- Linter (ESLint) -->
                        <div class="mb-6">
                            <label class="block text-sm font-medium text-gray-700 mb-3">
                                <i class="fas fa-check-circle mr-1"></i>Linting
                            </label>
                            <div class="flex items-start">
                                <input type="checkbox"
                                       name="frontendLinter"
                                       value="eslint"
                                       id="linter-eslint"
                                       checked
                                       class="mt-1 mr-3 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500">
                                <div class="flex-1">
                                    <label for="linter-eslint" class="text-sm font-medium text-gray-700 cursor-pointer">
                                        ESLint
                                    </label>
                                    <p class="text-xs text-gray-600">Find and fix problems in your JavaScript/TypeScript code</p>
                                </div>
                            </div>
                        </div>

                        <!-- Additional Features -->
                        <div class="mb-6">
                            <label class="block text-sm font-medium text-gray-700 mb-3">
                                <i class="fas fa-puzzle-piece mr-1"></i>Additional Features
                            </label>
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                                {{if .Options.Frontend}}
                                {{range .Options.Frontend.Features}}
                                <div class="flex items-start">
                                    <input type="checkbox"
                                           name="frontendFeatures"
                                           value="{{.ID}}"
                                           id="frontend-feature-{{.ID}}"
                                           class="mt-1 mr-3 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500">
                                    <div class="flex-1">
                                        <label for="frontend-feature-{{.ID}}" class="text-sm font-medium text-gray-700 cursor-pointer">
                                            {{.Name}}
                                        </label>
                                        <p class="text-xs text-gray-600">{{.Description}}</p>
                                    </div>
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                        </div>

                        <!-- npm Package Search -->
                        <div class="mb-6">
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                <i class="fas fa-search mr-1"></i>Add npm Dependencies
                            </label>
                            <input type="search"
                                   name="q"
                                   class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent mb-3"
                                   placeholder="Search npm registry for packages..."
                                   hx-get="/search-npm-packages"
                                   hx-trigger="keyup changed delay:500ms"
                                   hx-target="#npm-search-results"
                                   hx-swap="innerHTML"
                                   hx-indicator="#npm-search-loading"
                                   hx-vals='{"limit": 10}'>

                            <!-- Loading indicator -->
                            <div id="npm-search-loading" class="htmx-indicator text-center py-2">
                                <i class="fas fa-spinner fa-spin text-blue-500"></i>
                                <span class="ml-2 text-sm text-gray-600">Searching...</span>
                            </div>

                            <!-- Search results container -->
                            <div id="npm-search-results" class="max-h-48 overflow-y-auto mb-4 space-y-2"></div>

                            <!-- Selected npm packages -->
                            <div>
                                <h4 class="text-sm font-medium text-gray-700 mb-2">Selected Packages:</h4>
                                <div id="selected-npm-packages" class="space-y-2 min-h-[2rem] p-2 border border-gray-200 rounded-lg bg-gray-50">
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
                    </div>
                </form>
            </div>
        </div>

        <!-- Explore Modal -->
        <div id="explore-modal" class="fixed inset-0 bg-black bg-opacity-50 modal-backdrop hidden z-50" onclick="closeExploreModal(event)">
            <div class="flex items-center justify-center min-h-screen p-0 md:p-4">
                <div class="bg-white rounded-none md:rounded-lg shadow-2xl w-full h-full md:h-[85vh] md:max-w-7xl flex flex-col" onclick="event.stopPropagation()">
                    <!-- Modal Header -->
                    <div class="flex items-center justify-between p-4 md:p-6 border-b border-gray-200">
                        <div class="flex items-center">
                            <i class="fas fa-folder-open text-xl md:text-2xl text-blue-600 mr-2 md:mr-3"></i>
                            <h2 class="text-lg md:text-2xl font-bold text-gray-800">Project Explorer</h2>
                        </div>
                        <button onclick="closeExploreModal()" class="text-gray-400 hover:text-gray-600 transition duration-200 min-h-[44px] min-w-[44px] flex items-center justify-center">
                            <i class="fas fa-times text-xl md:text-2xl"></i>
                        </button>
                    </div>

                    <!-- Modal Body -->
                    <div class="flex flex-col md:flex-row flex-1 overflow-hidden">
                        <!-- File Tree Sidebar -->
                        <div class="w-full md:w-1/3 border-r-0 md:border-r border-b md:border-b-0 border-gray-200 bg-gray-50 md:max-h-full max-h-[50vh]">
                            <div class="p-3 md:p-4 h-full flex flex-col">
                                <h3 class="text-base md:text-lg font-semibold text-gray-700 mb-2 md:mb-3 flex items-center">
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
                        <div class="flex-1 flex flex-col bg-gray-900 md:max-h-full max-h-[50vh]">
                            <!-- File Header -->
                            <div class="p-2 md:p-3 border-b border-gray-700 bg-gray-800">
                                <div class="flex items-center justify-between gap-2">
                                    <div id="current-file-header" class="flex items-center text-gray-300 min-w-0 flex-1">
                                        <i class="fas fa-file-code mr-2 flex-shrink-0"></i>
                                        <span class="text-xs md:text-sm truncate">Select a file to preview</span>
                                    </div>
                                    <div class="flex items-center space-x-2 flex-shrink-0">
                                        <button onclick="copyFileContent()" id="copy-btn" class="hidden bg-blue-600 hover:bg-blue-700 text-white px-2 md:px-3 py-1 md:py-1 rounded text-xs transition duration-200 min-h-[32px]">
                                            <i class="fas fa-copy mr-1"></i><span class="hidden sm:inline">Copy</span>
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
                    <div class="p-3 md:p-4 border-t border-gray-200 bg-gray-50 flex flex-col md:flex-row justify-between items-center gap-3">
                        <div class="text-xs md:text-sm text-gray-600 text-center md:text-left">
                            <i class="fas fa-info-circle mr-1"></i>
                            <span class="hidden sm:inline">Click files to preview content • </span>All files will be included in the download
                        </div>
                        <div class="flex gap-2 md:gap-3 w-full md:w-auto">
                            <button onclick="closeExploreModal()" class="flex-1 md:flex-none bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 md:py-2 rounded transition duration-200 min-h-[44px]">
                                Close
                            </button>
                            <button onclick="downloadProject()" class="flex-1 md:flex-none bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 md:py-2 rounded transition duration-200 min-h-[44px]">
                                <i class="fas fa-download mr-2"></i>Download Project
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Footer -->
        <footer class="mt-12 text-center border-t pt-8" style="border-top-color: #404040; color: #b0b0b0;">
            <div class="flex justify-center items-center space-x-6">
                <a href="https://github.com/syst3mctl/go-ctl" class="transition duration-200" style="color: #11A32B;">
                    <i class="fab fa-github mr-2"></i>GitHub
                </a>
                <span style="color: #404040;">•</span>
                <span style="color: #b0b0b0;">Built with ❤️ by systemctl</span>
                <span style="color: #404040;">•</span>
                <span style="color: #b0b0b0;">Powered by Go + HTMX</span>
            </div>
            <p class="mt-2 text-sm" style="color: #b0b0b0;">Generate production-ready Go projects with clean architecture</p>
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

        // Tab switching function
        function switchTab(tabName) {
            // Update tab buttons
            document.getElementById('tab-backend').classList.toggle('active', tabName === 'backend');
            document.getElementById('tab-frontend').classList.toggle('active', tabName === 'frontend');
            
            // Update tab content
            document.getElementById('tab-content-backend').classList.toggle('active', tabName === 'backend');
            document.getElementById('tab-content-frontend').classList.toggle('active', tabName === 'frontend');
            
            // Update project type hidden input
            document.getElementById('projectType').value = tabName;
        }

        // Framework change handler - shows/hides language selection
        function handleFrameworkChange(frameworkId) {
            const languageSelection = document.getElementById('language-selection');
            const languageInputs = document.querySelectorAll('input[name="frontendLanguage"]');
            
            // Angular always uses TypeScript, so hide language selection
            if (frameworkId === 'angular') {
                languageSelection.style.display = 'none';
                // Set TypeScript as the language (hidden)
                const typescriptInput = document.getElementById('lang-typescript');
                if (typescriptInput) {
                    typescriptInput.checked = true;
                }
            } else {
                // Show language selection for other frameworks
                languageSelection.style.display = 'block';
            }
        }

        // Initialize language selection visibility on page load
        document.addEventListener('DOMContentLoaded', function() {
            const selectedFramework = document.querySelector('input[name="frontendFramework"]:checked');
            if (selectedFramework) {
                handleFrameworkChange(selectedFramework.value);
            }
        });

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
                    dbSection.className = 'border rounded-lg p-4';
                    dbSection.style.backgroundColor = '#2d2d2d';
                    dbSection.style.borderColor = '#404040';
                    const savedSelection = savedSelections[database];
                    const defaultDriver = driverOptions[database][0].id;
                    const selectedDriver = savedSelection || defaultDriver;
                    
                    dbSection.innerHTML =
                        '<h4 class="font-medium mb-3 capitalize" style="color: #ffffff;">' + database + ' Driver</h4>' +
                        '<div class="space-y-2">' +
                            driverOptions[database].map(function(driver) {
                                const isChecked = driver.id === selectedDriver ? 'checked' : '';
                                return '<label class="flex items-start space-x-3 p-3 rounded-lg border cursor-pointer transition-colors duration-150" style="background-color: #252525; border-color: #404040;">' +
                                    '<input type="radio" name="driver_' + database + '" value="' + driver.id + '" ' + isChecked + ' class="mt-1" style="accent-color: #11A32B;" required>' +
                                    '<div class="flex-1 min-w-0">' +
                                        '<div class="text-sm font-medium" style="color: #ffffff;">' + driver.name + '</div>' +
                                        '<div class="text-sm" style="color: #b0b0b0;">' + driver.description + '</div>' +
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
            
            // Handle HTTP Framework selection
            function updateHttpFrameworkSelection() {
                document.querySelectorAll('input[name="httpPackage"]').forEach(radio => {
                    const label = document.querySelector('label[for="' + radio.id + '"]');
                    if (label && label.classList.contains('http-framework-option')) {
                        if (radio.checked) {
                            label.classList.add('selected');
                        } else {
                            label.classList.remove('selected');
                        }
                    }
                });
            }
            
            // Add event listeners to HTTP framework radio buttons
            document.querySelectorAll('input[name="httpPackage"]').forEach(radio => {
                radio.addEventListener('change', updateHttpFrameworkSelection);
            });
            
            // Set initial state
            updateHttpFrameworkSelection();
            
            // Handle Frontend Framework selection
            function updateFrontendFrameworkSelection() {
                document.querySelectorAll('input[name="frontendFramework"]').forEach(radio => {
                    const label = document.querySelector('label[for="' + radio.id + '"]');
                    if (label && label.classList.contains('frontend-framework-option')) {
                        if (radio.checked) {
                            label.classList.add('selected');
                        } else {
                            label.classList.remove('selected');
                        }
                    }
                });
            }
            
            // Add event listeners to Frontend framework radio buttons
            document.querySelectorAll('input[name="frontendFramework"]').forEach(radio => {
                radio.addEventListener('change', function() {
                    updateFrontendFrameworkSelection();
                    handleFrameworkChange(radio.value);
                });
            });
            
            // Set initial state
            updateFrontendFrameworkSelection();
            
            // Handle Language selection
            function updateLanguageSelection() {
                document.querySelectorAll('input[name="frontendLanguage"]').forEach(radio => {
                    const label = document.querySelector('label[for="' + radio.id + '"]');
                    if (label && label.classList.contains('language-option')) {
                        if (radio.checked) {
                            label.classList.add('selected');
                        } else {
                            label.classList.remove('selected');
                        }
                    }
                });
            }
            
            // Add event listeners to Language radio buttons
            document.querySelectorAll('input[name="frontendLanguage"]').forEach(radio => {
                radio.addEventListener('change', updateLanguageSelection);
            });
            
            // Set initial state
            updateLanguageSelection();
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
            const projectType = formData.get('projectType') || 'backend';
            params.append('projectType', projectType);
            
            if (projectType === 'frontend') {
                const framework = formData.get('frontendFramework') || 'react';
                params.append('frontendProjectName', formData.get('frontendProjectName') || 'my-react-app');
                params.append('frontendFramework', framework);
                // For Angular, always use TypeScript; otherwise use selected language or default to TypeScript
                const language = framework === 'angular' ? 'typescript' : (formData.get('frontendLanguage') || 'typescript');
                params.append('frontendLanguage', language);
                params.append('frontendBuildTool', formData.get('frontendBuildTool') || 'vite');
                params.append('frontendLinter', formData.get('frontendLinter') || 'eslint');
            } else {
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
                case 'jsx': return 'jsx';
                case 'ts': return 'typescript';
                case 'tsx': return 'tsx';
                case 'json': return 'json';
                case 'yaml':
                case 'yml': return 'yaml';
                case 'toml': return 'toml';
                case 'md': return 'markdown';
                case 'html': return 'html';
                case 'css': return 'css';
                case 'sql': return 'sql';
                case 'cjs': return 'javascript';
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

        // Package management constants
        const MAX_PACKAGES = 20;

        // Package search input handler - clears results when input is empty
        function handlePackageSearchInput(input) {
            if (input.value.trim() === '') {
                const searchResults = document.getElementById('search-results');
                if (searchResults) {
                    searchResults.innerHTML = '';
                }
            }
        }

        // Check if package is already selected
        function isPackageSelected(pkgPath) {
            const selectedPackages = document.getElementById('selected-packages');
            if (!selectedPackages) return false;
            
            const hiddenInputs = selectedPackages.querySelectorAll('input[name="customPackages"]');
            for (let input of hiddenInputs) {
                if (input.value === pkgPath) {
                    return true;
                }
            }
            return false;
        }

        // Get current package count
        function getPackageCount() {
            const selectedPackages = document.getElementById('selected-packages');
            if (!selectedPackages) return 0;
            
            const hiddenInputs = selectedPackages.querySelectorAll('input[name="customPackages"]');
            return hiddenInputs.length;
        }

        // Add package to selection
        function addPackage(pkgPath) {
            // Check for duplicates
            if (isPackageSelected(pkgPath)) {
                alert('This package is already selected.');
                return;
            }

            // Check package limit
            if (getPackageCount() >= MAX_PACKAGES) {
                alert('Maximum of ' + MAX_PACKAGES + ' packages allowed. Please remove some packages before adding more.');
                return;
            }

            // Generate unique ID for the package element
            const pkgID = pkgPath.replace(/\//g, '-').replace(/\./g, '-');

            // Create package element
            const packageDiv = document.createElement('div');
            packageDiv.id = 'pkg-' + pkgID;
            packageDiv.className = 'flex items-center justify-between bg-blue-100 border border-blue-200 rounded-lg p-2';
            packageDiv.innerHTML = 
                '<div class="flex items-center">' +
                    '<i class="fas fa-cube text-blue-600 mr-2"></i>' +
                    '<span class="text-sm font-medium text-blue-800">' + escapeHtml(pkgPath) + '</span>' +
                '</div>' +
                '<input type="hidden" name="customPackages" value="' + escapeHtml(pkgPath) + '">' +
                '<button type="button" onclick="removePackage(\'' + pkgID + '\', \'' + escapeHtml(pkgPath) + '\')" class="text-red-500 hover:text-red-700 font-bold text-sm ml-2 transition duration-150">' +
                    '<i class="fas fa-times"></i>' +
                '</button>';

            // Remove placeholder if it exists
            const selectedPackages = document.getElementById('selected-packages');
            const placeholder = selectedPackages.querySelector('p.italic');
            if (placeholder) {
                placeholder.remove();
            }

            // Add package to container
            selectedPackages.appendChild(packageDiv);
        }

        // Remove package from selection
        function removePackage(pkgID, pkgPath) {
            const packageElement = document.getElementById('pkg-' + pkgID);
            if (packageElement) {
                packageElement.remove();
                
                // Show placeholder if no packages remain
                const selectedPackages = document.getElementById('selected-packages');
                const hiddenInputs = selectedPackages.querySelectorAll('input[name="customPackages"]');
                if (hiddenInputs.length === 0) {
                    const placeholder = document.createElement('p');
                    placeholder.className = 'text-sm text-gray-500 italic';
                    placeholder.textContent = 'No packages selected';
                    selectedPackages.appendChild(placeholder);
                }
            }
        }

        // Package search functionality
        // Remove placeholder text when packages are selected (for HTMX responses)
        document.addEventListener('htmx:afterRequest', function(event) {
            // Only handle add-package requests
            if (event.detail.pathInfo.requestPath === '/add-package') {
                const selectedPackages = document.getElementById('selected-packages');
                const placeholder = selectedPackages.querySelector('p.italic');
                if (placeholder && selectedPackages.children.length > 1) {
                    placeholder.remove();
                }
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
        window.addPackage = addPackage;
        window.removePackage = removePackage;
        window.handlePackageSearchInput = handlePackageSearchInput;
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
    <div class="flex items-center justify-between p-3 border rounded-lg transition duration-150" style="background-color: #2d2d2d; border-color: #404040;">
        <div class="flex-1 min-w-0">
            <div class="font-semibold text-sm" style="color: #ffffff;">{{.Path}}</div>
            <div class="text-xs mt-1 truncate" style="color: #b0b0b0;">{{.Synopsis}}</div>
        </div>
        <button type="button"
                onclick="addPackage('{{.Path}}')"
                class="ml-3 text-white px-3 py-1 rounded text-sm font-medium transition duration-150 flex items-center"
                style="background-color: #11A32B;">
            <i class="fas fa-plus mr-1"></i>Add
        </button>
    </div>
    {{end}}
{{else}}
    <div class="text-center py-4" style="color: #b0b0b0;">
        <i class="fas fa-search text-2xl mb-2"></i>
        <p class="text-sm">No packages found for "{{.Query}}"</p>
        <p class="text-xs mt-1" style="color: #808080;">Try a different search term</p>
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
            onclick="removePackage('{{.ID}}', '{{.PkgPath}}')"
            class="text-red-500 hover:text-red-700 font-bold text-sm ml-2 transition duration-150">
        <i class="fas fa-times"></i>
    </button>
</div>
`

const npmSearchResultsTemplate = `
{{if .Results}}
    {{range .Results}}
    <div class="flex items-center justify-between p-3 border rounded-lg transition duration-150" style="background-color: #2d2d2d; border-color: #404040;">
        <div class="flex-1 min-w-0">
            <div class="font-semibold text-sm" style="color: #ffffff;">{{.Name}}</div>
            <div class="text-xs mt-1 truncate" style="color: #b0b0b0;">{{.Description}}</div>
            <div class="flex items-center gap-3 mt-1 text-xs" style="color: #808080;">
                <span>v{{.Version}}</span>
                {{if .Downloads.Monthly}}
                <span>• <i class="fas fa-download mr-1"></i>{{.Downloads.Monthly}} downloads/month</span>
                {{end}}
            </div>
        </div>
        <button type="button"
                hx-post="/add-npm-package"
                hx-vals='{"pkgName": "{{.Name}}"}'
                hx-target="#selected-npm-packages"
                hx-swap="beforeend"
                class="ml-3 text-white px-3 py-1 rounded text-sm font-medium transition duration-150 flex items-center"
                style="background-color: #11A32B;">
            <i class="fas fa-plus mr-1"></i>Add
        </button>
    </div>
    {{end}}
{{else}}
    <div class="text-center py-4" style="color: #b0b0b0;">
        <i class="fas fa-search text-2xl mb-2"></i>
        <p class="text-sm">No packages found for "{{.Query}}"</p>
        <p class="text-xs mt-1" style="color: #808080;">Try a different search term</p>
    </div>
{{end}}
`

const selectedNpmPackageTemplate = `
<div id="npm-pkg-{{.ID}}" class="flex items-center justify-between bg-blue-100 border border-blue-200 rounded-lg p-2">
    <div class="flex items-center">
        <i class="fab fa-npm text-blue-600 mr-2"></i>
        <span class="text-sm font-medium text-blue-800">{{.PkgName}}</span>
    </div>

    <!-- Hidden input for form submission -->
    <input type="hidden" name="npmPackages" value="{{.PkgName}}">

    <button type="button"
            hx-target="#npm-pkg-{{.ID}}"
            hx-swap="delete"
            class="text-red-500 hover:text-red-700 font-bold text-sm ml-2 transition duration-150">
        <i class="fas fa-times"></i>
    </button>
</div>
`

const landingTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=5.0, user-scalable=yes">
    <title>SYSCTL - Go Project Initializer | Generate Production-Ready Go Projects Instantly</title>
    <link rel="icon" href="/static/Group110.svg" type="image/svg+xml">
    
    <!-- SEO Meta Tags -->
    <meta name="description" content="SYSCTL Go Project Initializer - Generate production-ready Go projects instantly. Free tool to scaffold Go web applications with clean architecture, multiple HTTP frameworks, and database support. From idea to launch in seconds.">
    <meta name="keywords" content="go project generator, go initializer, go project scaffold, go boilerplate, go web framework, clean architecture go, go project setup, golang generator, go project template, spring boot initializr go, go web application generator, scaffold go projects, free go initializer">
    <meta name="author" content="systemctl">
    <meta name="robots" content="index, follow">
    <meta name="language" content="English">
    <link rel="canonical" href="https://go-ctl.systemctl.dev/">
    
    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="website">
    <meta property="og:url" content="https://go-ctl.systemctl.dev/">
    <meta property="og:title" content="SYSCTL - Go Project Initializer | Generate Production-Ready Go Projects Instantly">
    <meta property="og:description" content="SYSCTL Go Project Initializer - Generate production-ready Go projects instantly. Free tool to scaffold Go web applications with clean architecture, multiple HTTP frameworks, and database support. From idea to launch in seconds.">
    <meta property="og:image" content="https://go-ctl.systemctl.dev/static/Group110.svg">
    <meta property="og:site_name" content="SYSCTL Go Project Initializer">
    
    <!-- Twitter -->
    <meta property="twitter:card" content="summary_large_image">
    <meta property="twitter:url" content="https://go-ctl.systemctl.dev/">
    <meta property="twitter:title" content="SYSCTL - Go Project Initializer | Generate Production-Ready Go Projects Instantly">
    <meta property="twitter:description" content="SYSCTL Go Project Initializer - Generate production-ready Go projects instantly. Free tool to scaffold Go web applications with clean architecture.">
    <meta property="twitter:image" content="https://go-ctl.systemctl.dev/static/Group110.svg">
    
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Manrope:wght@200..800&display=swap" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Manrope', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background-color: #1a1a1a;
            color: #ffffff;
            min-height: 100vh;
            height: 100vh;
            overflow-x: hidden;
            overflow-y: auto;
            margin: 0;
            padding: 0;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
        }
        
        html {
            scroll-behavior: smooth;
            overflow-x: hidden;
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 0 clamp(0.5rem, 2vw, 1rem);
            width: 100%;
        }
        
        /* Header */
        header {
            padding: clamp(1rem, 2vw, 1.5rem) 0;
            position: relative;
        }
        
        .header-content {
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 1rem;
        }
        
        .logo-section {
            display: flex;
            flex-direction: column;
        }
        
        .green-accent-line {
            width: clamp(40px, 8vw, 60px);
            height: 3px;
            background-color: #11A32B;
            margin-bottom: 0.5rem;
        }
        
        .logo {
            display: flex;
            align-items: center;
            gap: clamp(0.5rem, 2vw, 1rem);
        }
        
        .logo-icon {
            width: clamp(36px, 8vw, 48px);
            height: clamp(36px, 8vw, 48px);
            min-width: 36px;
            min-height: 36px;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        
        .logo-icon img {
            width: 100%;
            height: 100%;
            object-fit: contain;
        }
        
        .logo-text {
            font-size: clamp(1.25rem, 4vw, 2rem);
            font-weight: 700;
            letter-spacing: -0.5px;
        }
        
        .nav-links {
            display: flex;
            gap: clamp(0.75rem, 3vw, 1.5rem);
            align-items: center;
            font-size: clamp(0.8rem, 2vw, 0.9rem);
            color: #ffffff;
            flex-wrap: wrap;
        }
        
        .nav-links a {
            color: #ffffff;
            text-decoration: none;
            transition: opacity 0.2s;
            padding: 0.5rem;
            min-height: 44px;
            display: inline-flex;
            align-items: center;
        }
        
        .nav-links a:hover {
            opacity: 0.7;
        }
        
        .nav-links span {
            display: none;
        }
        
        @media (min-width: 480px) {
            .nav-links span {
                display: inline;
            }
        }
        
        /* Main Content Wrapper */
        .main-wrapper {
            min-height: calc(100vh - clamp(80px, 15vw, 120px));
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            overflow-x: hidden;
        }
        
        /* Hero Section */
        .hero {
            flex: 1;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            padding: clamp(1rem, 4vw, 2rem) 0;
            position: relative;
            gap: clamp(1.5rem, 5vw, 3rem);
        }
        
        @media (min-width: 1024px) {
            .hero {
                flex-direction: row;
                justify-content: space-between;
                align-items: center;
            }
        }
        
        .hero-left {
            flex: 1;
            width: 100%;
            max-width: 100%;
            text-align: center;
        }
        
        @media (min-width: 1024px) {
            .hero-left {
                max-width: 60%;
                text-align: left;
            }
        }
        
        .hero-heading {
            font-size: clamp(2rem, 8vw, 4rem);
            font-weight: 700;
            line-height: 1.1;
            margin-bottom: clamp(1rem, 3vw, 1.5rem);
            letter-spacing: clamp(-1px, -0.5vw, -2px);
        }
        
        .hero-right {
            flex: 1;
            width: 100%;
            max-width: 100%;
            text-align: center;
        }
        
        @media (min-width: 1024px) {
            .hero-right {
                max-width: 40%;
                text-align: right;
            }
        }
        
        .hero-subheading {
            font-size: clamp(2rem, 8vw, 4rem);
            font-weight: 700;
            line-height: 1.1;
            letter-spacing: clamp(-1px, -0.5vw, -2px);
        }
        
        /* Green Accent Bar */
        .accent-bar {
            background-color: #11A32B;
            width: 100%;
            padding: clamp(0.75rem, 2vw, 1rem) 0;
            margin: clamp(1rem, 3vw, 1.5rem) 0;
            overflow: hidden;
            position: relative;
            flex-shrink: 0;
        }
        
        .accent-bar-content {
            display: flex;
            gap: clamp(1.5rem, 5vw, 3rem);
            white-space: nowrap;
            animation: scroll 30s linear infinite;
        }
        
        .accent-bar-item {
            font-size: clamp(0.75rem, 2.5vw, 1.2rem);
            font-weight: 600;
            color: #ffffff;
            text-transform: uppercase;
            letter-spacing: clamp(1px, 0.5vw, 2px);
            opacity: 0.9;
            flex-shrink: 0;
        }
        
        @media (prefers-reduced-motion: reduce) {
            .accent-bar-content {
                animation: none;
            }
        }
        
        @keyframes scroll {
            0% {
                transform: translateX(0);
            }
            100% {
                transform: translateX(-50%);
            }
        }
        
        /* Stats Section */
        .stats-section {
            display: flex;
            flex-wrap: wrap;
            gap: clamp(1rem, 4vw, 2rem);
            margin: clamp(1rem, 3vw, 1.5rem) 0;
            font-size: clamp(0.9rem, 2vw, 1rem);
            justify-content: center;
        }
        
        @media (min-width: 1024px) {
            .stats-section {
                justify-content: flex-start;
            }
        }
        
        .stat-item {
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
            min-width: 120px;
        }
        
        .stat-number {
            font-size: clamp(1.5rem, 5vw, 2rem);
            font-weight: 700;
            color: #11A32B;
        }
        
        .stat-label {
            font-size: clamp(0.7rem, 1.8vw, 0.85rem);
            opacity: 0.8;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        /* CTA Buttons */
        .cta-section {
            display: flex;
            flex-direction: column;
            gap: clamp(0.75rem, 3vw, 1.5rem);
            margin-top: clamp(1rem, 3vw, 1.5rem);
            flex-shrink: 0;
            width: 100%;
        }
        
        @media (min-width: 640px) {
            .cta-section {
                flex-direction: row;
            }
        }
        
        .cta-button {
            padding: clamp(0.75rem, 2vw, 0.875rem) clamp(1.5rem, 4vw, 2rem);
            border: 2px solid #ffffff;
            background-color: transparent;
            color: #ffffff;
            font-size: clamp(0.85rem, 2vw, 0.95rem);
            font-weight: 600;
            text-decoration: none;
            transition: all 0.3s ease;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            gap: 0.5rem;
            cursor: pointer;
            min-height: 44px;
            width: 100%;
            text-align: center;
        }
        
        @media (min-width: 640px) {
            .cta-button {
                width: auto;
            }
        }
        
        .cta-button:hover {
            background-color: #ffffff;
            color: #1a1a1a;
            transform: translateY(-2px);
        }
        
        .cta-button:active {
            transform: translateY(0);
        }
        
        .cta-button.primary {
            background-color: #11A32B;
            border-color: #11A32B;
        }
        
        .cta-button.primary:hover {
            background-color: #0d8a22;
            border-color: #0d8a22;
        }
        
        /* Responsive Breakpoints */
        @media (max-width: 320px) {
            .container {
                padding: 0 0.5rem;
            }
            
            .hero-heading,
            .hero-subheading {
                font-size: 1.75rem;
            }
            
            .stats-section {
                flex-direction: column;
                align-items: center;
            }
        }
        
        @media (min-width: 480px) and (max-width: 639px) {
            .hero-heading,
            .hero-subheading {
                font-size: clamp(2.25rem, 6vw, 3rem);
            }
        }
        
        @media (min-width: 640px) and (max-width: 767px) {
            .hero-heading,
            .hero-subheading {
                font-size: clamp(2.5rem, 7vw, 3.5rem);
            }
        }
        
        @media (min-width: 768px) and (max-width: 1023px) {
            .hero-heading,
            .hero-subheading {
                font-size: clamp(3rem, 8vw, 3.75rem);
            }
        }
        
        @media (min-width: 1280px) {
            .container {
                max-width: 1200px;
            }
        }
        
        @media (min-width: 1536px) {
            .container {
                max-width: 1400px;
            }
        }
        
        @media (min-width: 1920px) {
            .container {
                max-width: 1600px;
            }
        }
        
        /* Landscape orientation adjustments */
        @media (orientation: landscape) and (max-height: 500px) {
            .main-wrapper {
                min-height: auto;
            }
            
            .hero {
                padding: 1rem 0;
            }
            
            .accent-bar {
                margin: 0.5rem 0;
                padding: 0.5rem 0;
            }
        }
    </style>
    
    <!-- Structured Data (JSON-LD) -->
    <script type="application/ld+json">
    {
        "@context": "https://schema.org",
        "@type": "SoftwareApplication",
        "name": "SYSCTL Go Project Initializer",
        "description": "SYSCTL Go Project Initializer - Generate production-ready Go projects instantly. Free tool to scaffold Go web applications with clean architecture, multiple HTTP frameworks, and database support. From idea to launch in seconds.",
        "url": "https://go-ctl.systemctl.dev/",
        "applicationCategory": "DeveloperApplication",
        "operatingSystem": "Web",
        "offers": {
            "@type": "Offer",
            "price": "0",
            "priceCurrency": "USD"
        },
        "featureList": [
            "Go project generation",
            "Clean architecture templates",
            "Multiple HTTP framework support (Gin, Echo, Fiber, Chi, net/http)",
            "Database integration (PostgreSQL, MySQL, SQLite, MongoDB, Redis)",
            "Production-ready project structure",
            "Interactive project explorer",
            "Package dependency management",
            "Frontend framework support (React, Angular, Vue, Svelte, Solid)"
        ],
        "creator": {
            "@type": "Organization",
            "name": "systemctl",
            "url": "https://github.com/syst3mctl/go-ctl"
        },
        "aggregateRating": {
            "@type": "AggregateRating",
            "ratingValue": "5",
            "ratingCount": "1"
        }
    }
    </script>
</head>
<body>
    <!-- Header -->
    <header>
        <div class="container">
            <div class="header-content">
                <div class="logo-section">
                    <div class="green-accent-line"></div>
                    <div class="logo">
                        <div class="logo-icon">
                            <img src="/static/Group110.svg" alt="SYSCTL Logo">
                        </div>
                    </div>
                </div>
                <div class="nav-links">
                    <a href="/generator">Generator Page</a>
                    <span>|</span>
                    <a href="#">English</a>
                </div>
            </div>
        </div>
    </header>

    <!-- Main Content -->
    <div class="main-wrapper">
        <!-- Hero Section -->
        <div class="container">
            <div class="hero">
                <div class="hero-left">
                    <h1 class="hero-heading">Short, Punchy, Solid</h1>
                    
                    <!-- Stats Section -->
                    <div class="stats-section">
                        <div class="stat-item">
                            <div class="stat-number">{{.TotalGenerations}}</div>
                            <div class="stat-label">Projects Generated</div>
                        </div>
                        <div class="stat-item">
                            <div class="stat-number">{{.TotalDownloads}}</div>
                            <div class="stat-label">Downloads</div>
                        </div>
                    </div>
                    
                    <!-- CTA Section -->
                    <div class="cta-section">
                        <a href="/generator" class="cta-button">
                            Try out →
                        </a>
                        <a href="https://github.com/syst3mctl/go-ctl" class="cta-button">
                            See Our Work
                        </a>
                    </div>
                </div>
                
                <div class="hero-right">
                    <h2 class="hero-subheading">From Idea<br>to Launch</h2>
                </div>
            </div>
        </div>

        <!-- Green Accent Bar -->
        <div class="accent-bar">
            <div class="accent-bar-content">
                <span class="accent-bar-item">Go Projects</span>
                <span class="accent-bar-item">Javascript & Typescript</span>
                <span class="accent-bar-item">Clean Architecture</span>
                <span class="accent-bar-item">Production Ready</span>
                <span class="accent-bar-item">Fast Setup</span>
                <span class="accent-bar-item">Javascript & Typescript</span>
                <span class="accent-bar-item">Go Projects</span>
                <span class="accent-bar-item">React Apps</span>
                <span class="accent-bar-item">Clean Architecture</span>
                <span class="accent-bar-item">Production Ready</span>
                <span class="accent-bar-item">Fast Setup</span>
                <span class="accent-bar-item">Best Practices</span>
            </div>
        </div>
    </div>


    <script>
        // Smooth animations and interactions
        document.addEventListener('DOMContentLoaded', () => {
            // Animate stats on load - fast animation (under 1 second)
            const statNumbers = document.querySelectorAll('.stat-number');
            statNumbers.forEach(stat => {
                const finalValue = parseInt(stat.textContent);
                if (!isNaN(finalValue) && finalValue > 0) {
                    let current = 0;
                    const duration = 800; // Total animation duration in ms
                    const steps = 20; // Number of animation steps
                    const increment = finalValue / steps;
                    const interval = duration / steps;
                    
                    const timer = setInterval(() => {
                        current += increment;
                        if (current >= finalValue) {
                            stat.textContent = finalValue;
                            clearInterval(timer);
                        } else {
                            stat.textContent = Math.floor(current);
                        }
                    }, interval);
                }
            });
        });
        
        // Smooth button hover effects
        const buttons = document.querySelectorAll('.cta-button');
        buttons.forEach(button => {
            button.addEventListener('mouseenter', () => {
                button.style.transform = 'translateY(-2px)';
            });
            button.addEventListener('mouseleave', () => {
                button.style.transform = 'translateY(0)';
            });
        });
    </script>
</body>
</html>
`
