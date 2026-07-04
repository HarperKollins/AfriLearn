package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/afrilearn/curriculum-api/internal/cache"
	"github.com/afrilearn/curriculum-api/internal/ingestion"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
)

// ServeAdminDashboard renders the internal management UI
// GET /admin
func ServeAdminDashboard(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AfriLearn Admin — Curriculum Management & Cache Control</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #060b14;
            --surface: #0d1526;
            --card: #111e35;
            --border: #1e3050;
            --green: #10b981;
            --cyan: #06b6d4;
            --purple: #8b5cf6;
            --orange: #f59e0b;
            --red: #ef4444;
            --text: #f0f6ff;
            --muted: #7b93b8;
            --mono: 'JetBrains Mono', monospace;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { background: var(--bg); color: var(--text); font-family: 'Outfit', sans-serif; min-height: 100vh; padding: 2rem 1.5rem; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 2rem; padding-bottom: 1rem; border-bottom: 1px solid var(--border); }
        .title { font-size: 1.8rem; font-weight: 700; background: linear-gradient(135deg, #fff 0%, var(--cyan) 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
        .badge { background: rgba(6,182,212,0.15); color: var(--cyan); border: 1px solid rgba(6,182,212,0.3); padding: 4px 12px; border-radius: 99px; font-size: 0.8rem; font-weight: 600; }
        
        .grid-4 { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1.25rem; margin-bottom: 2rem; }
        .stat-card { background: var(--card); border: 1px solid var(--border); border-radius: 12px; padding: 1.25rem; }
        .stat-val { font-size: 2.2rem; font-weight: 800; color: #fff; margin-top: 0.25rem; }
        .stat-lbl { font-size: 0.85rem; color: var(--muted); text-transform: uppercase; letter-spacing: 0.5px; }

        .actions-section { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; margin-bottom: 2rem; }
        @media (max-width: 768px) { .actions-section { grid-template-columns: 1fr; } }
        .action-card { background: var(--card); border: 1px solid var(--border); border-radius: 16px; padding: 1.5rem; }
        .action-card h3 { font-size: 1.2rem; font-weight: 600; margin-bottom: 0.5rem; display: flex; align-items: center; gap: 0.5rem; }
        .action-card p { font-size: 0.9rem; color: var(--muted); margin-bottom: 1.25rem; line-height: 1.5; }
        
        .btn { padding: 0.75rem 1.25rem; border-radius: 8px; font-weight: 600; font-size: 0.95rem; cursor: pointer; border: none; transition: transform 0.1s, opacity 0.2s; display: inline-flex; align-items: center; gap: 0.5rem; }
        .btn:hover { opacity: 0.9; transform: translateY(-1px); }
        .btn-green { background: linear-gradient(135deg, var(--green) 0%, #059669 100%); color: #fff; }
        .btn-purple { background: linear-gradient(135deg, var(--purple) 0%, #6d28d9 100%); color: #fff; }
        .btn-orange { background: linear-gradient(135deg, var(--orange) 0%, #d97706 100%); color: #fff; }
        
        .console-box { background: #000; border: 1px solid var(--border); border-radius: 12px; padding: 1.25rem; font-family: var(--mono); font-size: 0.85rem; color: var(--cyan); min-height: 160px; max-height: 300px; overflow-y: auto; white-space: pre-wrap; }
        .top-nav { display: flex; gap: 1rem; }
        .top-nav a { color: var(--muted); text-decoration: none; font-size: 0.9rem; font-weight: 500; }
        .top-nav a:hover { color: var(--cyan); }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div>
                <span class="badge">AfriLearn Infrastructure Control</span>
                <h1 class="title">Curriculum Admin & Operations</h1>
            </div>
            <div class="top-nav">
                <a href="/portal">🔑 Portal</a>
                <a href="/playground">🎮 Playground</a>
                <a href="/docs">📖 Docs</a>
            </div>
        </div>

        <!-- System Metrics Grid -->
        <div class="grid-4">
            <div class="stat-card">
                <div class="stat-lbl">Live Curricula Files</div>
                <div class="stat-val" style="color: var(--green);">56</div>
            </div>
            <div class="stat-card">
                <div class="stat-lbl">Institutions / Boards</div>
                <div class="stat-val" style="color: var(--cyan);">14</div>
            </div>
            <div class="stat-card">
                <div class="stat-lbl">Cache Total Items</div>
                <div class="stat-val" id="cacheItems" style="color: var(--purple);">0</div>
            </div>
            <div class="stat-card">
                <div class="stat-lbl">Cache Hit Ratio</div>
                <div class="stat-val" id="cacheRatio" style="color: var(--orange);">0%</div>
            </div>
        </div>

        <!-- Operations Actions -->
        <div class="actions-section">
            <div class="action-card">
                <h3>⚡ Data Ingestion Engine</h3>
                <p>Scan all 56 curriculum JSON datasets from <code>data/curricula/</code> and upsert into Neon PostgreSQL. Automatically purges the in-memory cache upon completion.</p>
                <button class="btn btn-green" onclick="runReIngestion()">Re-Ingest All Curricula</button>
                <button class="btn btn-purple" onclick="runValidation()" style="margin-left: 0.5rem;">Validate JSON Datasets</button>
            </div>
            <div class="action-card">
                <h3>🧹 Hot-Path Memory Cache</h3>
                <p>View cache metrics and clear cached API responses for <code>/curriculum/:board/:subject</code> and <code>/llm-prompt</code> endpoints.</p>
                <button class="btn btn-orange" onclick="purgeCache()">Purge Memory Cache</button>
                <button class="btn btn-purple" onclick="fetchMetrics()" style="margin-left: 0.5rem;">Refresh Metrics</button>
            </div>
        </div>

        <!-- Console Log Output -->
        <div style="margin-bottom: 0.5rem; font-weight: 600; font-size: 0.9rem; color: var(--muted);">System Output Console:</div>
        <div class="console-box" id="console">Ready. Click an action above to execute operations...</div>
    </div>

    <script>
        function log(msg) {
            const el = document.getElementById('console');
            const timestamp = new Date().toLocaleTimeString();
            el.innerText = '[' + timestamp + '] ' + msg + '\n\n' + el.innerText;
        }

        async function fetchMetrics() {
            try {
                const res = await fetch('/api/v1/admin/cache/stats');
                const json = await res.json();
                if (json.success) {
                    document.getElementById('cacheItems').innerText = json.data.total_items;
                    document.getElementById('cacheRatio').innerText = json.data.hit_ratio.toFixed(1) + '%';
                    log('Cache metrics updated: ' + json.data.total_items + ' items, ' + json.data.hits + ' hits, ' + json.data.misses + ' misses.');
                }
            } catch (e) {
                log('Error fetching metrics: ' + e.message);
            }
        }

        async function runReIngestion() {
            log('🚀 Triggering curriculum re-ingestion...');
            try {
                const res = await fetch('/api/v1/admin/reingest', { method: 'POST' });
                const json = await res.json();
                log(json.message);
                fetchMetrics();
            } catch (e) {
                log('❌ Re-ingestion failed: ' + e.message);
            }
        }

        async function runValidation() {
            log('🔍 Running curriculum data validation...');
            try {
                const res = await fetch('/api/v1/admin/validate');
                const json = await res.json();
                log(json.message + ' (Scanned ' + json.data.total_files + ' files, ' + json.data.errors_count + ' errors)');
            } catch (e) {
                log('❌ Validation check failed: ' + e.message);
            }
        }

        async function purgeCache() {
            log('🧹 Purging memory cache...');
            try {
                const res = await fetch('/api/v1/admin/cache/purge', { method: 'POST' });
                const json = await res.json();
                log(json.message);
                fetchMetrics();
            } catch (e) {
                log('❌ Cache purge failed: ' + e.message);
            }
        }

        fetchMetrics();
    </script>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// GetCacheStats returns current cache metrics
// GET /api/v1/admin/cache/stats
func GetCacheStats(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    cache.GetCache().Stats(),
	})
}

// PurgeCache clears the in-memory cache
// POST /api/v1/admin/cache/purge
func PurgeCache(c *gin.Context) {
	cache.GetCache().Clear()
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "In-memory cache purged successfully",
	})
}

// TriggerReIngestion runs the ingestion engine to re-ingest all JSON files
// POST /api/v1/admin/reingest
func TriggerReIngestion(c *gin.Context) {
	dataDir := "data/curricula"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		dataDir = filepath.Join("..", "..", "data", "curricula")
	}

	engine := ingestion.NewEngine(dataDir)
	if err := engine.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Re-ingestion failed: " + err.Error(),
		})
		return
	}

	// Purge cache after new data ingestion
	cache.GetCache().Clear()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "All 56 curriculum files re-ingested into database & cache purged successfully",
	})
}

// ValidateCurriculaDatasets runs validation across all JSON files
// GET /api/v1/admin/validate
func ValidateCurriculaDatasets(c *gin.Context) {
	dataDir := "data/curricula"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		dataDir = filepath.Join("..", "..", "data", "curricula")
	}

	engine := ingestion.NewEngine(dataDir)
	count, valErrs, err := engine.ValidateAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Validation process failed: " + err.Error(),
		})
		return
	}

	msg := "✨ All curriculum datasets passed quality validation with 0 errors!"
	if len(valErrs) > 0 {
		msg = "⚠️ Found validation warnings/errors in datasets."
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: msg,
		Data: gin.H{
			"total_files":  count,
			"errors_count": len(valErrs),
			"errors":       valErrs,
		},
	})
}
