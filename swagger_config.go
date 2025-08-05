package main

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"multiplescrape/docs"
)

// SwaggerConfig holds dynamic swagger configuration
type SwaggerConfig struct {
	Host    string   `json:"host"`
	Schemes []string `json:"schemes"`
}

// getDynamicSwaggerJSON returns swagger JSON with dynamic host
func getDynamicSwaggerJSON(c *gin.Context) string {
	// Get current host and scheme
	host := c.Request.Host
	scheme := "http"
	
	// Check for HTTPS via headers (for reverse proxy/load balancer)
	if c.Request.TLS != nil || 
	   c.GetHeader("X-Forwarded-Proto") == "https" ||
	   c.GetHeader("X-Forwarded-Ssl") == "on" ||
	   c.GetHeader("X-Url-Scheme") == "https" ||
	   strings.HasPrefix(c.Request.Header.Get("Referer"), "https://") {
		scheme = "https"
	}
	
	// Get original swagger JSON
	swaggerJSON := docs.SwaggerInfo.ReadDoc()
	
	// Parse JSON
	var swaggerDoc map[string]interface{}
	json.Unmarshal([]byte(swaggerJSON), &swaggerDoc)
	
	// Update host and schemes
	swaggerDoc["host"] = host
	swaggerDoc["schemes"] = []string{scheme}
	
	// Convert back to JSON
	updatedJSON, _ := json.Marshal(swaggerDoc)
	return string(updatedJSON)
}

// setupSwaggerRoutes configures swagger routes with dynamic host detection
func setupSwaggerRoutes(router *gin.Engine) {
	// Custom swagger.json endpoint
	router.GET("/swagger/doc.json", func(c *gin.Context) {
		swaggerJSON := getDynamicSwaggerJSON(c)
		c.Header("Content-Type", "application/json")
		c.String(200, swaggerJSON)
	})
	
	// Custom Swagger UI endpoint
	router.GET("/swagger/index.html", func(c *gin.Context) {
		// Get current host and scheme for swagger UI
		host := c.Request.Host
		scheme := "http"
		if c.Request.TLS != nil || 
		   c.GetHeader("X-Forwarded-Proto") == "https" ||
		   c.GetHeader("X-Forwarded-Ssl") == "on" ||
		   c.GetHeader("X-Url-Scheme") == "https" ||
		   strings.HasPrefix(c.Request.Header.Get("Referer"), "https://") {
			scheme = "https"
		}
		
		swaggerUIHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '` + scheme + `://` + host + `/swagger/doc.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
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
		
		c.Header("Content-Type", "text/html")
		c.String(200, swaggerUIHTML)
	})
	
	// Redirect /swagger/ to /swagger/index.html
	router.GET("/swagger/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
	
	// Redirect /swagger to /swagger/index.html
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
}