package httpdelivery

import (
	"net/http"
	"strings"

	swaggerFiles "github.com/swaggo/files"
)

const swaggerIndexPage = `<!doctype html>
<html lang="ru">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>AnalogYouTube API</title>
    <link rel="stylesheet" href="/swagger/swagger-ui.css">
    <link rel="icon" type="image/png" href="/swagger/favicon-32x32.png" sizes="32x32">
    <link rel="icon" type="image/png" href="/swagger/favicon-16x16.png" sizes="16x16">
    <style>
        html { box-sizing: border-box; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="/swagger/swagger-ui-bundle.js"></script>
    <script src="/swagger/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function () {
            window.ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: "#swagger-ui",
                deepLinking: true,
                persistAuthorization: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

func (h *Handler) swaggerUI(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/swagger")

	switch path {
	case "", "/", "/index.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(swaggerIndexPage))
	case "/doc.json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFile(w, r, "docs/swagger.json")
	default:
		http.StripPrefix("/swagger/", http.FileServer(swaggerFiles.HTTP)).ServeHTTP(w, r)
	}
}
