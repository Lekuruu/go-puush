package routes

import (
	"encoding/json"
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
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Response.WriteHeader(http.StatusOK)

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
		http.Error(ctx.Response, "Internal server error", http.StatusInternalServerError)
		ctx.State.Logger.Logf("Template execution error: %v", err)
	}
}

func renderErrorTemplate(title string, message string, ctx *app.Context) {
	renderTemplate(ctx, "public/response", map[string]interface{}{
		"Title":           "error",
		"ResponseTitle":   title,
		"ResponseMessage": message,
	})
}

func renderResponseTemplate(title string, message string, siteTitle string, ctx *app.Context) {
	renderTemplate(ctx, "public/response", map[string]interface{}{
		"Title":           siteTitle,
		"ResponseTitle":   title,
		"ResponseMessage": message,
	})
}

func renderRaw(status int, contentType string, data []byte, ctx *app.Context) {
	ctx.Response.Header().Set("Content-Type", contentType)
	ctx.Response.WriteHeader(status)
	ctx.Response.Write(data)
}

func renderText(status int, text string, ctx *app.Context) {
	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.WriteHeader(status)
	ctx.Response.Write([]byte(text))
}

func renderJson(status int, object any, ctx *app.Context) {
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.WriteHeader(status)
	err := json.NewEncoder(ctx.Response).Encode(object)
	if err != nil {
		ctx.State.Logger.Logf("JSON marshal error: %v", err)
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
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"iterate": func(count int) []int {
			result := make([]int, count)
			for i := range result {
				result[i] = i
			}
			return result
		},
	}

	var err error
	templates, err = template.New("").Funcs(funcs).ParseGlob("web/templates/**/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
}
