// Package user provides user management functionality.
package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

// Password validation constants.
const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
)

// UserService handles user-related operations with explicit dependencies.
type UserService struct {
	db *Database
}

// NewUserService creates a new UserService with the given database connection.
func NewUserService(db *Database) *UserService {
	return &UserService{db: db}
}

// GetUser returns the user with the given ID or returns an error
// if the user doesn't exist or is inactive.
func (s *UserService) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	user, err := fetchUser(id)
	if err != nil {
		return nil, fmt.Errorf("fetching user %d: %w", id, err)
	}

	if !user.Active {
		return nil, errors.New("user inactive")
	}

	return user, nil
}

// ProcessFile reads and processes a file at the given path.
// It returns an error with context if the file cannot be read.
func ProcessFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", path, err)
	}
	// process data
	_ = data
	return nil
}

// ReadConfig reads and validates configuration from a file.
// It uses defer for proper resource cleanup.
func ReadConfig(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file %s: %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	if err := validateConfig(data); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return data, nil
}

// ValidatePassword checks if the password meets length requirements.
// It returns an error if the password is too short or too long.
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return errors.New("password too short")
	}
	if len(password) > MaxPasswordLength {
		return errors.New("password too long")
	}
	return nil
}

// BuildMessage concatenates string parts with spaces using strings.Builder
// for efficient memory allocation.
func BuildMessage(parts []string) string {
	var builder strings.Builder
	for i, part := range parts {
		builder.WriteString(part)
		if i < len(parts)-1 {
			builder.WriteString(" ")
		}
	}
	return builder.String()
}

// ProcessItems transforms a slice of items into results.
// It preallocates the result slice for better performance.
func ProcessItems(items []Item) []Result {
	results := make([]Result, 0, len(items))
	for _, item := range items {
		results = append(results, process(item))
	}
	return results
}

// FormatID converts an integer ID to its string representation.
// It uses strconv.Itoa for better performance than fmt.Sprintf.
func FormatID(id int) string {
	return strconv.Itoa(id)
}

// WaitForTask pauses execution for the specified duration.
// It accepts time.Duration for type safety and clarity.
func WaitForTask(d time.Duration) {
	time.Sleep(d)
}

// ProcessData loads, transforms, and saves data.
// Variables are declared close to their usage for better readability.
func ProcessData() error {
	data, err := loadData()
	if err != nil {
		return fmt.Errorf("loading data: %w", err)
	}

	result, err := transform(data)
	if err != nil {
		return fmt.Errorf("transforming data: %w", err)
	}

	if err := save(result); err != nil {
		return fmt.Errorf("saving result: %w", err)
	}

	return nil
}

// StartWorker starts a background worker that performs work periodically.
// It uses context for graceful shutdown and proper goroutine lifecycle management.
func StartWorker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				doWork()
			}
		}
	}()
}

// FetchAll retrieves users for all given IDs concurrently.
// It uses errgroup for cleaner error handling and context cancellation.
func FetchAll(ctx context.Context, ids []int) ([]*User, error) {
	g, ctx := errgroup.WithContext(ctx)
	users := make([]*User, len(ids))

	for i, id := range ids {
		i, id := i, id // capture loop variables
		g.Go(func() error {
			user, err := fetchUser(id)
			if err != nil {
				return fmt.Errorf("fetching user %d: %w", id, err)
			}
			users[i] = user
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return users, nil
}

// ProcessItemsInPlace processes items and returns a new slice.
// It copies the input slice to prevent modifying the caller's data.
func ProcessItemsInPlace(items []Item) []Item {
	result := make([]Item, len(items))
	copy(result, items)

	for i := range result {
		result[i].Process()
	}
	return result
}

// GetStringValue safely extracts a string value from an interface.
// It returns an error if the type assertion fails instead of panicking.
func GetStringValue(x interface{}) (string, error) {
	val, ok := x.(string)
	if !ok {
		return "", errors.New("expected string type")
	}
	return val, nil
}

// calculateTotal computes the total price of items after applying a discount.
// The discount should be a value between 0 and 1 (e.g., 0.1 for 10% off).
func calculateTotal(items []Item, discount float64) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price
	}
	return total * (1 - discount)
}

// Placeholder types for compilation.
type User struct {
	Active bool
}
type Database struct{}
type Item struct {
	Price float64
}
type Result struct{}

func fetchUser(id int) (*User, error)       { return nil, nil }
func validateConfig(data []byte) error      { return nil }
func process(item Item) Result              { return Result{} }
func loadData() ([]byte, error)             { return nil, nil }
func transform(data []byte) (string, error) { return "", nil }
func save(result string) error              { return nil }
func doWork()                               {}
func (i *Item) Process()                    {}
