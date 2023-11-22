package templates

import (
	"html/template"
	"net/http"
)

var templates *template.Template

// LoadTemplates parses the templates stored in the web directory into a public variable.
func LoadTemplates() *template.Template {
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
	return templates
}

// renderTemplate renders the specified HTML template with data
func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	err := templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
