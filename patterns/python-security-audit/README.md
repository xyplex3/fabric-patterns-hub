# python-security-audit

A comprehensive Fabric pattern for performing security audits on Python codebases. This pattern identifies vulnerabilities, security anti-patterns, and potential exploits, providing actionable remediation guidance based on industry best practices.

## Features

- **Comprehensive Vulnerability Detection**: Scans for OWASP Top 10 and Python-specific vulnerabilities
- **Severity-Based Reporting**: Categorizes findings as CRITICAL, HIGH, MEDIUM, LOW, or INFO
- **Actionable Remediation**: Provides vulnerable code snippets alongside secure implementations
- **Knowledge Base Integration**: Leverages extensive Python security knowledge covering:
  - Injection attacks (SQL, Command, Template/SSTI)
  - Cross-Site Scripting (XSS)
  - Insecure deserialization (pickle, yaml.load, eval)
  - Path traversal vulnerabilities
  - Server-Side Request Forgery (SSRF)
  - Web framework security (Django, Flask, FastAPI)
  - Cryptography standards and password hashing
  - Dependency security and supply chain attacks
  - Secrets management
- **Tool Recommendations**: Suggests appropriate security tools (bandit, pip-audit, safety, semgrep)
- **Tailored Checklists**: Generates security checklists specific to your codebase
- **Quick Wins**: Identifies high-impact, low-effort security improvements

## Installation

This pattern requires [Fabric](https://github.com/danielmiessler/fabric) to be installed.

```bash
# Clone the patterns hub
git clone https://github.com/yourusername/fabric-patterns-hub.git

# Link or copy to your Fabric patterns directory
ln -s $(pwd)/fabric-patterns-hub/patterns/python-security-audit ~/.config/fabric/patterns/python-security-audit

# Or copy directly
cp -r fabric-patterns-hub/patterns/python-security-audit ~/.config/fabric/patterns/
```

## Usage

### Audit a Single Python File

```bash
cat app.py | fabric -p python-security-audit
```

### Audit an Entire Package

```bash
cat src/**/*.py | fabric -p python-security-audit
```

### Audit Multiple Files

```bash
cat $(find . -name "*.py" -not -path "*/venv/*") | fabric -p python-security-audit
```

### Audit with Security Tool Output

```bash
# Combine your code with bandit output for enhanced analysis
bandit -r . -f txt > bandit.txt
cat app.py bandit.txt | fabric -p python-security-audit
```

### Audit with Dependency Scan

```bash
# Combine code audit with pip-audit findings
pip-audit -r requirements.txt > pip-audit.txt
cat app.py requirements.txt pip-audit.txt | fabric -p python-security-audit
```

### Save Report to File

```bash
cat app.py | fabric -p python-security-audit > security-audit-report.md
```

### Audit Specific Module

```bash
find . -path "*/auth/*.py" | xargs cat | fabric -p python-security-audit
```

## Example Input

```bash
cat << 'EOF' | fabric -p python-security-audit
import sqlite3
from flask import Flask, request

app = Flask(__name__)
app.secret_key = "dev-secret"

def get_user(username):
    conn = sqlite3.connect("app.db")
    query = f"SELECT * FROM users WHERE username='{username}'"
    return conn.execute(query).fetchone()

@app.route('/user')
def user():
    username = request.args.get('username')
    return str(get_user(username))

if __name__ == '__main__':
    app.run(debug=True)
EOF
```

## Example Output

The pattern generates a comprehensive security audit report including:

```markdown
# Python Security Audit Report

## Executive Summary

**Code Type:** Web Application (Flask)
**Audit Date:** 2026-04-14
**Python Version:** 3.10+

### Findings Overview

| Severity | Count | Status |
|----------|-------|--------|
| CRITICAL | 1 | 🔴 Immediate action required |
| HIGH | 0 | 🟠 Address urgently |
| MEDIUM | 2 | 🟡 Schedule remediation |
| LOW | 0 | 🔵 Consider fixing |
| INFO | 1 | ℹ️ Recommendations |

**Total Issues:** 4

### Most Critical Issues

1. SQL Injection vulnerability in get_user function
2. Hardcoded Flask secret_key
3. Debug mode enabled in production

## CRITICAL Severity Findings

### SQL Injection - CWE-89

**File:** `app.py:8`

**Description:**
The code uses an f-string to construct a SQL query with user input,
allowing attackers to inject SQL and access or modify arbitrary data.

**Vulnerable Code:**
```python
query = f"SELECT * FROM users WHERE username='{username}'"
conn.execute(query)
```

**Security Impact:**
- Attackers can bypass authentication
- Unauthorized data access or modification
- Potential full database compromise

**Secure Implementation:**
```python
query = "SELECT * FROM users WHERE username = ?"
conn.execute(query, (username,))
```

**Remediation Steps:**
1. Replace f-string query with parameterized statement
2. Use `?` placeholders for SQLite or `%s` for psycopg2
3. Validate and sanitize all user input before use
4. Run `bandit -t B608 app.py` to detect similar issues

**Detection Tool:** `bandit -t B608 .`
...
```

## When to Use

Use this pattern when you need to:

- **Security audit** Python codebases before production deployment
- **Code review** with security focus for pull requests
- **Compliance** checks against security standards
- **Vulnerability assessment** of third-party Python libraries
- **Security training** by analyzing example vulnerable code
- **Penetration testing** preparation and threat modeling
- **Regular security scanning** as part of CI/CD pipeline

## When NOT to Use

This pattern is not suitable for:

- **Automated security gates** - While comprehensive, this pattern provides educational context best reviewed by humans
- **Real-time production monitoring** - Use dedicated security monitoring tools instead
- **Performance profiling** - This focuses on security, not performance
- **Functional correctness** - This checks security, not business logic
- **Dependency vulnerability scanning only** - Use pip-audit directly for faster dependency-only scans

## Integration with Security Tools

This pattern complements but doesn't replace dedicated security tools:

### bandit (Static Security Analysis)
```bash
# Run bandit first, then analyze with pattern
bandit -r . -f txt > bandit.txt
cat $(find . -name "*.py") bandit.txt | fabric -p python-security-audit
```

### pip-audit (Dependency Scanning)
```bash
# Combine dependency audit with code audit
pip-audit -r requirements.txt > pip-audit.txt
cat app.py requirements.txt pip-audit.txt | fabric -p python-security-audit
```

### safety (Dependency Database)
```bash
# Include safety output
safety check -r requirements.txt > safety.txt
cat app.py safety.txt | fabric -p python-security-audit
```

### semgrep (Advanced SAST)
```bash
# Use semgrep Python security rules
semgrep --config=p/python-security . --json > semgrep.json
cat app.py semgrep.json | fabric -p python-security-audit
```

## Best Practices

1. **Run Regularly**: Integrate into your development workflow
   ```bash
   # Git pre-commit hook
   git diff HEAD --name-only --diff-filter=ACM | grep '\.py$' | xargs cat | fabric -p python-security-audit
   ```

2. **Prioritize Findings**: Address CRITICAL and HIGH severity issues first

3. **Combine with Tools**: Use alongside bandit and pip-audit for comprehensive coverage

4. **Review Context**: This pattern provides educational context — review findings with your security requirements

5. **Update Knowledge**: Keep the python-security-guide.md knowledge base updated with latest CVEs

6. **Customize Checks**: Modify system.md to add organization-specific security requirements

7. **Track Remediation**: Save reports to track progress over time
   ```bash
   cat app.py | fabric -p python-security-audit > audits/$(date +%Y-%m-%d)-audit.md
   ```

## Understanding Severity Levels

- **CRITICAL**: Immediate remote exploitation possible (RCE via eval/pickle, SQL injection, command injection, auth bypass)
- **HIGH**: Significant security risk requiring urgent attention (XSS, path traversal, SSRF, weak crypto)
- **MEDIUM**: Security weaknesses that should be fixed (missing validation, CSRF, debug mode)
- **LOW**: Minor issues or information disclosure (verbose errors, missing headers)
- **INFO**: Best practice recommendations and hardening opportunities

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Python Security Audit
on: [push, pull_request]

jobs:
  security-audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: '3.12'

      - name: Install security tools
        run: |
          pip install bandit pip-audit safety

      - name: Install Fabric
        run: |
          curl -sL https://get.fabric.sh | bash
          export PATH="$HOME/.local/bin:$PATH"

      - name: Install Pattern
        run: |
          mkdir -p ~/.config/fabric/patterns
          cp -r patterns/python-security-audit ~/.config/fabric/patterns/

      - name: Run bandit
        run: bandit -r . -f txt --exclude ./tests,./venv > bandit.txt || true

      - name: Run pip-audit
        run: pip-audit -r requirements.txt > pip-audit.txt || true

      - name: Run Security Audit
        run: |
          find . -name "*.py" -not -path "*/venv/*" -not -path "*/.git/*" | \
            xargs cat | \
            cat - bandit.txt pip-audit.txt | \
            fabric -p python-security-audit > security-audit.md

      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: security-audit-report
          path: security-audit.md

      - name: Check for Critical Issues
        run: |
          if grep -q "## CRITICAL" security-audit.md; then
            echo "::error::Critical security issues found"
            exit 1
          fi
```

## Troubleshooting

### Pattern Not Found
```bash
# Verify pattern installation
ls ~/.config/fabric/patterns/python-security-audit

# Refresh Fabric patterns
fabric --update
```

### Large Codebase Timeout
```bash
# Audit modules individually
for module in src/auth src/api src/models; do
  echo "Auditing $module"
  find "$module" -name "*.py" | xargs cat | fabric -p python-security-audit > "audit-$(basename $module).md"
done
```

### False Positives
The pattern is designed to be thorough. Review each finding in context — some may not apply to your specific use case. Consider:
- Input validation already performed elsewhere in the call chain
- Code running in trusted, internal-only contexts
- False positives from security tools included in input

## Related Patterns

- **go-security-audit**: Equivalent pattern for Go codebases
- **go-cobra**: Go CLI best practices review
- **ansible-review**: Automated Ansible code review

## Knowledge Base

This pattern includes `python-security-guide.md`, a comprehensive security reference covering:

- Security tools and analysis (bandit, pip-audit, safety, semgrep)
- Common vulnerabilities with CWE references
- Insecure deserialization (pickle, yaml, eval)
- Web framework security for Django, Flask, and FastAPI
- Cryptography best practices and password hashing
- Dependency security and supply chain attacks
- Secrets management patterns
- Complete security checklists

The knowledge base covers established Python security research and common CVE patterns. Update `python-security-guide.md` manually to reflect newly discovered vulnerabilities.

## Contributing

To improve this pattern:

1. Update `python-security-guide.md` with new vulnerabilities and CVEs
2. Enhance `system.md` with additional vulnerability detection patterns
3. Add test cases in `test-vulnerable.py` and `test-secure.py`
4. Submit pull requests with improvements

## References

- [Python Security Documentation](https://python.org/dev/security/)
- [OWASP Python Security Project](https://owasp.org/www-project-python-security/)
- [Bandit Documentation](https://bandit.readthedocs.io/)
- [pip-audit GitHub Repository](https://github.com/pypa/pip-audit)
- [Python Cryptography Library](https://cryptography.io/)
- [OWASP Top 10](https://owasp.org/Top10/)
- [PyPA Advisory Database](https://github.com/pypa/advisory-database)

## License

This pattern is part of the Fabric Patterns Hub and follows the same license terms.

## Support

For issues, questions, or contributions:
- Open an issue in the fabric-patterns-hub repository
- Refer to the comprehensive `python-security-guide.md` for detailed security information
- Check official Python security resources for the latest updates
