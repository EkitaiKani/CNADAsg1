package templates

import (
    "html/template"
    "log"
    "path/filepath"
)

var Templates *template.Template

// Initialize templates by loading all .html files from the templates directory
func InitializeTemplates() {
    // Create a function map with custom functions
    funcMap := template.FuncMap{
        "add": func(a, b int) int {
            return a + b
        },
    }

    var err error
    // Use ParseGlob with Funcs to load all HTML templates and add custom functions
    Templates, err = template.New("").Funcs(funcMap).ParseGlob(filepath.Join("templates", "*.html"))
    
    if err != nil {
        log.Fatalf("Error loading templates: %v", err)
    }
    log.Println("Templates successfully loaded.")
}