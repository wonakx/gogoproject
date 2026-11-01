package httpapi

import (
	"io"
	"net/http"
	"strings"

	apidoc "example.com/gogoproject/api"
	swaggerFiles "github.com/swaggo/files/v2"
)

// DocumentationHandler serves documentation only, without a database or API mutations.
func DocumentationHandler() http.Handler { return documentationMux(false) }
func documentationMux(apiEnabled bool) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		_, _ = w.Write(apidoc.Spec)
	})
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("GET /docs/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		page := docsHTML
		if !apiEnabled {
			page = strings.Replace(page, "tryItOutEnabled:false", "tryItOutEnabled:false,supportedSubmitMethods:[]", 1)
		}
		_, _ = io.WriteString(w, page)
	})
	mux.Handle("GET /docs/assets/", http.StripPrefix("/docs/assets/", http.FileServer(http.FS(swaggerFiles.FS))))

	return mux
}

const docsHTML = `<!doctype html>
<html lang="ko"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Dealer API 문서</title><link rel="stylesheet" href="/docs/assets/swagger-ui.css"></head>
<body><div id="swagger-ui"></div><script src="/docs/assets/swagger-ui-bundle.js"></script>
<script>SwaggerUIBundle({url:'/openapi.yaml',dom_id:'#swagger-ui',deepLinking:true,validatorUrl:null,displayRequestDuration:true,tryItOutEnabled:false});</script></body></html>`
