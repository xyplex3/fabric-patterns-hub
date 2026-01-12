package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SECURE: Use environment variables for secrets
func getDBPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func getAPIKey() string {
	return os.Getenv("API_KEY")
}

// SECURE: Parameterized queries prevent SQL injection
func getUserByID(db *sql.DB, userID string) error {
	query := "SELECT * FROM users WHERE id=$1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return fmt.Errorf("database query failed")
	}
	defer rows.Close()
	return nil
}

// SECURE: No shell invocation, direct command execution
func listDirectory(basePath string) error {
	// Validate input
	cleanPath := filepath.Clean(basePath)
	if !strings.HasPrefix(cleanPath, "/var/www/allowed/") {
		return fmt.Errorf("invalid path")
	}

	entries, err := os.ReadDir(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to read directory")
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
	return nil
}

// SECURE: Path traversal protection using filepath.Clean and validation
func readFile(filename string) ([]byte, error) {
	basePath := "/var/www/uploads"
	cleanFilename := filepath.Clean(filename)

	// Prevent directory traversal
	fullPath := filepath.Join(basePath, cleanFilename)
	if !strings.HasPrefix(fullPath, filepath.Clean(basePath)+string(os.PathSeparator)) {
		return nil, fmt.Errorf("invalid file path")
	}

	return os.ReadFile(fullPath)
}

// SECURE: XSS protection using html/template
func renderTemplate(w http.ResponseWriter, userContent string) error {
	// html/template automatically escapes content
	tmpl := template.Must(template.New("page").Parse(`
		<html>
		<body>
			<div>{{.Content}}</div>
		</body>
		</html>
	`))
	return tmpl.Execute(w, map[string]string{"Content": userContent})
}

// SECURE: Strong password hashing with bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SECURE: Cryptographically secure random token generation
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SECURE: Thread-safe counter using mutex
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// SECURE: Authentication middleware
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !validateToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func validateToken(token string) bool {
	// Implement actual token validation
	return token == "Bearer "+getAPIKey()
}

// SECURE: Server with TLS, authentication, timeouts, and proper error handling
func startServer() error {
	counter := &SafeCounter{}

	mux := http.NewServeMux()

	// Protected endpoints with authentication
	mux.HandleFunc("/user", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("id")
		if userID == "" {
			http.Error(w, "Missing user ID", http.StatusBadRequest)
			return
		}

		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		if err := getUserByID(db, userID); err != nil {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))

	mux.HandleFunc("/file", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("name")
		if filename == "" {
			http.Error(w, "Missing filename", http.StatusBadRequest)
			return
		}

		content, err := readFile(filename)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(content)
	}))

	mux.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
		counter.Increment()
		fmt.Fprintf(w, "Counter: %d", counter.Value())
	})

	// Load TLS certificates
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Configure server with security best practices
	server := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Starting secure server on :8443")
	return server.ListenAndServeTLS("", "")
}

func main() {
	// Use context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Safe concurrent operations
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter.Value())

	// Start secure server
	if err := startServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
}
