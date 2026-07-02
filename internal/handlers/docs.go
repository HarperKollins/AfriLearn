package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServeOpenAPISpec returns the complete OpenAPI 3.0 specification in JSON format
func ServeOpenAPISpec(c *gin.Context) {
	spec := gin.H{
		"openapi": "3.0.3",
		"info": gin.H{
			"title":       "AfriLearn Curriculum API",
			"description": "The foundational data layer for African educational technology, AI Tutors, universities, polytechnics, and schools.",
			"version":     "1.0.0",
			"contact": gin.H{
				"name":  "AfriLearn Engineering Team",
				"email": "engineering@afrilearn.org",
				"url":   "https://github.com/HarperKollins/AfriLearn",
			},
		},
		"servers": []gin.H{
			{
				"url":         "http://localhost:8080",
				"description": "Local Development Server",
			},
		},
		"components": gin.H{
			"securitySchemes": gin.H{
				"ApiKeyAuth": gin.H{
					"type":        "apiKey",
					"in":          "header",
					"name":        "X-API-Key",
					"description": "Developer API Key (Use 'afr_live_demo_9f8e2b7a' for free tier testing)",
				},
			},
		},
		"security": []gin.H{
			{"ApiKeyAuth": []string{}},
		},
		"paths": gin.H{
			"/health": gin.H{
				"get": gin.H{
					"summary":     "Health Check",
					"description": "Check API status and database connection health",
					"responses": gin.H{
						"200": gin.H{"description": "API is healthy"},
					},
				},
			},
			"/api/v1/subjects": gin.H{
				"get": gin.H{
					"summary":     "List All Subjects & Degree Programs",
					"description": "Get all 46 subjects, university degree programs, and polytechnic diplomas across all 17 NUC disciplines",
					"responses": gin.H{
						"200": gin.H{"description": "List of subjects"},
					},
				},
			},
			"/api/v1/exam-boards": gin.H{
				"get": gin.H{
					"summary":     "List All Exam Boards & Institutions",
					"description": "Get all 22 registered examination boards, regulatory bodies, polytechnics, and universities",
					"responses": gin.H{
						"200": gin.H{"description": "List of exam boards and institutions"},
					},
				},
			},
			"/api/v1/curriculum/{board}/{subject}": gin.H{
				"get": gin.H{
					"summary":     "Get Full Curriculum Tree",
					"description": "Fetch complete 3-level tree (Topics -> Subtopics -> Learning Objectives) for any board/subject slug pair",
					"parameters": []gin.H{
						{
							"name":        "board",
							"in":          "path",
							"required":    true,
							"description": "Exam Board or Institution slug (e.g. waec, jamb, bece, nuc, yabatech, unilag, futo, ebsu, funai, unec, unn)",
							"schema":      gin.H{"type": "string", "example": "waec"},
						},
						{
							"name":        "subject",
							"in":          "path",
							"required":    true,
							"description": "Subject or Degree slug (e.g. physics, mathematics, computer-science, medicine-and-surgery, law, computer-engineering-tech)",
							"schema":      gin.H{"type": "string", "example": "physics"},
						},
					},
					"responses": gin.H{
						"200": gin.H{"description": "Full curriculum tree response"},
						"404": gin.H{"description": "Curriculum not found"},
					},
				},
			},
			"/api/v1/curriculum/{board}/{subject}/llm-prompt": gin.H{
				"get": gin.H{
					"summary":     "🤖 Generate AI Tutor LLM System Prompt",
					"description": "Generates structured System Directives, Context Windows, and Module Blocks formatted for OpenAI GPT-4, Claude, Gemini, and LLaMA AI tutors",
					"parameters": []gin.H{
						{
							"name":        "board",
							"in":          "path",
							"required":    true,
							"schema":      gin.H{"type": "string", "example": "waec"},
						},
						{
							"name":        "subject",
							"in":          "path",
							"required":    true,
							"schema":      gin.H{"type": "string", "example": "physics"},
						},
					},
					"responses": gin.H{
						"200": gin.H{"description": "AI Tutor LLM prompt and context window"},
					},
				},
			},
			"/api/v1/search": gin.H{
				"get": gin.H{
					"summary":     "Search Topics",
					"description": "Search topics across all secondary, university, and polytechnic curricula",
					"parameters": []gin.H{
						{
							"name":        "q",
							"in":          "query",
							"required":    true,
							"description": "Search query keyword",
							"schema":      gin.H{"type": "string", "example": "data structures"},
						},
					},
					"responses": gin.H{
						"200": gin.H{"description": "Matching topic search results"},
					},
				},
			},
		},
	}

	c.JSON(http.StatusOK, spec)
}

// ServeSwaggerUI serves the interactive Swagger UI HTML page
func ServeSwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>AfriLearn API — Interactive Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; font-family: sans-serif; }
        .swagger-ui .topbar { display: none; }
        .custom-header { background: #0f172a; color: #ffffff; padding: 20px 40px; border-bottom: 3px solid #10b981; }
        .custom-header h1 { margin: 0; font-size: 24px; font-weight: 700; color: #10b981; }
        .custom-header p { margin: 5px 0 0; color: #94a3b8; font-size: 14px; }
    </style>
</head>
<body>
    <div class="custom-header">
        <h1>🌍 AfriLearn Curriculum API Infrastructure</h1>
        <p>Interactive OpenAPI 3.0 Playground — BECE, WAEC, JAMB, NUC University Degrees & NBTE Polytechnics</p>
    </div>
    <div id="swagger-ui"></div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui-bundle.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/docs/openapi.json",
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

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
