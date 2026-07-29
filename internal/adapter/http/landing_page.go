package http

import (
	"net/http"
	"strings"
)

const landingHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>U.S. DOL Prevailing Wage API Documentation</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&family=Fira+Code:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-primary: #0a0d14;
            --bg-card: rgba(22, 27, 38, 0.7);
            --bg-code: #0d1117;
            --border: rgba(255, 255, 255, 0.08);
            --accent-primary: #6366f1;
            --accent-secondary: #06b6d4;
            --text-main: #f3f4f6;
            --text-muted: #9ca3af;
            --badge-get: #10b981;
            --badge-post: #3b82f6;
        }

        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg-primary);
            color: var(--text-main);
            line-height: 1.6;
            background-image: 
                radial-gradient(circle at 15% 20%, rgba(99, 102, 241, 0.15) 0%, transparent 40%),
                radial-gradient(circle at 85% 80%, rgba(6, 182, 212, 0.15) 0%, transparent 40%);
            background-attachment: fixed;
            min-height: 100vh;
        }

        .header {
            border-bottom: 1px solid var(--border);
            backdrop-filter: blur(12px);
            position: sticky;
            top: 0;
            z-index: 100;
            background: rgba(10, 13, 20, 0.8);
            padding: 1.25rem 2rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .brand {
            display: flex;
            align-items: center;
            gap: 0.75rem;
        }
        .brand-icon {
            background: linear-gradient(135deg, var(--accent-primary), var(--accent-secondary));
            width: 36px;
            height: 36px;
            border-radius: 8px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 800;
            color: white;
            font-size: 1.1rem;
        }
        .brand-title { font-weight: 700; font-size: 1.25rem; letter-spacing: -0.02em; }
        .status-pill {
            background: rgba(16, 185, 129, 0.15);
            color: #34d399;
            border: 1px solid rgba(16, 185, 129, 0.3);
            padding: 0.35rem 0.85rem;
            border-radius: 999px;
            font-size: 0.85rem;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        .status-dot {
            width: 8px;
            height: 8px;
            background: #34d399;
            border-radius: 50%;
            box-shadow: 0 0 10px #34d399;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 3rem 1.5rem;
        }

        .hero {
            text-align: center;
            margin-bottom: 4rem;
        }
        .hero h1 {
            font-size: 3rem;
            font-weight: 800;
            letter-spacing: -0.03em;
            margin-bottom: 1rem;
            background: linear-gradient(to right, #ffffff, #94a3b8);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        .hero p {
            font-size: 1.15rem;
            color: var(--text-muted);
            max-width: 750px;
            margin: 0 auto;
        }

        .grid {
            display: flex;
            flex-direction: column;
            gap: 2rem;
        }

        .card {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 16px;
            backdrop-filter: blur(16px);
            padding: 1.75rem;
            transition: border-color 0.3s ease, transform 0.2s ease;
        }
        .card:hover {
            border-color: rgba(99, 102, 241, 0.4);
        }

        .card-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 1rem;
            flex-wrap: wrap;
            gap: 1rem;
        }
        .endpoint-title {
            font-size: 1.2rem;
            font-weight: 700;
            display: flex;
            align-items: center;
            gap: 0.75rem;
        }
        .badge {
            padding: 0.25rem 0.6rem;
            border-radius: 6px;
            font-size: 0.75rem;
            font-weight: 700;
            text-transform: uppercase;
        }
        .badge-get { background: rgba(16, 185, 129, 0.2); color: var(--badge-get); border: 1px solid rgba(16, 185, 129, 0.4); }
        .badge-post { background: rgba(59, 130, 246, 0.2); color: var(--badge-post); border: 1px solid rgba(59, 130, 246, 0.4); }

        .path { font-family: 'Fira Code', monospace; color: #cbd5e1; font-size: 0.95rem; }

        .description {
            color: var(--text-muted);
            font-size: 0.95rem;
            margin-bottom: 1.25rem;
        }

        .section-label {
            font-size: 0.8rem;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--text-muted);
            margin-bottom: 0.5rem;
        }

        .code-block {
            background: var(--bg-code);
            border: 1px solid rgba(255, 255, 255, 0.05);
            border-radius: 10px;
            padding: 1rem;
            font-family: 'Fira Code', monospace;
            font-size: 0.85rem;
            color: #e2e8f0;
            overflow-x: auto;
            position: relative;
            margin-bottom: 1rem;
        }

        .copy-btn {
            position: absolute;
            top: 0.75rem;
            right: 0.75rem;
            background: rgba(255, 255, 255, 0.1);
            border: none;
            color: var(--text-main);
            padding: 0.35rem 0.65rem;
            border-radius: 6px;
            font-size: 0.75rem;
            cursor: pointer;
            transition: background 0.2s ease;
        }
        .copy-btn:hover { background: rgba(99, 102, 241, 0.5); }

        .btn-try {
            background: linear-gradient(135deg, var(--accent-primary), #4f46e5);
            color: white;
            border: none;
            padding: 0.6rem 1.2rem;
            border-radius: 8px;
            font-weight: 600;
            font-size: 0.85rem;
            cursor: pointer;
            transition: opacity 0.2s ease;
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
        }
        .btn-try:hover { opacity: 0.9; }

        .response-box {
            display: none;
            margin-top: 1rem;
        }

        footer {
            text-align: center;
            padding: 3rem;
            border-top: 1px solid var(--border);
            color: var(--text-muted);
            font-size: 0.9rem;
        }
    </style>
</head>
<body>

    <header class="header">
        <div class="brand">
            <div class="brand-icon">DOL</div>
            <div class="brand-title">Prevailing Wage API</div>
        </div>
        <div class="status-pill">
            <div class="status-dot"></div>
            System Online (Live)
        </div>
    </header>

    <main class="container">
        <section class="hero">
            <h1>U.S. DOL Prevailing Wage REST API</h1>
            <p>Automated 4-Tier Wage Determinations, O*NET SOC Metadata, Geographic FIPS Resolvers, and Immutable Immigration Audit Certificates.</p>
        </section>

        <div class="grid">

            <!-- API 1: Health Check -->
            <article class="card">
                <div class="card-header">
                    <div class="endpoint-title">
                        <span class="badge badge-get">GET</span>
                        <span class="path">/healthz</span>
                    </div>
                    <button class="btn-try" onclick="runApi('/healthz', 'GET', 'res-healthz')">▶ Try Live Request</button>
                </div>
                <p class="description">Checks service readiness, uptime, and system allocation statistics.</p>
                <div class="section-label">cURL Command</div>
                <div class="code-block">
                    <button class="copy-btn" onclick="copyCode(this)">Copy</button>
                    curl -s https://wage-determine.onrender.com/healthz
                </div>
                <div id="res-healthz" class="response-box">
                    <div class="section-label">Live Response</div>
                    <div class="code-block"><pre></pre></div>
                </div>
            </article>

            <!-- API 2: Wage Lookup -->
            <article class="card">
                <div class="card-header">
                    <div class="endpoint-title">
                        <span class="badge badge-get">GET</span>
                        <span class="path">/api/v1/wages/lookup</span>
                    </div>
                    <button class="btn-try" onclick="runApi('/api/v1/wages/lookup?soc_code=15-1252.00&zip_code=94103', 'GET', 'res-lookup')">▶ Try Live Request</button>
                </div>
                <p class="description">Fetches 4-Tier hourly and annual prevailing wages for a specific SOC Code and ZIP/Area.</p>
                <div class="section-label">cURL Command</div>
                <div class="code-block">
                    <button class="copy-btn" onclick="copyCode(this)">Copy</button>
                    curl -s "https://wage-determine.onrender.com/api/v1/wages/lookup?soc_code=15-1252.00&zip_code=94103"
                </div>
                <div id="res-lookup" class="response-box">
                    <div class="section-label">Live Response</div>
                    <div class="code-block"><pre></pre></div>
                </div>
            </article>

            <!-- API 3: Determine Wage Level -->
            <article class="card">
                <div class="card-header">
                    <div class="endpoint-title">
                        <span class="badge badge-post">POST</span>
                        <span class="path">/api/v1/wages/determine-level</span>
                    </div>
                    <button class="btn-try" onclick="runApiPost('/api/v1/wages/determine-level', sampleDetermineBody, 'res-determine')">▶ Try Live Request</button>
                </div>
                <p class="description">Automated 4-Tier evaluation engine based on job requirements (Education, Experience, Skills, Supervision). Calculates assigned wage level and generates an audit tracking number.</p>
                
                <div class="section-label">Request Body (JSON)</div>
                <div class="code-block">
{
  "soc_code": "15-1252.00",
  "zip_code": "94103",
  "job_title": "Senior Software Engineer",
  "education": { "required_degree": "Master" },
  "experience_months": 48,
  "special_skills": ["Go", "Kubernetes"],
  "supervises_employees": true,
  "number_of_subordinates": 3
}
                </div>

                <div class="section-label">cURL Command</div>
                <div class="code-block">
                    <button class="copy-btn" onclick="copyCode(this)">Copy</button>
curl -s -X POST https://wage-determine.onrender.com/api/v1/wages/determine-level \
  -H "Content-Type: application/json" \
  -d '{
    "soc_code": "15-1252.00",
    "zip_code": "94103",
    "job_title": "Senior Software Engineer",
    "education": { "required_degree": "Master" },
    "experience_months": 48,
    "special_skills": ["Go", "Kubernetes"],
    "supervises_employees": true,
    "number_of_subordinates": 3
  }'
                </div>
                <div id="res-determine" class="response-box">
                    <div class="section-label">Live Response</div>
                    <div class="code-block"><pre></pre></div>
                </div>
            </article>

            <!-- API 4: Location Resolver -->
            <article class="card">
                <div class="card-header">
                    <div class="endpoint-title">
                        <span class="badge badge-get">GET</span>
                        <span class="path">/api/v1/locations/resolve</span>
                    </div>
                    <button class="btn-try" onclick="runApi('/api/v1/locations/resolve?zip_code=10001', 'GET', 'res-loc')">▶ Try Live Request</button>
                </div>
                <p class="description">Resolves a 5-digit US ZIP Code to its corresponding FIPS County and BLS Area (MSA).</p>
                <div class="section-label">cURL Command</div>
                <div class="code-block">
                    <button class="copy-btn" onclick="copyCode(this)">Copy</button>
                    curl -s "https://wage-determine.onrender.com/api/v1/locations/resolve?zip_code=10001"
                </div>
                <div id="res-loc" class="response-box">
                    <div class="section-label">Live Response</div>
                    <div class="code-block"><pre></pre></div>
                </div>
            </article>

            <!-- API 5: Occupation Search -->
            <article class="card">
                <div class="card-header">
                    <div class="endpoint-title">
                        <span class="badge badge-get">GET</span>
                        <span class="path">/api/v1/occupations/search</span>
                    </div>
                    <button class="btn-try" onclick="runApi('/api/v1/occupations/search?q=Software', 'GET', 'res-occ')">▶ Try Live Request</button>
                </div>
                <p class="description">Searches O*NET Standard Occupational Classification codes and titles.</p>
                <div class="section-label">cURL Command</div>
                <div class="code-block">
                    <button class="copy-btn" onclick="copyCode(this)">Copy</button>
                    curl -s "https://wage-determine.onrender.com/api/v1/occupations/search?q=Software"
                </div>
                <div id="res-occ" class="response-box">
                    <div class="section-label">Live Response</div>
                    <div class="code-block"><pre></pre></div>
                </div>
            </article>

            <!-- API 6: Batch Lookup -->
            <article class="card">
                <div class="card-header">
                    <div class="endpoint-title">
                        <span class="badge badge-post">POST</span>
                        <span class="path">/api/v1/wages/batch-lookup</span>
                    </div>
                    <button class="btn-try" onclick="runApiPost('/api/v1/wages/batch-lookup', sampleBatchBody, 'res-batch')">▶ Try Live Request</button>
                </div>
                <p class="description">Queries prevailing wages in bulk across multiple SOC codes and geographic locations.</p>
                <div class="section-label">cURL Command</div>
                <div class="code-block">
                    <button class="copy-btn" onclick="copyCode(this)">Copy</button>
curl -s -X POST https://wage-determine.onrender.com/api/v1/wages/batch-lookup \
  -H "Content-Type: application/json" \
  -d '[
    {"soc_code": "15-1252.00", "zip_code": "94103"},
    {"soc_code": "15-1252.00", "zip_code": "10001"}
  ]'
                </div>
                <div id="res-batch" class="response-box">
                    <div class="section-label">Live Response</div>
                    <div class="code-block"><pre></pre></div>
                </div>
            </article>

        </div>
    </main>

    <footer>
        <p>U.S. Department of Labor Prevailing Wage Service &bull; Built in Golang &bull; Deployed on Render</p>
    </footer>

    <script>
        const sampleDetermineBody = {
            soc_code: "15-1252.00",
            zip_code: "94103",
            job_title: "Senior Software Engineer",
            education: { required_degree: "Master" },
            experience_months: 48,
            special_skills: ["Go", "Kubernetes"],
            supervises_employees: true,
            number_of_subordinates: 3
        };

        const sampleBatchBody = [
            { soc_code: "15-1252.00", zip_code: "94103" },
            { soc_code: "15-1252.00", zip_code: "10001" }
        ];

        function copyCode(btn) {
            const text = btn.parentElement.innerText.replace('Copy', '').trim();
            navigator.clipboard.writeText(text);
            btn.innerText = 'Copied!';
            setTimeout(() => btn.innerText = 'Copy', 2000);
        }

        async function runApi(url, method, targetId) {
            const box = document.getElementById(targetId);
            const pre = box.querySelector('pre');
            box.style.display = 'block';
            pre.innerText = 'Loading...';
            try {
                const res = await fetch(url);
                const data = await res.json();
                pre.innerText = JSON.stringify(data, null, 2);
            } catch (err) {
                pre.innerText = 'Error: ' + err.message;
            }
        }

        async function runApiPost(url, payload, targetId) {
            const box = document.getElementById(targetId);
            const pre = box.querySelector('pre');
            box.style.display = 'block';
            pre.innerText = 'Loading...';
            try {
                const res = await fetch(url, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });
                const data = await res.json();
                pre.innerText = JSON.stringify(data, null, 2);
            } catch (err) {
                pre.innerText = 'Error: ' + err.message;
            }
        }
    </script>
</body>
</html>`

func ServeLandingPage(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		// Return JSON API summary if client specifically requests JSON
		RespondJSON(w, http.StatusOK, map[string]interface{}{
			"service": "U.S. DOL Prevailing Wage Backend Service",
			"status":  "RUNNING",
			"version": "v1.0.0",
			"docs":    "https://github.com/satyamjiyadav/Wage-determine",
		})
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(landingHTML))
}
