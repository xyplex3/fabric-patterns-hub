"""
Secure Python code demonstrating best practices for the python-security-audit pattern.
Use this as a reference for secure implementations.
"""

import hashlib
import logging
import os
import re
import secrets
import sqlite3
from functools import wraps
from pathlib import Path
from typing import Optional

import bcrypt
import yaml
from cryptography.fernet import Fernet
from flask import Flask, abort, g, jsonify, request, render_template
from markupsafe import escape

# SECURE: Load secrets from environment variables, never hardcode
app = Flask(__name__)
app.secret_key = os.environ["FLASK_SECRET_KEY"]
app.config.update(
    SESSION_COOKIE_SECURE=True,
    SESSION_COOKIE_HTTPONLY=True,
    SESSION_COOKIE_SAMESITE="Lax",
)

# SECURE: Parameterized SQL queries
def get_user_by_id(user_id: str) -> Optional[tuple]:
    conn = sqlite3.connect("app.db")
    cursor = conn.cursor()
    cursor.execute("SELECT id, name, email FROM users WHERE id = ?", (user_id,))
    return cursor.fetchone()

# SECURE: Parameterized search
def search_users(name: str) -> list:
    conn = sqlite3.connect("app.db")
    cursor = conn.cursor()
    cursor.execute("SELECT id, name FROM users WHERE name = ?", (name,))
    return cursor.fetchall()

# SECURE: No shell=True, input validated, arguments as list
def run_ping(host: str) -> str:
    if not re.match(r'^[a-zA-Z0-9.\-]+$', host):
        raise ValueError("Invalid hostname format")
    import subprocess
    result = subprocess.run(
        ["ping", "-c", "4", host],
        capture_output=True,
        text=True,
        check=True,
        timeout=10,
    )
    return result.stdout

# SECURE: Path traversal protection
BASE_UPLOAD_DIR = Path("/var/www/uploads").resolve()

def read_file(filename: str) -> str:
    requested = (BASE_UPLOAD_DIR / filename).resolve()
    if not str(requested).startswith(str(BASE_UPLOAD_DIR) + os.sep):
        raise PermissionError("Access denied: path traversal detected")
    with open(requested) as f:
        return f.read()

# SECURE: Safe YAML loading
def load_config(config_text: str) -> dict:
    return yaml.safe_load(config_text)

# SECURE: No eval(); use ast.literal_eval for safe parsing of Python literals
def parse_literal(data: str):
    import ast
    return ast.literal_eval(data)

# SECURE: Strong password hashing with bcrypt
def hash_password(password: str) -> str:
    return bcrypt.hashpw(password.encode(), bcrypt.gensalt(rounds=12)).decode()

def verify_password(password: str, hashed: str) -> bool:
    return bcrypt.checkpw(password.encode(), hashed.encode())

# SECURE: Cryptographically secure random token
def generate_token() -> str:
    return secrets.token_urlsafe(32)

def generate_session_id() -> str:
    return secrets.token_hex(32)

# SECURE: XSS protection — escape user input before rendering
@app.route("/greet")
def greet():
    name = request.args.get("name", "")
    safe_name = escape(name)
    return f"<h1>Hello {safe_name}!</h1>"

# SECURE: Use render_template with a file (Jinja2 auto-escapes)
@app.route("/profile")
def profile():
    username = request.args.get("username", "")
    return render_template("profile.html", username=username)

# SECURE: TLS verification enabled (default in requests)
def fetch_data(url: str) -> dict:
    import requests
    import urllib.parse
    parsed = urllib.parse.urlparse(url)
    if parsed.scheme not in ("https",):
        raise ValueError("Only HTTPS URLs are allowed")
    return requests.get(url, timeout=10).json()

# SECURE: Authentication middleware
def require_auth(f):
    @wraps(f)
    def decorated(*args, **kwargs):
        token = request.headers.get("Authorization", "")
        if not token.startswith("Bearer "):
            abort(401)
        bearer_token = token[7:]
        if not validate_token(bearer_token):
            abort(401)
        return f(*args, **kwargs)
    return decorated

def validate_token(token: str) -> bool:
    expected = os.environ.get("API_TOKEN", "")
    return secrets.compare_digest(token, expected)

# SECURE: Protected admin endpoint with authentication
@app.route("/admin/users")
@require_auth
def list_all_users():
    conn = sqlite3.connect("app.db")
    cursor = conn.cursor()
    cursor.execute("SELECT id, name FROM users")
    users = [{"id": row[0], "name": row[1]} for row in cursor.fetchall()]
    return jsonify(users)

# SECURE: Redact sensitive data from logs
class SensitiveDataFilter(logging.Filter):
    PATTERNS = [
        (re.compile(r'password=\S+', re.I), 'password=***'),
        (re.compile(r'token=\S+', re.I), 'token=***'),
    ]

    def filter(self, record: logging.LogRecord) -> bool:
        for pattern, replacement in self.PATTERNS:
            record.msg = pattern.sub(replacement, str(record.msg))
        return True

logging.basicConfig(level=logging.INFO)
logging.getLogger().addFilter(SensitiveDataFilter())
logger = logging.getLogger(__name__)

def authenticate(username: str, password: str) -> bool:
    logger.info("Login attempt: username=%s", username)  # Password NOT logged
    user = get_user_by_id(username)
    if not user:
        return False
    stored_hash = os.environ.get("USER_HASH", "")
    return verify_password(password, stored_hash)

# SECURE: SSRF prevention — validate and restrict allowed URLs
ALLOWED_HOSTS = {"api.example.com", "data.example.com"}

@app.route("/proxy")
@require_auth
def proxy():
    import requests
    import urllib.parse
    url = request.args.get("url", "")
    parsed = urllib.parse.urlparse(url)
    if parsed.scheme != "https" or parsed.hostname not in ALLOWED_HOSTS:
        abort(400, description="URL not allowed")
    response = requests.get(url, timeout=5)
    return jsonify(response.json())

# SECURE: Symmetric encryption using Fernet
def encrypt_data(plaintext: str) -> bytes:
    key = os.environ["ENCRYPTION_KEY"].encode()
    fernet = Fernet(key)
    return fernet.encrypt(plaintext.encode())

def decrypt_data(ciphertext: bytes) -> str:
    key = os.environ["ENCRYPTION_KEY"].encode()
    fernet = Fernet(key)
    return fernet.decrypt(ciphertext).decode()

if __name__ == "__main__":
    # SECURE: debug mode controlled by environment variable, never hardcoded True
    debug_mode = os.environ.get("FLASK_DEBUG", "False").lower() == "true"
    app.run(debug=debug_mode, host="127.0.0.1", port=5000)
