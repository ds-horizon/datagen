import json
import os
import re
import subprocess
import sys
from typing import Optional
from urllib import request as urlrequest


def get_env(name: str, required: bool = False) -> Optional[str]:
    v = os.getenv(name)
    if required and not v:
        raise RuntimeError(f"Missing required env var: {name}")
    return v


GITHUB_TOKEN = get_env("GITHUB_TOKEN", True)
GITHUB_API_URL = get_env("GITHUB_API_URL") or "https://api.github.com"
GITHUB_REPOSITORY = get_env("GITHUB_REPOSITORY", True)  # owner/repo
GITHUB_SHA = get_env("GITHUB_SHA", True)
COVERAGE_TOTAL = get_env("COVERAGE_TOTAL")
COVERAGE_DETAILS = get_env("COVERAGE_DETAILS")
MIN_COVERAGE = get_env("MIN_COVERAGE") or "N/A"

owner, repo = GITHUB_REPOSITORY.split("/")


def github_request(method: str, path: str, body: Optional[dict] = None):
    url = f"{GITHUB_API_URL}{path}"
    data = None
    headers = {
        "Authorization": f"Bearer {GITHUB_TOKEN}",
        "Accept": "application/vnd.github+json",
        "User-Agent": "coverage-comment-script",
    }
    if body is not None:
        data = json.dumps(body).encode("utf-8")
        headers["Content-Type"] = "application/json"
    req = urlrequest.Request(url, data=data, headers=headers, method=method)
    try:
        with urlrequest.urlopen(req) as resp:
            raw = resp.read().decode("utf-8")
            try:
                return json.loads(raw) if raw else None
            except Exception:
                return None
    except Exception as e:
        raise RuntimeError(f"Request failed: {e}")


def build_details() -> str:
    if COVERAGE_DETAILS:
        return COVERAGE_DETAILS
    try:
        out = subprocess.check_output(
            ["bash", "-lc", "go tool cover -func=coverage.out"],
            text=True,
        )
        return out
    except Exception:
        return 'Coverage details unavailable (could not run "go tool cover -func").'


def make_body(details: str) -> str:
    marker = "<!-- coverage-comment -->"
    if COVERAGE_TOTAL:
        body = f"{marker}\n**Coverage**: {COVERAGE_TOTAL}% (min: {MIN_COVERAGE}%)"
    else:
        body = f"{marker}\n**Coverage**: unavailable (tests failed before coverage generation)"

    table = ""
    try:
        rows = []
        for raw in (details or "").splitlines():
            line = (raw or "").strip()
            if not line:
                continue
            if line.startswith("total:"):
                continue
            if "%" not in line:
                continue
            m_pct = re.search(r"([0-9]+(?:\\.[0-9]+)?)%\\s*$", line)
            if not m_pct:
                continue
            pct = m_pct.group(1) + "%"
            left = re.sub(r"\\s*[0-9]+(?:\\.[0-9]+)?%\\s*$", "", line).strip()

            file = ""
            fn = ""
            m = re.match(r"^(.+?):(?:\\d+:)?\\s*(.+)$", left)
            if m:
                file = m.group(1).strip()
                fn = (m.group(2) or "").strip()
            else:
                idx = left.find(":")
                if idx != -1:
                    file = left[:idx].strip()
                    fn = left[idx + 1 :].strip()
            fn = re.sub(r"\\s+\\d+$", "", fn)
            fn = re.sub(r":$", "", fn)
            if not file or not fn:
                continue
            rows.append((file, fn, pct))

        if rows:
            max_rows = 300
            slice_rows = rows[:max_rows]
            table_lines = ["| File | Function | Coverage |", "| --- | --- | ---: |"]
            table_lines.extend([f"| {f} | {fn} | {p} |" for f, fn, p in slice_rows])
            if len(rows) > max_rows:
                table_lines.append("| … | … | … |")
            table = "\n".join(table_lines)
    except Exception:
        table = ""

    if table:
        body += f"\n\n<details>\n<summary>Coverage details</summary>\n\n{table}\n</details>"
    else:
        body += f"\n\n<details>\n<summary>Coverage details</summary>\n\n```\n{details}\n```\n</details>"
    return body


def main():
    details = build_details()
    body = make_body(details)
    sha = GITHUB_SHA

    pr_number = None
    try:
        prs = github_request("GET", f"/repos/{owner}/{repo}/commits/{sha}/pulls")
        if isinstance(prs, list) and prs:
            open_pr = next((p for p in prs if p.get("state") == "open"), None)
            pr_number = (open_pr or prs[0]).get("number")
    except Exception:
        pass

    marker = "<!-- coverage-comment -->"
    if pr_number:
        comments = github_request("GET", f"/repos/{owner}/{repo}/issues/{pr_number}/comments")
        existing = next((c for c in (comments or []) if c.get("body") and marker in c.get("body")), None)
        if existing:
            github_request("PATCH", f"/repos/{owner}/{repo}/issues/comments/{existing.get('id')}", {"body": body})
        else:
            github_request("POST", f"/repos/{owner}/{repo}/issues/{pr_number}/comments", {"body": body})
    else:
        updated = False
        try:
            commit_comments = github_request("GET", f"/repos/{owner}/{repo}/commits/{sha}/comments")
            existing = next((c for c in (commit_comments or []) if c.get("body") and marker in c.get("body")), None)
            if existing:
                github_request("PATCH", f"/repos/{owner}/{repo}/comments/{existing.get('id')}", {"body": body})
                updated = True
        except Exception:
            pass
        if not updated:
            github_request("POST", f"/repos/{owner}/{repo}/commits/{sha}/comments", {"body": body})


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(e, file=sys.stderr)
        sys.exit(1)


