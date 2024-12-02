package templates

import (
    "html/template"
    "log"
    "path/filepath"
)

var Templates *template.Template

// Initialize templates by loading all .html files from the templates directory
func InitializeTemplates() {
    var err error
    // Use ParseGlob to load all HTML templates in the templates folder
    Templates, err = template.ParseGlob(filepath.Join("templates", "*.html"))
    
    if err != nil {
        log.Fatalf("Error loading templates: %v", err)
    }

    log.Println("Templates successfully loaded.")
}
