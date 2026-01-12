package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"text/template"
	"time"
)

// VULNERABILITY: Hardcoded credentials
const (
	dbPassword = "SuperSecret123!"
	apiKey     = "sk-1234567890abcdef"
)

// VULNERABILITY: SQL Injection
func getUserByID(db *sql.DB, userID string) error {
	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s'", userID)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

// VULNERABILITY: Command Injection
func runCommand(userInput string) error {
	cmd := exec.Command("/bin/sh", "-c", "ls "+userInput)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// VULNERABILITY: Path Traversal
func readFile(filename string) ([]byte, error) {
	filepath := "/var/www/uploads/" + filename
	return ioutil.ReadFile(filepath)
}

// VULNERABILITY: XSS using text/template
func renderTemplate(w http.ResponseWriter, userContent string) error {
	tmpl := template.Must(template.New("page").Parse(`
		<html>
		<body>
			<div>{{.Content}}</div>
		</body>
		</html>
	`))
	return tmpl.Execute(w, map[string]string{"Content": userContent})
}

// VULNERABILITY: Weak cryptography (MD5)
func hashPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// VULNERABILITY: Insecure random number generation
func generateToken() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Int())
}

// VULNERABILITY: Race condition
var counter int

func incrementCounter() {
	counter++
}

// VULNERABILITY: Missing TLS, no authentication
func startServer() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("id")
		getUserByID(nil, userID)
	})

	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		runCommand(cmd)
	})

	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("name")
		content, _ := readFile(filename)
		w.Write(content)
	})

	// Missing TLS
	http.ListenAndServe(":8080", nil)
}

func main() {
	// Launch goroutines without synchronization
	for i := 0; i < 100; i++ {
		go incrementCounter()
	}

	startServer()
}
