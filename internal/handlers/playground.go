package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServePlayground serves the full interactive AfriLearn AI Playground
// GET /playground
func ServePlayground(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AfriLearn AI Playground — Test Everything</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #060b14;
            --surface: #0d1526;
            --card: #111e35;
            --border: #1e3050;
            --green: #10b981;
            --green2: #059669;
            --cyan: #06b6d4;
            --purple: #8b5cf6;
            --orange: #f59e0b;
            --red: #ef4444;
            --text: #f0f6ff;
            --muted: #7b93b8;
            --mono: 'JetBrains Mono', monospace;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { background: var(--bg); color: var(--text); font-family: 'Outfit', sans-serif; min-height: 100vh; }

        /* ── Layout ── */
        .app { display: grid; grid-template-columns: 320px 1fr; grid-template-rows: 60px 1fr; height: 100vh; }
        .topbar { grid-column: 1 / -1; background: var(--surface); border-bottom: 1px solid var(--border); display: flex; align-items: center; padding: 0 1.5rem; gap: 1rem; }
        .sidebar { background: var(--surface); border-right: 1px solid var(--border); overflow-y: auto; padding: 1rem; display: flex; flex-direction: column; gap: 0.75rem; }
        .main { display: flex; flex-direction: column; overflow: hidden; }

        /* ── Topbar ── */
        .logo { font-size: 1.1rem; font-weight: 700; background: linear-gradient(135deg, #fff 0%, var(--green) 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
        .badge-live { background: rgba(16,185,129,0.15); color: var(--green); padding: 3px 10px; border-radius: 99px; font-size: 0.75rem; font-weight: 600; border: 1px solid rgba(16,185,129,0.3); }
        .topbar-links { margin-left: auto; display: flex; gap: 0.75rem; }
        .topbar-links a { color: var(--muted); text-decoration: none; font-size: 0.85rem; }
        .topbar-links a:hover { color: var(--green); }

        /* ── Sidebar ── */
        .section-label { font-size: 0.7rem; font-weight: 700; text-transform: uppercase; letter-spacing: 1.5px; color: var(--muted); margin-top: 0.5rem; }
        .config-group { background: var(--card); border: 1px solid var(--border); border-radius: 10px; padding: 1rem; }
        .config-group label { font-size: 0.8rem; font-weight: 500; color: var(--muted); display: block; margin-bottom: 0.35rem; }
        .config-group input, .config-group select { width: 100%; background: var(--bg); border: 1px solid var(--border); border-radius: 6px; color: var(--text); padding: 0.5rem 0.75rem; font-size: 0.85rem; outline: none; transition: border 0.2s; margin-bottom: 0.6rem; }
        .config-group input:focus, .config-group select:focus { border-color: var(--green); }
        .config-group input:last-child, .config-group select:last-child { margin-bottom: 0; }

        .nav-btn { width: 100%; text-align: left; padding: 0.6rem 0.85rem; background: transparent; border: 1px solid transparent; border-radius: 8px; color: var(--muted); font-family: inherit; font-size: 0.85rem; cursor: pointer; transition: all 0.15s; display: flex; align-items: center; gap: 0.6rem; }
        .nav-btn:hover { background: var(--card); color: var(--text); border-color: var(--border); }
        .nav-btn.active { background: rgba(16,185,129,0.1); color: var(--green); border-color: rgba(16,185,129,0.3); font-weight: 600; }
        .nav-btn .method { font-size: 0.65rem; font-weight: 700; padding: 1px 5px; border-radius: 3px; font-family: var(--mono); }
        .get { background: rgba(6,182,212,0.2); color: var(--cyan); }
        .post { background: rgba(139,92,246,0.2); color: var(--purple); }

        /* ── Main panels ── */
        .panel { display: none; flex-direction: column; height: 100%; }
        .panel.active { display: flex; }

        /* ── Chat Panel ── */
        .chat-messages { flex: 1; overflow-y: auto; padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
        .msg { max-width: 78%; animation: fadeUp 0.3s ease; }
        .msg.user { align-self: flex-end; }
        .msg.assistant { align-self: flex-start; }
        .msg.system { align-self: center; max-width: 100%; }
        @keyframes fadeUp { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

        .bubble { padding: 0.85rem 1.1rem; border-radius: 14px; font-size: 0.9rem; line-height: 1.6; }
        .msg.user .bubble { background: linear-gradient(135deg, var(--green) 0%, var(--green2) 100%); color: #fff; border-bottom-right-radius: 4px; }
        .msg.assistant .bubble { background: var(--card); border: 1px solid var(--border); border-bottom-left-radius: 4px; color: var(--text); }
        .msg.system .bubble { background: rgba(245,158,11,0.1); border: 1px solid rgba(245,158,11,0.3); color: var(--orange); font-size: 0.8rem; border-radius: 8px; text-align: center; }

        .meta-pill { font-size: 0.7rem; color: var(--muted); margin-top: 4px; }
        .msg.user .meta-pill { text-align: right; }

        .step-chips { display: flex; gap: 0.4rem; flex-wrap: wrap; margin-top: 0.5rem; }
        .chip { font-size: 0.7rem; padding: 2px 8px; border-radius: 99px; font-family: var(--mono); }
        .chip.green { background: rgba(16,185,129,0.15); color: var(--green); border: 1px solid rgba(16,185,129,0.3); }
        .chip.cyan  { background: rgba(6,182,212,0.15);  color: var(--cyan);  border: 1px solid rgba(6,182,212,0.3);  }
        .chip.purple{ background: rgba(139,92,246,0.15); color: var(--purple);border: 1px solid rgba(139,92,246,0.3); }
        .chip.orange{ background: rgba(245,158,11,0.15); color: var(--orange);border: 1px solid rgba(245,158,11,0.3); }

        .chat-input-row { padding: 1rem 1.5rem; border-top: 1px solid var(--border); display: flex; gap: 0.75rem; background: var(--surface); }
        .chat-input { flex: 1; background: var(--card); border: 1px solid var(--border); border-radius: 10px; color: var(--text); padding: 0.75rem 1rem; font-size: 0.95rem; font-family: inherit; outline: none; resize: none; transition: border 0.2s; }
        .chat-input:focus { border-color: var(--green); }
        .send-btn { background: linear-gradient(135deg, var(--green) 0%, var(--green2) 100%); color: #fff; border: none; border-radius: 10px; padding: 0.75rem 1.25rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; font-family: inherit; transition: opacity 0.2s, transform 0.1s; white-space: nowrap; }
        .send-btn:hover { opacity: 0.9; transform: translateY(-1px); }
        .send-btn:disabled { opacity: 0.4; cursor: not-allowed; transform: none; }

        /* ── API Tester Panel ── */
        .tester { flex: 1; overflow-y: auto; padding: 1.5rem; display: flex; flex-direction: column; gap: 1.25rem; }
        .tester-url { display: flex; gap: 0.5rem; align-items: center; background: var(--card); border: 1px solid var(--border); border-radius: 10px; padding: 0.5rem 1rem; }
        .tester-url .method-tag { font-family: var(--mono); font-size: 0.75rem; font-weight: 700; color: var(--cyan); }
        .tester-url input { flex: 1; background: transparent; border: none; color: var(--text); font-family: var(--mono); font-size: 0.85rem; outline: none; }
        .run-btn { background: var(--cyan); color: #000; border: none; padding: 0.4rem 1rem; border-radius: 6px; font-weight: 700; cursor: pointer; font-size: 0.85rem; font-family: inherit; }
        .run-btn:hover { opacity: 0.85; }
        .params-grid { display: grid; grid-template-columns: 1fr 2fr; gap: 0.5rem; }
        .params-grid label { font-size: 0.8rem; color: var(--muted); display: flex; align-items: center; }
        .params-grid input { background: var(--card); border: 1px solid var(--border); border-radius: 6px; color: var(--text); padding: 0.4rem 0.7rem; font-size: 0.85rem; outline: none; }
        .params-grid input:focus { border-color: var(--cyan); }
        .response-box { background: var(--bg); border: 1px solid var(--border); border-radius: 10px; padding: 1rem; font-family: var(--mono); font-size: 0.8rem; color: #a8d8a8; white-space: pre-wrap; overflow-x: auto; max-height: 400px; overflow-y: auto; min-height: 100px; }
        .response-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.5rem; }
        .response-header span { font-size: 0.8rem; color: var(--muted); }
        .status-ok { color: var(--green); font-weight: 700; }
        .status-err { color: var(--red); font-weight: 700; }

        /* ── Typing indicator ── */
        .typing { display: flex; gap: 4px; align-items: center; padding: 0.85rem 1.1rem; background: var(--card); border: 1px solid var(--border); border-radius: 14px; border-bottom-left-radius: 4px; width: fit-content; }
        .typing span { width: 7px; height: 7px; border-radius: 50%; background: var(--muted); animation: bounce 1.2s infinite; }
        .typing span:nth-child(2) { animation-delay: 0.2s; }
        .typing span:nth-child(3) { animation-delay: 0.4s; }
        @keyframes bounce { 0%,60%,100% { transform: translateY(0); } 30% { transform: translateY(-6px); } }

        /* ── Misc ── */
        hr { border: none; border-top: 1px solid var(--border); }
        .clarify-btns { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 0.75rem; }
        .clarify-btn { background: var(--card); border: 1px solid var(--border); border-radius: 8px; color: var(--text); padding: 0.5rem 0.75rem; cursor: pointer; font-family: inherit; font-size: 0.85rem; text-align: left; transition: border 0.15s; }
        .clarify-btn:hover { border-color: var(--green); color: var(--green); }

        @media (max-width: 768px) {
            .app { grid-template-columns: 1fr; grid-template-rows: 60px auto 1fr; }
            .sidebar { max-height: 280px; }
        }
    </style>
</head>
<body>
<div class="app">

    <!-- Topbar -->
    <div class="topbar">
        <span class="logo">⚡ AfriLearn</span>
        <span class="badge-live">LIVE</span>
        <span style="color:var(--muted);font-size:0.85rem;">Intelligence Playground</span>
        <div class="topbar-links">
            <a href="/portal" target="_blank">🔑 Portal</a>
            <a href="/docs" target="_blank">📖 Docs</a>
            <a href="/health" target="_blank">💚 Health</a>
            <a href="https://github.com/HarperKollins/AfriLearn" target="_blank">🐙 GitHub</a>
        </div>
    </div>

    <!-- Sidebar -->
    <div class="sidebar">
        <div class="section-label">🔧 Configuration</div>
        <div class="config-group">
            <label>AfriLearn API Key</label>
            <input type="text" id="afriKey" placeholder="afr_live_..." value="afr_live_demo_9f8e2b7a">
        </div>
        <div class="config-group">
            <label>AI Provider</label>
            <select id="aiProvider" onchange="updateProviderHint()">
                <option value="groq" selected>⚡ Groq (Llama 3.3 70B - Free & Ultra Fast)</option>
                <option value="openai">OpenAI (GPT-4o / GPT-4o-mini)</option>
                <option value="gemini">Google Gemini</option>
                <option value="anthropic">Anthropic Claude</option>
            </select>
            <label>AI API Key</label>
            <input type="password" id="aiKey" placeholder="gsk_... or sk-... or AIza...">
            <label style="margin-top:0.25rem;font-size:0.72rem;color:var(--muted)" id="providerHint">Get free key: console.groq.com</label>
        </div>

        <hr>
        <div class="section-label">🧠 Panels</div>
        <button class="nav-btn active" onclick="showPanel('chat')" id="btn-chat">
            💬 AI Tutor Chat
        </button>
        <button class="nav-btn" onclick="showPanel('query')" id="btn-query">
            <span class="method post">POST</span> /api/v1/query
        </button>
        <button class="nav-btn" onclick="showPanel('match')" id="btn-match">
            <span class="method get">GET</span> /curriculum/match/:topic
        </button>
        <button class="nav-btn" onclick="showPanel('pathway')" id="btn-pathway">
            <span class="method get">GET</span> /curriculum/pathway
        </button>
        <button class="nav-btn" onclick="showPanel('curriculum')" id="btn-curriculum">
            <span class="method get">GET</span> /curriculum/:board/:subject
        </button>
        <button class="nav-btn" onclick="showPanel('llmprompt')" id="btn-llmprompt">
            <span class="method get">GET</span> /llm-prompt
        </button>

        <hr>
        <div class="section-label">⚡ Quick Tests</div>
        <button class="nav-btn" onclick="quickChat('What topics do I need for JAMB Physics?')">JAMB Physics topics</button>
        <button class="nav-btn" onclick="quickChat('Give me the learning pathway for mathematics from BECE to JAMB')">Math pathway BECE→JAMB</button>
        <button class="nav-btn" onclick="quickChat('Explain quadratic equations using WAEC curriculum')">Quadratic equations</button>
        <button class="nav-btn" onclick="quickChat('What is the NUC Computer Science degree curriculum?')">NUC CS degree</button>
    </div>

    <!-- ── CHAT PANEL ── -->
    <div class="panel active" id="panel-chat">
        <div class="chat-messages" id="chatMessages">
            <div class="msg system">
                <div class="bubble">
                    👋 Ask me anything about Nigerian curriculum — I'll use AfriLearn to fetch the curriculum context, then send it to your chosen AI to give you a real answer. Set your AI key in the sidebar first.
                </div>
            </div>
        </div>
        <div class="chat-input-row">
            <textarea class="chat-input" id="chatInput" rows="2" placeholder="Ask anything... e.g. 'Explain WAEC Physics Waves to me' or 'What topics in JAMB Math are hardest?'" onkeydown="handleChatKey(event)"></textarea>
            <button class="send-btn" id="sendBtn" onclick="sendChat()">Ask AI ✨</button>
        </div>
    </div>

    <!-- ── QUERY BRAIN PANEL ── -->
    <div class="panel" id="panel-query">
        <div class="tester">
            <h3 style="color:var(--purple)">🧠 Curriculum Query Brain</h3>
            <p style="color:var(--muted);font-size:0.9rem">Natural language query with intent parsing, clarification loops, and smart caching.</p>
            <div class="tester-url">
                <span class="method-tag" style="color:var(--purple)">POST</span>
                <input type="text" id="queryUrl" value="/api/v1/query" readonly>
                <button class="run-btn" style="background:var(--purple)" onclick="runQuery()">Run</button>
            </div>
            <div class="config-group">
                <label>Question</label>
                <input type="text" id="queryQuestion" value="What topics do I need for WAEC Physics?" placeholder="Ask a curriculum question...">
                <label>Session ID (for follow-up clarification)</label>
                <input type="text" id="querySession" placeholder="Leave empty for first question">
            </div>
            <div class="response-header"><span>Response</span><span id="queryStatus" class="status-ok"></span></div>
            <div class="response-box" id="queryResponse">// Hit "Run" to see the response...</div>
        </div>
    </div>

    <!-- ── MATCH PANEL ── -->
    <div class="panel" id="panel-match">
        <div class="tester">
            <h3 style="color:var(--cyan)">🔍 Cross-Board Curriculum Match</h3>
            <p style="color:var(--muted);font-size:0.9rem">Query all 22 boards simultaneously and see where a topic appears across BECE → WAEC → JAMB → NUC.</p>
            <div class="tester-url">
                <span class="method-tag">GET</span>
                <input type="text" id="matchUrl" value="/api/v1/curriculum/match/algebra">
                <button class="run-btn" onclick="runMatch()">Run</button>
            </div>
            <div class="params-grid">
                <label>Topic</label>
                <input type="text" id="matchTopic" value="algebra" oninput="document.getElementById('matchUrl').value='/api/v1/curriculum/match/'+this.value">
            </div>
            <div class="response-header"><span>Response</span><span id="matchStatus"></span></div>
            <div class="response-box" id="matchResponse">// Hit "Run" to see cross-board coverage...</div>
        </div>
    </div>

    <!-- ── PATHWAY PANEL ── -->
    <div class="panel" id="panel-pathway">
        <div class="tester">
            <h3 style="color:var(--orange)">🗺️ Learning Pathway Engine</h3>
            <p style="color:var(--muted);font-size:0.9rem">Returns an ordered step-by-step learning journey between curriculum boards for a subject.</p>
            <div class="tester-url">
                <span class="method-tag">GET</span>
                <input type="text" id="pathwayUrl" value="/api/v1/curriculum/pathway?subject=mathematics&from=bece&to=jamb">
                <button class="run-btn" onclick="runPathway()">Run</button>
            </div>
            <div class="params-grid">
                <label>Subject</label>
                <input type="text" id="pathwaySubject" value="mathematics" oninput="updatePathwayUrl()">
                <label>From Board</label>
                <input type="text" id="pathwayFrom" value="bece" oninput="updatePathwayUrl()">
                <label>To Board</label>
                <input type="text" id="pathwayTo" value="jamb" oninput="updatePathwayUrl()">
            </div>
            <div class="response-header"><span>Response</span><span id="pathwayStatus"></span></div>
            <div class="response-box" id="pathwayResponse">// Hit "Run" to generate learning pathway...</div>
        </div>
    </div>

    <!-- ── CURRICULUM PANEL ── -->
    <div class="panel" id="panel-curriculum">
        <div class="tester">
            <h3 style="color:var(--green)">📚 Full Curriculum Tree</h3>
            <p style="color:var(--muted);font-size:0.9rem">Fetch the complete curriculum with topics, subtopics, and learning objectives.</p>
            <div class="tester-url">
                <span class="method-tag">GET</span>
                <input type="text" id="currUrl" value="/api/v1/curriculum/waec/physics">
                <button class="run-btn" onclick="runCurriculum()">Run</button>
            </div>
            <div class="params-grid">
                <label>Board</label>
                <input type="text" id="currBoard" value="waec" oninput="updateCurrUrl()">
                <label>Subject</label>
                <input type="text" id="currSubject" value="physics" oninput="updateCurrUrl()">
            </div>
            <div class="response-header"><span>Response</span><span id="currStatus"></span></div>
            <div class="response-box" id="currResponse">// Hit "Run" to fetch curriculum...</div>
        </div>
    </div>

    <!-- ── LLM PROMPT PANEL ── -->
    <div class="panel" id="panel-llmprompt">
        <div class="tester">
            <h3 style="color:var(--purple)">🤖 AI Tutor LLM Prompt Generator</h3>
            <p style="color:var(--muted);font-size:0.9rem">Generates a pre-formatted system prompt, context window, and module blocks for any AI model.</p>
            <div class="tester-url">
                <span class="method-tag">GET</span>
                <input type="text" id="llmUrl" value="/api/v1/curriculum/waec/physics/llm-prompt">
                <button class="run-btn" onclick="runLLM()">Run</button>
            </div>
            <div class="params-grid">
                <label>Board</label>
                <input type="text" id="llmBoard" value="waec" oninput="updateLLMUrl()">
                <label>Subject</label>
                <input type="text" id="llmSubject" value="physics" oninput="updateLLMUrl()">
            </div>
            <div class="response-header"><span>Response</span><span id="llmStatus"></span></div>
            <div class="response-box" id="llmResponse">// Hit "Run" to generate AI Tutor system prompt...</div>
        </div>
    </div>

</div>

<script>
const BASE = window.location.origin;
let sessionId = null;

// ── Utility ──────────────────────────────────────────────────────────────────
function getAfriKey() { return document.getElementById('afriKey').value.trim(); }
function getAIKey()   { return document.getElementById('aiKey').value.trim(); }
function getProvider(){ return document.getElementById('aiProvider').value; }

function afriHeaders() {
    return { 'X-API-Key': getAfriKey(), 'Content-Type': 'application/json' };
}

function updateProviderHint() {
    const hints = {
        groq: 'Get free key: console.groq.com',
        openai: 'Get key: platform.openai.com',
        gemini: 'Get key: aistudio.google.com',
        anthropic: 'Get key: console.anthropic.com'
    };
    document.getElementById('providerHint').innerText = hints[getProvider()] || '';
}

function showPanel(name) {
    document.querySelectorAll('.panel').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.nav-btn').forEach(b => b.classList.remove('active'));
    document.getElementById('panel-' + name).classList.add('active');
    document.getElementById('btn-' + name).classList.add('active');
}

function fmtJSON(obj) { return JSON.stringify(obj, null, 2); }

function setResponse(id, statusId, statusCode, data) {
    document.getElementById(id).innerText = fmtJSON(data);
    const el = document.getElementById(statusId);
    if (el) {
        el.innerText = statusCode + (statusCode < 300 ? ' OK' : ' ERROR');
        el.className = statusCode < 300 ? 'status-ok' : 'status-err';
    }
}

// ── URL updaters ──────────────────────────────────────────────────────────────
function updatePathwayUrl() {
    const s = document.getElementById('pathwaySubject').value;
    const f = document.getElementById('pathwayFrom').value;
    const t = document.getElementById('pathwayTo').value;
    document.getElementById('pathwayUrl').value = '/api/v1/curriculum/pathway?subject='+s+'&from='+f+'&to='+t;
}
function updateCurrUrl() {
    document.getElementById('currUrl').value = '/api/v1/curriculum/'+document.getElementById('currBoard').value+'/'+document.getElementById('currSubject').value;
}
function updateLLMUrl() {
    document.getElementById('llmUrl').value = '/api/v1/curriculum/'+document.getElementById('llmBoard').value+'/'+document.getElementById('llmSubject').value+'/llm-prompt';
}

// ── API Testers ───────────────────────────────────────────────────────────────
async function runQuery() {
    const q = document.getElementById('queryQuestion').value;
    const sid = document.getElementById('querySession').value;
    const body = sid ? { question: q, session_id: sid } : { question: q };
    const res = await fetch(BASE + '/api/v1/query', { method: 'POST', headers: afriHeaders(), body: JSON.stringify(body) });
    const data = await res.json();
    setResponse('queryResponse', 'queryStatus', res.status, data);
    if (data.data && data.data.session_id) {
        document.getElementById('querySession').value = data.data.session_id;
    }
}

async function runMatch() {
    const topic = document.getElementById('matchTopic').value;
    const res = await fetch(BASE + '/api/v1/curriculum/match/' + topic, { headers: afriHeaders() });
    const data = await res.json();
    setResponse('matchResponse', 'matchStatus', res.status, data);
}

async function runPathway() {
    const url = document.getElementById('pathwayUrl').value;
    const res = await fetch(BASE + url, { headers: afriHeaders() });
    const data = await res.json();
    setResponse('pathwayResponse', 'pathwayStatus', res.status, data);
}

async function runCurriculum() {
    const url = document.getElementById('currUrl').value;
    const res = await fetch(BASE + url, { headers: afriHeaders() });
    const data = await res.json();
    setResponse('currResponse', 'currStatus', res.status, data);
}

async function runLLM() {
    const url = document.getElementById('llmUrl').value;
    const res = await fetch(BASE + url, { headers: afriHeaders() });
    const data = await res.json();
    setResponse('llmResponse', 'llmStatus', res.status, data);
}

// ── Chat ──────────────────────────────────────────────────────────────────────
function handleChatKey(e) {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendChat(); }
}

function quickChat(q) {
    showPanel('chat');
    document.getElementById('chatInput').value = q;
    sendChat();
}

function addMsg(role, html, chips) {
    const box = document.getElementById('chatMessages');
    const div = document.createElement('div');
    div.className = 'msg ' + role;
    let inner = '<div class="bubble">' + html + '</div>';
    if (chips && chips.length) {
        inner += '<div class="step-chips">' + chips.map(c => '<span class="chip '+c.color+'">'+c.label+'</span>').join('') + '</div>';
    }
    div.innerHTML = inner;
    box.appendChild(div);
    box.scrollTop = box.scrollHeight;
    return div;
}

function addTyping() {
    const box = document.getElementById('chatMessages');
    const div = document.createElement('div');
    div.className = 'msg assistant';
    div.id = 'typingIndicator';
    div.innerHTML = '<div class="typing"><span></span><span></span><span></span></div>';
    box.appendChild(div);
    box.scrollTop = box.scrollHeight;
}

function removeTyping() {
    const el = document.getElementById('typingIndicator');
    if (el) el.remove();
}

// Topic keyword → subject slug mapping (client-side smart inference)
const topicSubjectMap = {
    // Biology
    photosynthesis: 'biology', osmosis: 'biology', mitosis: 'biology', meiosis: 'biology',
    respiration: 'biology', genetics: 'biology', ecosystem: 'biology', cell: 'biology',
    dna: 'biology', rna: 'biology', enzyme: 'biology', diffusion: 'biology',
    chlorophyll: 'biology', nucleus: 'biology', chromosome: 'biology', virus: 'biology',
    bacteria: 'biology', fungi: 'biology', evolution: 'biology', excretion: 'biology',
    // Mathematics
    quadratic: 'mathematics', algebra: 'mathematics', trigonometry: 'mathematics',
    calculus: 'mathematics', statistics: 'mathematics', probability: 'mathematics',
    geometry: 'mathematics', differentiation: 'mathematics', integration: 'mathematics',
    logarithm: 'mathematics', matrix: 'mathematics', vector: 'mathematics',
    fraction: 'mathematics', equation: 'mathematics', arithmetic: 'mathematics',
    // Physics
    newton: 'physics', velocity: 'physics', acceleration: 'physics', momentum: 'physics',
    gravity: 'physics', electricity: 'physics', magnetism: 'physics', optics: 'physics',
    wave: 'physics', thermodynamics: 'physics', pressure: 'physics', energy: 'physics',
    force: 'physics', motion: 'physics', nuclear: 'physics', electromagnetic: 'physics',
    // Chemistry
    acid: 'chemistry', base: 'chemistry', periodic: 'chemistry', element: 'chemistry',
    compound: 'chemistry', reaction: 'chemistry', mole: 'chemistry', bond: 'chemistry',
    titration: 'chemistry', electrolysis: 'chemistry', hydrocarbon: 'chemistry',
    oxidation: 'chemistry', reduction: 'chemistry', salt: 'chemistry', alkane: 'chemistry',
    // Economics
    supply: 'economics', demand: 'economics', inflation: 'economics', gdp: 'economics',
    market: 'economics', elasticity: 'economics', unemployment: 'economics',
    monopoly: 'economics', fiscal: 'economics', monetary: 'economics',
    // Government
    democracy: 'government', constitution: 'government', federalism: 'government',
    legislature: 'government', judiciary: 'government', executive: 'government',
    election: 'government', sovereignty: 'government', citizenship: 'government',
    // Literature / English
    poem: 'literature-in-english', novel: 'literature-in-english', prose: 'literature-in-english',
    shakespeare: 'literature-in-english', metaphor: 'literature-in-english', sonnet: 'literature-in-english',
    grammar: 'english-studies', comprehension: 'english-studies', essay: 'english-studies',
    // Computer Science
    algorithm: 'computer-science', programming: 'computer-science', database: 'computer-science',
    network: 'computer-science', binary: 'computer-science', software: 'computer-science',
    hardware: 'computer-science', operating: 'computer-science', data: 'computer-science',
    // Medicine
    anatomy: 'medicine-surgery', physiology: 'medicine-surgery', pharmacology: 'medicine-surgery',
    diagnosis: 'medicine-surgery', pathology: 'medicine-surgery', surgery: 'medicine-surgery',
};

function guessSubjectFromText(text) {
    const lower = text.toLowerCase();
    for (const [keyword, subject] of Object.entries(topicSubjectMap)) {
        if (lower.includes(keyword)) return subject;
    }
    return null;
}

async function sendChat() {
    const input = document.getElementById('chatInput');
    const question = input.value.trim();
    if (!question) return;
    input.value = '';
    document.getElementById('sendBtn').disabled = true;

    addMsg('user', question);

    // Client-side subject inference — tries to resolve subject before hitting AfriLearn
    const guessedSubject = guessSubjectFromText(question);

    // Step 1: Call AfriLearn Query Brain
    addTyping();
    let afriData = null;
    try {
        const body = sessionId ? { question, session_id: sessionId } : { question };
        const r = await fetch(BASE + '/api/v1/query', {
            method: 'POST', headers: afriHeaders(), body: JSON.stringify(body)
        });
        afriData = await r.json();
    } catch(e) {
        removeTyping();
        addMsg('system', '❌ Could not reach AfriLearn API. Check your AfriLearn key in the sidebar.');
        document.getElementById('sendBtn').disabled = false;
        return;
    }
    removeTyping();

    // Handle clarification needed — but inject guessed subject if we have one
    if (afriData.data && afriData.data.needs_clarification) {
        const remaining = afriData.data.clarification_required || [];

        // If we could guess the subject from topic keywords, auto-retry with it
        if (guessedSubject && remaining.some(q => q.toLowerCase().includes('subject'))) {
            const sessionIdLocal = afriData.data.session_id;
            addTyping();
            try {
                const body2 = { question: guessedSubject, session_id: sessionIdLocal };
                const r2 = await fetch(BASE + '/api/v1/query', {
                    method: 'POST', headers: afriHeaders(), body: JSON.stringify(body2)
                });
                afriData = await r2.json();
            } catch(e) {}
            removeTyping();
            // If still needs clarification, fall through to show it
            if (!afriData.data || afriData.data.needs_clarification) {
                sessionId = afriData.data && afriData.data.session_id ? afriData.data.session_id : sessionIdLocal;
                const qs = afriData.data && afriData.data.clarification_required ? afriData.data.clarification_required : remaining;
                let html = '<strong>🤔 Just one more thing:</strong><br><br>';
                html += qs.map(q => '• ' + q).join('<br>');
                addMsg('assistant', html, [{ label: '💬 clarification', color: 'orange' }]);
                document.getElementById('sendBtn').disabled = false;
                return;
            }
        } else {
            sessionId = afriData.data.session_id;
            const qs = remaining;
            let html = '<strong>🤔 Quick question:</strong><br><br>';
            html += qs.map(q => '• ' + q).join('<br>');
            html += '<br><br><em style="color:var(--muted);font-size:0.82rem">Tip: You can ask like "WAEC Biology" or "JAMB Physics" to skip this.</em>';
            addMsg('assistant', html, [{ label: '💬 clarification', color: 'orange' }]);
            document.getElementById('sendBtn').disabled = false;
            return;
        }
    }

    sessionId = null;

    // Check if no AI key
    const aiKey = getAIKey();
    if (!aiKey) {
        // Just show AfriLearn result without AI
        const result = afriData.data && afriData.data.result ? afriData.data.result : afriData.data;
        const topics = result && result.topics ? result.topics.slice(0, 8).map(t => typeof t === 'string' ? t : t.name).join(', ') : '';
        let html = '<strong>📚 AfriLearn Result</strong>';
        if (result && result.total_topics) html += ' <span style="color:var(--muted);font-size:0.85rem">('+result.total_topics+' topics)</span>';
        html += '<br><br>';
        if (topics) html += '<strong>Topics:</strong> ' + topics + (result.total_topics > 8 ? ', ...' : '');
        else html += '<pre style="white-space:pre-wrap;font-size:0.8rem">' + JSON.stringify(result, null, 2).slice(0, 600) + '...</pre>';
        html += '<br><br><em style="color:var(--muted);font-size:0.82rem">💡 Add an AI key in the sidebar to get a real AI answer using this curriculum!</em>';
        addMsg('assistant', html, [
            { label: '✅ AfriLearn', color: 'green' },
            { label: afriData.data.intent ? afriData.data.intent.board + '/' + afriData.data.intent.subject : '', color: 'cyan' },
            { label: '🔑 no AI key', color: 'orange' }
        ]);
        document.getElementById('sendBtn').disabled = false;
        return;
    }

    // Step 2: Get LLM system prompt from AfriLearn
    const intent = afriData.data && afriData.data.intent ? afriData.data.intent : {};
    let systemPrompt = 'You are an expert African curriculum AI Tutor. Answer the student\'s question clearly and helpfully.';
    if (intent.board && intent.subject) {
        try {
            const llmR = await fetch(BASE + '/api/v1/curriculum/' + intent.board + '/' + intent.subject + '/llm-prompt', {
                headers: afriHeaders()
            });
            const llmData = await llmR.json();
            if (llmData.data && llmData.data.system_prompt) {
                systemPrompt = llmData.data.system_prompt;
            }
        } catch(e) { /* use default */ }
    }

    // Step 3: Call AI
    addTyping();
    let aiAnswer = '';
    const provider = getProvider();
    try {
        if (provider === 'groq') {
            let r = await fetch('https://api.groq.com/openai/v1/chat/completions', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + aiKey },
                body: JSON.stringify({
                    model: 'llama-3.3-70b-versatile',
                    messages: [
                        { role: 'system', content: systemPrompt },
                        { role: 'user', content: question }
                    ],
                    temperature: 0.7,
                    max_tokens: 1024
                })
            });
            let d = await r.json();
            // Fallback model if 70b is busy
            if (!d.choices && d.error) {
                r = await fetch('https://api.groq.com/openai/v1/chat/completions', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + aiKey },
                    body: JSON.stringify({
                        model: 'llama3-8b-8192',
                        messages: [
                            { role: 'system', content: systemPrompt },
                            { role: 'user', content: question }
                        ],
                        temperature: 0.7,
                        max_tokens: 1024
                    })
                });
                d = await r.json();
            }
            if (d.choices && d.choices[0].message) aiAnswer = d.choices[0].message.content;
            else aiAnswer = '⚠️ Groq Error: ' + JSON.stringify(d.error || d);
        } else if (provider === 'openai') {
            const r = await fetch('https://api.openai.com/v1/chat/completions', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + aiKey },
                body: JSON.stringify({
                    model: 'gpt-4o-mini',
                    messages: [
                        { role: 'system', content: systemPrompt },
                        { role: 'user', content: question }
                    ],
                    temperature: 0.7,
                    max_tokens: 800
                })
            });
            const d = await r.json();
            if (d.choices) aiAnswer = d.choices[0].message.content;
            else aiAnswer = '⚠️ OpenAI Error: ' + JSON.stringify(d.error || d);
        } else if (provider === 'gemini') {
            // Try gemini-2.0-flash first (latest free model), fall back to gemini-1.5-flash-latest
            let geminiModel = 'gemini-2.0-flash';
            let r = await fetch('https://generativelanguage.googleapis.com/v1beta/models/' + geminiModel + ':generateContent?key=' + aiKey, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    contents: [{ parts: [{ text: systemPrompt + '\n\nStudent question: ' + question }] }],
                    generationConfig: { temperature: 0.7, maxOutputTokens: 1024 }
                })
            });
            let d = await r.json();
            // Fallback to gemini-1.5-flash-latest if 2.0-flash not available
            if (!d.candidates && d.error && d.error.code === 404) {
                geminiModel = 'gemini-1.5-flash-latest';
                r = await fetch('https://generativelanguage.googleapis.com/v1beta/models/' + geminiModel + ':generateContent?key=' + aiKey, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        contents: [{ parts: [{ text: systemPrompt + '\n\nStudent question: ' + question }] }],
                        generationConfig: { temperature: 0.7, maxOutputTokens: 1024 }
                    })
                });
                d = await r.json();
            }
            if (d.candidates && d.candidates[0].content) aiAnswer = d.candidates[0].content.parts[0].text;
            else aiAnswer = '⚠️ Gemini Error: ' + JSON.stringify(d.error || d);
        } else if (provider === 'anthropic') {
            const r = await fetch('https://api.anthropic.com/v1/messages', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'x-api-key': aiKey, 'anthropic-version': '2023-06-01', 'anthropic-dangerous-direct-browser-access': 'true' },
                body: JSON.stringify({
                    model: 'claude-3-haiku-20240307',
                    max_tokens: 800,
                    system: systemPrompt,
                    messages: [{ role: 'user', content: question }]
                })
            });
            const d = await r.json();
            if (d.content) aiAnswer = d.content[0].text;
            else aiAnswer = '⚠️ Claude Error: ' + JSON.stringify(d.error || d);
        }
    } catch(e) {
        aiAnswer = '⚠️ AI call failed: ' + e.message + '. Make sure your AI key is valid and CORS allows this request.';
    }

    removeTyping();

    // Format answer as HTML (preserve line breaks)
    const formatted = aiAnswer.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>').replace(/\n/g, '<br>');
    const chips = [
        { label: '✅ AfriLearn', color: 'green' },
        { label: intent.board && intent.subject ? intent.board + '/' + intent.subject : 'query', color: 'cyan' },
        { label: provider, color: 'purple' }
    ];
    if (afriData.data && afriData.data.cache_hit) chips.push({ label: '⚡ cached', color: 'orange' });
    addMsg('assistant', formatted, chips);
    document.getElementById('sendBtn').disabled = false;
}
</script>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
