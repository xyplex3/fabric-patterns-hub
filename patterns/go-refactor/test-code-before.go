package user

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Test file with various Go anti-patterns for refactoring
// This file demonstrates common issues that the go-refactor pattern should fix

var db *Database // global variable anti-pattern

func Init() {
	db = connectDatabase()
}

// Deep nesting anti-pattern
func GetUser(id int) (*User, error) {
	if id > 0 {
		user, err := fetchUser(id)
		if err == nil {
			if user.Active {
				return user, nil
			} else {
				return nil, errors.New("user inactive")
			}
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid id")
	}
}

// Poor error handling - loses context
func ProcessFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.New("failed to read file")
	}
	// process data
	_ = data
	return nil
}

// Resource leak - no defer
func ReadConfig(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	err = validateConfig(data)
	if err != nil {
		file.Close()
		return nil, err
	}

	file.Close()
	return data, nil
}

// Magic numbers
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password too short")
	}
	if len(password) > 128 {
		return errors.New("password too long")
	}
	return nil
}

// Inefficient string concatenation
func BuildMessage(parts []string) string {
	msg := ""
	for _, part := range parts {
		msg += part + " "
	}
	return msg
}

// No preallocation
func ProcessItems(items []Item) []Result {
	var results []Result
	for _, item := range items {
		results = append(results, process(item))
	}
	return results
}

// Using fmt.Sprintf for simple conversion
func FormatID(id int) string {
	return fmt.Sprintf("%d", id)
}

// Using int for duration
func WaitForTask(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// Variables declared too early
func ProcessData() error {
	var data []byte
	var err error
	var result string

	data, err = loadData()
	if err != nil {
		return err
	}

	result, err = transform(data)
	if err != nil {
		return err
	}

	return save(result)
}

// Goroutine without cancellation
func StartWorker() {
	go func() {
		for {
			doWork()
			time.Sleep(time.Second)
		}
	}()
}

// Concurrent fetch with manual sync
func FetchAll(ids []int) ([]*User, error) {
	var wg sync.WaitGroup
	users := make([]*User, len(ids))
	var firstErr error
	var errOnce sync.Once

	for i, id := range ids {
		wg.Add(1)
		go func(i, id int) {
			defer wg.Done()
			user, err := fetchUser(id)
			if err != nil {
				errOnce.Do(func() { firstErr = err })
				return
			}
			users[i] = user
		}(i, id)
	}

	wg.Wait()
	return users, firstErr
}

// Modifying input slice
func ProcessItemsInPlace(items []Item) []Item {
	for i := range items {
		items[i].Process()
	}
	return items
}

// Unchecked type assertion
func GetStringValue(x interface{}) string {
	return x.(string) // panics if not string
}

// Missing documentation
func calculateTotal(items []Item, discount float64) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price
	}
	return total * (1 - discount)
}

// Placeholder types for compilation
type User struct {
	Active bool
}
type Database struct{}
type Item struct {
	Price float64
}
type Result struct{}

func connectDatabase() *Database          { return nil }
func fetchUser(id int) (*User, error)     { return nil, nil }
func validateConfig(data []byte) error    { return nil }
func process(item Item) Result            { return Result{} }
func loadData() ([]byte, error)           { return nil, nil }
func transform(data []byte) (string, error) { return "", nil }
func save(result string) error            { return nil }
func doWork()                             {}
func (i *Item) Process()                  {}
