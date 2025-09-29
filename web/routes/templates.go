package routes

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Lekuruu/go-puush/internal/app"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var templates *template.Template
var printer = message.NewPrinter(language.English)

func renderTemplate(ctx *app.Context, tmpl string, pageData map[string]interface{}) {
	data := map[string]interface{}{
		"Config":  ctx.State.Config,
		"Query":   ctx.Request.URL.Query(),
		"Printer": printer,
	}
	for k, v := range pageData {
		data[k] = v
	}

	err := templates.ExecuteTemplate(ctx.Response, tmpl, data)
	if err != nil {
		// TODO: Error templates
		http.Error(ctx.Response, "Template execution error", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func init() {
	funcs := template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"sub":   func(a, b int) int { return a - b },
		"mod":   func(a, b int) int { return a % b },
		"mul":   func(a, b int) int { return a * b },
		"div":   func(a, b int) int { return a / b },
		"lower": func(s string) string { return strings.ToLower(s) },
		"upper": func(s string) string { return strings.ToUpper(s) },
	}

	var err error
	templates, err = template.New("").Funcs(funcs).ParseGlob("web/templates/**/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
}
