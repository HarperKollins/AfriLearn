package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/middleware"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
)

type KeyGenRequest struct {
	DeveloperName string `json:"developer_name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Tier          string `json:"tier"`
}

// GenerateAPIKey generates a new developer API key
// POST /api/v1/keys/generate
func GenerateAPIKey(c *gin.Context) {
	var req KeyGenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "developer_name and email are required fields",
		})
		return
	}

	req.DeveloperName = strings.TrimSpace(req.DeveloperName)
	req.Email = strings.TrimSpace(req.Email)
	if req.Tier == "" {
		req.Tier = "free"
	}

	// Generate random 16-byte hex key string: afr_live_<hex>
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate cryptographically secure API key",
		})
		return
	}
	apiKey := "afr_live_" + hex.EncodeToString(bytes)

	// Save to PostgreSQL
	_, err := database.DB.Exec(`
		INSERT INTO api_keys (api_key, developer_name, email, tier, is_active)
		VALUES ($1, $2, $3, $4, true)
	`, apiKey, req.DeveloperName, req.Email, req.Tier)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to save API key to database",
		})
		return
	}

	// Update in-memory cache immediately
	middleware.UpdateKeyCache(apiKey, req.DeveloperName, req.Email, req.Tier, true)

	rateLimit := "1,000 req/min"
	if req.Tier == "pro" {
		rateLimit = "50,000 req/min"
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Developer API Key generated successfully",
		Data: gin.H{
			"api_key":        apiKey,
			"developer_name": req.DeveloperName,
			"email":          req.Email,
			"tier":           req.Tier,
			"rate_limit":     rateLimit,
		},
	})
}

// ServeDeveloperPortal serves the self-service Developer Portal HTML dashboard
// GET /portal
func ServeDeveloperPortal(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AfriLearn Developer Portal — Self-Service API Keys</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-dark: #090d16;
            --card-bg: #111827;
            --accent-green: #10b981;
            --accent-cyan: #06b6d4;
            --accent-purple: #8b5cf6;
            --text-main: #f8fafc;
            --text-muted: #94a3b8;
            --border: #1e293b;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Outfit', sans-serif; }
        body { background: var(--bg-dark); color: var(--text-main); min-height: 100vh; padding: 2rem 1rem; }
        .container { max-width: 1100px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 3rem; }
        .badge { display: inline-block; background: rgba(16, 185, 129, 0.15); color: var(--accent-green); padding: 6px 16px; border-radius: 99px; font-size: 0.85rem; font-weight: 600; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 1rem; border: 1px solid rgba(16, 185, 129, 0.3); }
        .header h1 { font-size: 2.8rem; font-weight: 700; background: linear-gradient(135deg, #ffffff 0%, #10b981 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; margin-bottom: 0.75rem; }
        .header p { color: var(--text-muted); font-size: 1.15rem; max-width: 650px; margin: 0 auto; line-height: 1.6; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; margin-bottom: 3rem; }
        @media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }
        .card { background: var(--card-bg); border: 1px solid var(--border); border-radius: 16px; padding: 2rem; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5); }
        .card h2 { font-size: 1.4rem; font-weight: 600; margin-bottom: 1.25rem; color: #ffffff; display: flex; align-items: center; gap: 0.5rem; }
        .form-group { margin-bottom: 1.25rem; }
        .form-group label { display: block; font-size: 0.9rem; font-weight: 500; color: var(--text-muted); margin-bottom: 0.5rem; }
        .form-control { width: 100%; padding: 0.85rem 1rem; background: #0b0f19; border: 1px solid var(--border); border-radius: 8px; color: #ffffff; font-size: 1rem; outline: none; transition: border 0.2s; }
        .form-control:focus { border-color: var(--accent-green); }
        select.form-control { appearance: none; cursor: pointer; }
        .btn { width: 100%; padding: 0.9rem; background: linear-gradient(135deg, var(--accent-green) 0%, #059669 100%); color: #ffffff; font-weight: 600; font-size: 1rem; border: none; border-radius: 8px; cursor: pointer; transition: transform 0.1s, opacity 0.2s; }
        .btn:hover { opacity: 0.95; transform: translateY(-1px); }
        .result-box { margin-top: 1.5rem; display: none; background: rgba(16, 185, 129, 0.08); border: 1px solid rgba(16, 185, 129, 0.3); border-radius: 10px; padding: 1.25rem; }
        .key-display { font-family: 'JetBrains Mono', monospace; background: #000; padding: 0.75rem 1rem; border-radius: 6px; color: var(--accent-green); font-size: 0.95rem; margin: 0.75rem 0; word-break: break-all; display: flex; justify-content: space-between; align-items: center; }
        .copy-btn { background: #1e293b; color: #fff; border: none; padding: 4px 10px; border-radius: 4px; font-size: 0.75rem; cursor: pointer; margin-left: 10px; }
        .stats-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; margin-bottom: 2rem; }
        .stat-card { background: #0b0f19; border: 1px solid var(--border); border-radius: 12px; padding: 1.25rem; text-align: center; }
        .stat-num { font-size: 2rem; font-weight: 700; color: var(--accent-cyan); }
        .stat-label { font-size: 0.85rem; color: var(--text-muted); margin-top: 0.25rem; }
        .code-block { font-family: 'JetBrains Mono', monospace; background: #000; color: #38bdf8; padding: 1rem; border-radius: 8px; font-size: 0.85rem; overflow-x: auto; white-space: pre-wrap; margin-top: 1rem; }
        .links { text-align: center; margin-top: 3rem; color: var(--text-muted); }
        .links a { color: var(--accent-green); text-decoration: none; font-weight: 500; margin: 0 10px; }
        .links a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <span class="badge">AfriLearn Infrastructure API</span>
            <h1>Developer Portal</h1>
            <p>Generate instant self-service API keys to access 36 African curriculum datasets, NUC University Degrees, and AI Tutor LLM System Prompt generators.</p>
        </div>

        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-num">36</div>
                <div class="stat-label">Live Curricula Datasets</div>
            </div>
            <div class="stat-card">
                <div class="stat-num">22</div>
                <div class="stat-label">Institutions & Universities</div>
            </div>
            <div class="stat-card">
                <div class="stat-num">46</div>
                <div class="stat-label">Subjects & Degree Programs</div>
            </div>
        </div>

        <div class="grid">
            <!-- Key Generation Form -->
            <div class="card">
                <h2>🔑 Generate API Key</h2>
                <form id="keyForm">
                    <div class="form-group">
                        <label>Developer Name / Organization</label>
                        <input type="text" id="devName" class="form-control" placeholder="e.g. Acme EdTech Inc" required>
                    </div>
                    <div class="form-group">
                        <label>Developer Email</label>
                        <input type="email" id="devEmail" class="form-control" placeholder="e.g. dev@acme.com" required>
                    </div>
                    <div class="form-group">
                        <label>Access Tier</label>
                        <select id="devTier" class="form-control">
                            <option value="free">Free Tier (1,000 req/min)</option>
                            <option value="pro">Pro Partner Tier (50,000 req/min)</option>
                        </select>
                    </div>
                    <button type="submit" class="btn">Generate My API Key</button>
                </form>

                <div id="resultBox" class="result-box">
                    <p style="font-size: 0.9rem; font-weight: 600; color: #fff;">🎉 Your Developer API Key is Ready!</p>
                    <div class="key-display">
                        <span id="keySpan">afr_live_...</span>
                        <button class="copy-btn" onclick="copyKey()">Copy</button>
                    </div>
                    <p style="font-size: 0.8rem; color: var(--text-muted);" id="tierInfo">Tier: Free (1,000 req/min)</p>
                </div>
            </div>

            <!-- Quick Code Example -->
            <div class="card">
                <h2>⚡ Quick Integration Code</h2>
                <p style="font-size: 0.9rem; color: var(--text-muted);">Use your API key in the <code>X-API-Key</code> HTTP header:</p>
                <div class="code-block" id="curlBlock"># 🤖 Get AI Tutor LLM System Prompt
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/waec/physics/llm-prompt

# 🎓 Get University Degree Curriculum (NUC Computer Science)
curl -H "X-API-Key: afr_live_demo_9f8e2b7a" \
  http://localhost:8080/api/v1/curriculum/nuc/computer-science</div>
            </div>
        </div>

        <div class="links">
            <a href="/docs" target="_blank">📖 Interactive Swagger Docs</a> |
            <a href="/health" target="_blank">💚 API Health Status</a> |
            <a href="https://github.com/HarperKollins/AfriLearn" target="_blank">🐙 GitHub Repository</a>
        </div>
    </div>

    <script>
        document.getElementById('keyForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const devName = document.getElementById('devName').value;
            const devEmail = document.getElementById('devEmail').value;
            const devTier = document.getElementById('devTier').value;

            try {
                const res = await fetch('/api/v1/keys/generate', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ developer_name: devName, email: devEmail, tier: devTier })
                });
                const json = await res.json();

                if (json.success) {
                    const apiKey = json.data.api_key;
                    document.getElementById('keySpan').innerText = apiKey;
                    document.getElementById('tierInfo').innerText = 'Tier: ' + json.data.tier.toUpperCase() + ' (' + json.data.rate_limit + ')';
                    document.getElementById('resultBox').style.display = 'block';

                    // Update curl snippet
                    document.getElementById('curlBlock').innerText = 
'# 🤖 Get AI Tutor LLM System Prompt\n' +
'curl -H "X-API-Key: ' + apiKey + '" \\\n' +
'  http://localhost:8080/api/v1/curriculum/waec/physics/llm-prompt\n\n' +
'# 🎓 Get University Degree Curriculum\n' +
'curl -H "X-API-Key: ' + apiKey + '" \\\n' +
'  http://localhost:8080/api/v1/curriculum/nuc/computer-science';
                } else {
                    alert('Failed: ' + json.message);
                }
            } catch (err) {
                alert('Error generating API Key: ' + err.message);
            }
        });

        function copyKey() {
            const keyText = document.getElementById('keySpan').innerText;
            navigator.clipboard.writeText(keyText);
            alert('API Key copied to clipboard!');
        }
    </script>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
