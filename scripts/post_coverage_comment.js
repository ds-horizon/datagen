const { execSync } = require('child_process');
const https = require('https');

function getEnv(name, required = false) {
  const v = process.env[name];
  if (required && !v) {
    throw new Error(`Missing required env var: ${name}`);
  }
  return v;
}

const GITHUB_TOKEN = getEnv('GITHUB_TOKEN', true);
const GITHUB_API_URL = getEnv('GITHUB_API_URL') || 'https://api.github.com';
const GITHUB_REPOSITORY = getEnv('GITHUB_REPOSITORY', true); // owner/repo
const GITHUB_SHA = getEnv('GITHUB_SHA', true);
const COVERAGE_TOTAL = getEnv('COVERAGE_TOTAL');
const COVERAGE_DETAILS = getEnv('COVERAGE_DETAILS');
const MIN_COVERAGE = getEnv('MIN_COVERAGE') || 'N/A';

const [owner, repo] = GITHUB_REPOSITORY.split('/');

function request(method, path, body) {
  const data = body ? Buffer.from(JSON.stringify(body)) : null;
  const opts = {
    method,
    headers: {
      'Authorization': `Bearer ${GITHUB_TOKEN}`,
      'Accept': 'application/vnd.github+json',
      'User-Agent': 'coverage-comment-script',
    },
  };
  return new Promise((resolve, reject) => {
    const req = https.request(`${GITHUB_API_URL}${path}`, opts, (res) => {
      const chunks = [];
      res.on('data', (c) => chunks.push(c));
      res.on('end', () => {
        const raw = Buffer.concat(chunks).toString('utf8');
        let json;
        try { json = raw ? JSON.parse(raw) : null; } catch (_) { json = null; }
        if (res.statusCode && res.statusCode >= 200 && res.statusCode < 300) {
          resolve(json);
        } else {
          reject(new Error(`Request failed ${res.statusCode}: ${raw}`));
        }
      });
    });
    req.on('error', reject);
    if (data) req.write(data);
    req.end();
  });
}

function buildDetails() {
  if (COVERAGE_DETAILS) return COVERAGE_DETAILS;
  try {
    return execSync('go tool cover -func=coverage.out', { encoding: 'utf8' });
  } catch (_) {
    return 'Coverage details unavailable (could not run "go tool cover -func").';
  }
}

function makeBody(details) {
  const marker = '<!-- coverage-comment -->';
  let body = `${marker}\n**Coverage**: ${COVERAGE_TOTAL || 'N/A'}% (min: ${MIN_COVERAGE}%)`;
  if (!COVERAGE_TOTAL) {
    body = `${marker}\n**Coverage**: unavailable (tests failed before coverage generation)`;
  }
  let table = '';
  try {
    const lines = (details || '').split('\n');
    const rows = [];
    for (const raw of lines) {
      const line = (raw || '').trim();
      if (!line) continue;
      if (line.startsWith('total:')) continue;
      if (!line.includes('%')) continue;
      const pm = line.match(/([0-9]+(?:\.[0-9]+)?)%\s*$/);
      if (!pm) continue;
      const pct = pm[1] + '%';
      const left = line.replace(/\s*[0-9]+(?:\.[0-9]+)?%\s*$/, '').trim();
      let file = '';
      let fn = '';
      const m = left.match(/^(.+?):(?:\d+:)?\s*(.+)$/);
      if (m) {
        file = m[1].trim();
        fn = (m[2] || '').trim();
      } else {
        const idx = left.indexOf(':');
        if (idx !== -1) {
          file = left.slice(0, idx).trim();
          fn = left.slice(idx + 1).trim();
        }
      }
      fn = fn.replace(/\s+\d+$/, '');
      fn = fn.replace(/:$/, '');
      if (!file || !fn) continue;
      rows.push({ file, fn, pct });
    }
    if (rows.length > 0) {
      const maxRows = 300;
      const slice = rows.slice(0, maxRows);
      table += `| File | Function | Coverage |\n`;
      table += `| --- | --- | ---: |\n`;
      for (const r of slice) {
        table += `| ${r.file} | ${r.fn} | ${r.pct} |\n`;
      }
      if (rows.length > maxRows) {
        table += `| … | … | … |\n`;
      }
    }
  } catch (_) {}
  if (table) {
    body += `\n\n<details>\n<summary>Coverage details</summary>\n\n${table}\n</details>`;
  } else {
    body += `\n\n<details>\n<summary>Coverage details</summary>\n\n\`\`\`\n${details}\n\`\`\`\n</details>`;
  }
  return body;
}

async function run() {
  const details = buildDetails();
  const body = makeBody(details);
  const sha = GITHUB_SHA;

  let prNumber = null;
  try {
    const prs = await request('GET', `/repos/${owner}/${repo}/commits/${sha}/pulls`);
    if (Array.isArray(prs) && prs.length > 0) {
      const open = prs.find(p => p.state === 'open');
      prNumber = (open || prs[0]).number;
    }
  } catch (e) {
  }

  const marker = '<!-- coverage-comment -->';
  if (prNumber) {
    const comments = await request('GET', `/repos/${owner}/${repo}/issues/${prNumber}/comments`);
    const existing = (comments || []).find(c => c.body && c.body.includes(marker));
    if (existing) {
      await request('PATCH', `/repos/${owner}/${repo}/issues/comments/${existing.id}`, { body });
    } else {
      await request('POST', `/repos/${owner}/${repo}/issues/${prNumber}/comments`, { body });
    }
  } else {
    let updated = false;
    try {
      const commitComments = await request('GET', `/repos/${owner}/${repo}/commits/${sha}/comments`);
      const existing = (commitComments || []).find(c => c.body && c.body.includes(marker));
      if (existing) {
        await request('PATCH', `/repos/${owner}/${repo}/comments/${existing.id}`, { body });
        updated = true;
      }
    } catch (_) {}
    if (!updated) {
      await request('POST', `/repos/${owner}/${repo}/commits/${sha}/comments`, { body });
    }
  }
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});