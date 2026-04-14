"""
Deliberately vulnerable Python code for testing the python-security-audit pattern.
DO NOT use this code in production.
"""

import hashlib
import os
import pickle
import random
import sqlite3
import subprocess
import yaml
from flask import Flask, request, render_template_string
from markupsafe import Markup

app = Flask(__name__)

# VULNERABILITY: Hardcoded credentials
app.secret_key = "dev-secret-key-hardcoded"
DB_PASSWORD = "SuperSecret123!"
API_KEY = "sk-1234567890abcdef"
JWT_SECRET = "my-hardcoded-jwt-secret"

# VULNERABILITY: SQL Injection
def get_user_by_id(user_id: str):
    conn = sqlite3.connect("app.db")
    cursor = conn.cursor()
    query = f"SELECT * FROM users WHERE id = '{user_id}'"
    cursor.execute(query)
    return cursor.fetchone()

# VULNERABILITY: SQL Injection via % formatting
def search_users(name: str):
    conn = sqlite3.connect("app.db")
    cursor = conn.cursor()
    query = "SELECT * FROM users WHERE name = '%s'" % name
    cursor.execute(query)
    return cursor.fetchall()

# VULNERABILITY: Command Injection
def run_ping(host: str):
    os.system(f"ping -c 4 {host}")

# VULNERABILITY: Command Injection via subprocess shell=True
def list_files(directory: str):
    result = subprocess.run(f"ls {directory}", shell=True, capture_output=True, text=True)
    return result.stdout

# VULNERABILITY: Path Traversal
def read_file(filename: str) -> str:
    path = "/var/www/uploads/" + filename
    with open(path) as f:
        return f.read()

# VULNERABILITY: Insecure deserialization (pickle)
def load_session(data: bytes):
    return pickle.loads(data)

# VULNERABILITY: Unsafe YAML loading
def load_config(config_text: str):
    return yaml.load(config_text)

# VULNERABILITY: eval() with user input
def calculate(expression: str):
    return eval(expression)

# VULNERABILITY: Weak password hashing (MD5)
def hash_password(password: str) -> str:
    return hashlib.md5(password.encode()).hexdigest()

# VULNERABILITY: SHA-1 without salt
def legacy_hash(password: str) -> str:
    return hashlib.sha1(password.encode()).hexdigest()

# VULNERABILITY: Insecure random number generation
def generate_token() -> str:
    return str(random.randint(100000, 999999))

def generate_session_id() -> str:
    return "".join(random.choices("abcdefghijklmnopqrstuvwxyz0123456789", k=32))

# VULNERABILITY: XSS via Markup
@app.route("/greet")
def greet():
    name = request.args.get("name", "")
    return Markup(f"<h1>Hello {name}</h1>")

# VULNERABILITY: Server-Side Template Injection (SSTI)
@app.route("/template")
def ssti_endpoint():
    name = request.args.get("name", "")
    template = f"<h1>Hello {name}!</h1>"
    return render_template_string(template)

# VULNERABILITY: Missing TLS verification
def fetch_data(url: str):
    import requests
    return requests.get(url, verify=False).json()

# VULNERABILITY: Debug mode enabled
@app.route("/debug")
def debug_info():
    import sys
    # Exposes internal state
    return str(sys.path) + str(os.environ)

# VULNERABILITY: Missing authentication on sensitive endpoint
@app.route("/admin/users")
def list_all_users():
    conn = sqlite3.connect("app.db")
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM users")
    return str(cursor.fetchall())

# VULNERABILITY: Logging sensitive data
import logging
logging.basicConfig(level=logging.DEBUG)

def authenticate(username: str, password: str) -> bool:
    logging.info(f"Login attempt: username={username} password={password}")
    return True

# VULNERABILITY: SSRF - user-controlled URL
@app.route("/proxy")
def proxy():
    import requests
    url = request.args.get("url")
    return requests.get(url).text

if __name__ == "__main__":
    # VULNERABILITY: debug=True in production
    app.run(debug=True, host="0.0.0.0", port=5000)
