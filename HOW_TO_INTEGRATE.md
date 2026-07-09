# AfriLearn × HK AI — Integration Guide

> **AfriLearn Base URL**: `https://afrilearn-yka7.onrender.com`
> **API Key**: `afr_live_b99390afb1fcb6005379f803f7103782`
> **Demo Key** (rate-limited, public): `afr_live_demo_9f8e2b7a`

AfriLearn is a Nigerian curriculum intelligence API. It provides official WAEC, JAMB, BECE, NUC, and NBTE curriculum data, LLM system prompts, Bloom's taxonomy breakdowns, and subject-specific pedagogy rules — all pre-packaged as a single HTTP call your AI model can consume directly.

---

## How It Works

```
HK AI (Node.js / Fastify)
        │
        │  POST https://afrilearn-yka7.onrender.com/api/v1/bridge/enrich
        │  Header: X-API-Key: afr_live_b99390afb1fcb6005379f803f7103782
        │  Body:   { messages: [...chatHistory], workspace_title: "WAEC Physics" }
        ▼
AfriLearn API returns:
  {
    board: "waec",
    subject: "physics",
    detection_score: "high",
    injection_prompt: "=== AFRILEARN CURRICULUM CONTEXT === ..."  ← paste into Groq/Gemini
    topics: [...],          ← official topic tree for lesson structuring
    blooms_taxonomy: {...}, ← cognitive level breakdown
    search_hits: [...]      ← FTS results from the conversation
  }
        │
        ▼
HK AI prepends injection_prompt to its Groq/Gemini system message
→ Every AI response is now grounded in official Nigerian curriculum
```

---

## Step 1 — Add Environment Variables

In **`server/.env`**, add these 2 lines:

```ini
# AfriLearn Curriculum Intelligence API
AFRILEARN_API_URL=https://afrilearn-yka7.onrender.com
AFRILEARN_API_KEY=afr_live_b99390afb1fcb6005379f803f7103782
```

---

## Step 2 — Create `server/utils/afrilearn.js`

Create this **new file** in HK AI. Do not modify any existing file yet.

```javascript
// server/utils/afrilearn.js
// AfriLearn Curriculum Intelligence Client for HK AI

import axios from 'axios';

const AFRILEARN_URL = process.env.AFRILEARN_API_URL || 'https://afrilearn-yka7.onrender.com';
const AFRILEARN_KEY = process.env.AFRILEARN_API_KEY || 'afr_live_demo_9f8e2b7a';

const client = axios.create({
  baseURL: AFRILEARN_URL,
  headers: {
    'X-API-Key': AFRILEARN_KEY,
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10s — never block course generation if AfriLearn is slow
});

/**
 * enrichWithCurriculum
 *
 * Call BEFORE passing messages to Groq/Gemini.
 * AfriLearn auto-detects the board + subject from the conversation,
 * then returns the full curriculum context as a ready-to-inject prompt string.
 *
 * @param {Object} options
 * @param {Array}  options.messages        - Chat history: [{ role: 'user'|'assistant', content: '...' }]
 * @param {string} options.workspaceTitle  - Workspace title for extra detection signal
 * @param {string} [options.board]         - Optional: explicit board slug e.g. 'waec', 'jamb', 'unilag'
 * @param {string} [options.subject]       - Optional: explicit subject slug e.g. 'physics', 'law'
 *
 * @returns {Object|null} AfriLearn context:
 *   {
 *     injection_prompt: string,   ← PASTE THIS into Groq/Gemini system message
 *     board: string,              ← e.g. 'waec'
 *     subject: string,            ← e.g. 'physics'
 *     board_full_name: string,    ← e.g. 'West African Examinations Council'
 *     subject_name: string,       ← e.g. 'Physics'
 *     detection_score: string,    ← 'explicit' | 'high' | 'medium' | 'low' | 'none'
 *     topics: Array,              ← official topic tree
 *     blooms_taxonomy: Object,    ← cognitive breakdown
 *     search_hits: Array          ← FTS hits from conversation
 *   }
 *
 * Returns null if AfriLearn is unreachable — HK AI continues normally (fail-safe).
 */
export async function enrichWithCurriculum({
  messages = [],
  workspaceTitle = '',
  board = null,
  subject = null,
}) {
  try {
    const body = { messages, workspace_title: workspaceTitle };
    if (board && subject) {
      body.board = board;
      body.subject = subject;
    }
    const res = await client.post('/api/v1/bridge/enrich', body);
    if (res.data?.success && res.data?.data) {
      return res.data.data;
    }
    return null;
  } catch (err) {
    console.warn('[AfriLearn] enrichWithCurriculum failed:', err.message);
    return null;
  }
}

/**
 * getAvailableBoards
 *
 * Returns all AfriLearn boards + their available subjects.
 * Use this to build the board/subject picker in the workspace creation UI.
 *
 * @returns {Array}
 *   [
 *     { slug: 'waec', name: 'WAEC', full_name: '...', subjects: ['physics', 'chemistry', ...] },
 *     { slug: 'jamb', ... },
 *     ... 12 total boards
 *   ]
 */
export async function getAvailableBoards() {
  try {
    const res = await client.get('/api/v1/bridge/boards');
    if (res.data?.success && res.data?.data) {
      return res.data.data;
    }
    return [];
  } catch (err) {
    console.warn('[AfriLearn] getAvailableBoards failed:', err.message);
    return [];
  }
}
```

---

## Step 3 — Inject into Course Generation (`prompts.js`)

### 3a. Add the import at the very top of `server/utils/prompts.js`

```javascript
import { enrichWithCurriculum } from './afrilearn.js';
```

### 3b. Replace `courseSkeleton` with an async version

Find the existing `export const courseSkeleton = ...` (around line 192) and **replace it**:

```javascript
export const buildCourseSkeletonPrompt = async (
  memoryContext,
  chatContext,
  messages = [],
  workspaceTitle = '',
  board = null,
  subject = null
) => {
  // Call AfriLearn — fails silently if unreachable
  const afrilearn = await enrichWithCurriculum({ messages, workspaceTitle, board, subject });

  const curriculumContext = afrilearn?.injection_prompt
    ? `=== OFFICIAL NIGERIAN CURRICULUM CONTEXT (AfriLearn) ===
${afrilearn.injection_prompt}
=== END CURRICULUM CONTEXT ===

`
    : '';

  return `ROLE: Expert Curriculum Designer
TASK: Generate a course skeleton WITHOUT lesson content or quizzes.
IMPORTANT: Output strictly in valid JSON format. Do NOT include explanations or extra text.

${curriculumContext}USER MEMORY:
${memoryContext}

CHAT HISTORY:
${chatContext}

RETURN ONLY JSON:
{
  "course": {
    "title": "Course Title",
    "description": "Brief description",
    "difficulty": "Beginner | Intermediate | Advanced",
    "estimated_time": "e.g. '5 hours'"
  },
  "lessons": [
    {
      "title": "Lesson Title",
      "objectives": ["Objective 1", "Objective 2"],
      "videoQuery": "YouTube search query for this lesson"
    }
  ],
  "new_memory": {
    "key": "current_course_topic",
    "value": "The specific topic from the conversation"
  }
}
`;
};
```

---

## Step 4 — Update `inMemoryQueue.js`

**Find** where `courseSkeleton` is imported and called. It looks like:

```javascript
// EXISTING (find these lines):
import { courseSkeleton, lessonGenPrompt, ... } from './prompts.js';
// ...
const prompt = courseSkeleton(memoryContext, chatContext);
```

**Replace** with:

```javascript
// UPDATED import — swap courseSkeleton for buildCourseSkeletonPrompt:
import { buildCourseSkeletonPrompt, lessonGenPrompt, chatPrompt, quizGenPrompt, videoQueryPrompt } from './prompts.js';

// ...

// UPDATED call — now async, passes messages + workspace info:
const prompt = await buildCourseSkeletonPrompt(
  memoryContext,
  chatContext,
  messages,              // the chat history array you already have in context
  workspace.title,       // workspace title from your DB query
  workspace.board_slug   || null,  // if you store board on workspace row
  workspace.subject_slug || null   // if you store subject on workspace row
);
```

> ⚠️ Make sure the worker function that calls this is declared `async` so `await` works.

---

## Step 5 — Enrich Chat Messages (Optional)

**File**: `server/routes/messages.js`

For regular chat (not just course generation), also enrich the AI system prompt:

```javascript
// At top of messages.js:
import { enrichWithCurriculum } from '../utils/afrilearn.js';

// Inside POST /workspace/:id/message handler, BEFORE calling Groq/Gemini:
const afrilearn = await enrichWithCurriculum({
  messages: recentMessages,    // the last N messages you already load
  workspaceTitle: workspace.title,
});

const curriculumAddition = afrilearn?.injection_prompt
  ? '\n\n' + afrilearn.injection_prompt
  : '';

// Append to your existing system prompt string:
const finalSystemPrompt = yourExistingSystemPrompt + curriculumAddition;
```

---

## Test the Connection

Run this from your HK AI `server/` directory:

```javascript
// test-afrilearn.mjs
import axios from 'axios';

const res = await axios.post(
  'https://afrilearn-yka7.onrender.com/api/v1/bridge/enrich',
  {
    messages: [
      { role: 'user', content: 'I need help with WAEC physics, specifically Newton laws and motion' }
    ],
    workspace_title: 'WAEC Physics Revision'
  },
  {
    headers: {
      'X-API-Key': 'afr_live_b99390afb1fcb6005379f803f7103782',
      'Content-Type': 'application/json'
    }
  }
);

const d = res.data.data;
console.log('Board:',            d.board);
console.log('Subject:',          d.subject);
console.log('Detection score:',  d.detection_score);
console.log('Topics count:',     d.topics?.length);
console.log('\nInjection prompt preview:');
console.log(d.injection_prompt?.substring(0, 400) + '...');
```

Run it:
```bash
node test-afrilearn.mjs
```

**Expected output:**
```
Board:            waec
Subject:          physics
Detection score:  high
Topics count:     6

Injection prompt preview:
=== AFRILEARN CURRICULUM INTELLIGENCE CONTEXT ===
You are an expert AI Tutor specialized in the official
West African Examinations Council Physics curriculum for Nigerian students.

OFFICIAL CURRICULUM TOPICS:
1. Interaction of Matter, Space, and Time
2. Energy (Work, Energy, and Power, Thermal Energy...)
...
```

---

## All Available API Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/health` | GET | ❌ None | Health check |
| `/api/v1/bridge/enrich` | POST | ✅ | **Main integration point** |
| `/api/v1/bridge/boards` | GET | ✅ | All boards + subjects for UI picker |
| `/api/v1/curriculum/:board/:subject` | GET | ✅ | Full topic tree |
| `/api/v1/curriculum/:board/:subject/llm-prompt` | GET | ✅ | Raw LLM system prompt |
| `/api/v1/search?q=<query>` | GET | ✅ | FTS across all 56 curricula |
| `/api/v1/curriculum/pathway/:topic` | GET | ✅ | BECE→WAEC→JAMB→Uni pathway |
| `/api/v1/curriculum/match/:topic` | GET | ✅ | Same topic across all 22 boards |

**Auth header for all protected endpoints:**
```
X-API-Key: afr_live_b99390afb1fcb6005379f803f7103782
```

---

## Available Boards + Subjects

| Board slug | Full Name | Subjects |
|---|---|---|
| `waec` | West African Examinations Council | biology, chemistry, economics, government, literature, mathematics, physics |
| `jamb` | Joint Admissions and Matriculation Board | biology, chemistry, economics, government, mathematics, physics |
| `bece` | Basic Education Certificate Examination | basic-science, basic-technology, business-studies, english-language, mathematics, social-studies |
| `nuc` | National Universities Commission (CCMAS) | accounting, business-administration, computer-science, electrical-engineering, law, mass-communication, mechanical-engineering, medicine-and-surgery, nursing-science |
| `unilag` | University of Lagos | accounting, computer-science, law, mass-communication, medicine-and-surgery |
| `futo` | Federal University of Technology, Owerri | computer-science, electrical-engineering, mechanical-engineering, petroleum-engineering |
| `unn` | University of Nigeria, Nsukka | computer-science, law, mechanical-engineering, medicine-and-surgery |
| `unec` | University of Nigeria, Enugu Campus | accounting, computer-science, law |
| `ebsu` | Ebonyi State University | accounting, computer-science, law, nursing-science |
| `funai` | Alex Ekwueme Federal University | accounting, computer-science, electrical-engineering, law |
| `yabatech` | Yaba College of Technology | computer-engineering-tech |
| `imt` | Institute of Management and Technology | accounting, computer-engineering-tech, science-laboratory-tech |

---

## What HK AI Gets After Integration

| Before | After |
|---|---|
| Generic AI-fabricated course topics | Official WAEC/JAMB topic order |
| Global examples | Nigerian analogies, Lagos context |
| No exam awareness | WAEC mark scheme rules, JAMB tricks |
| No domain rules | IRAC (Law), SOAP (Medicine), GIVEN/REQUIRED (Engineering) |
| No misconception handling | Proactive misconception flags per subject |
| Random lesson structure | Lessons aligned to real syllabus |

---

*AfriLearn — The open curriculum intelligence layer for African EdTech.*
*Built by HarperKollins. MIT License. https://github.com/HarperKollins/AfriLearn*
