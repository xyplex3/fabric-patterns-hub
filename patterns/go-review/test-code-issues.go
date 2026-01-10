package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Test file with various Go issues for code review
// This file demonstrates common problems that the go-review pattern should identify

// Global variable - should be dependency injected
var db *sql.DB

// Missing documentation for exported type
type UserService struct {
	cache map[int]*User
}

// Missing documentation for exported function
func NewUserService() *UserService {
	return &UserService{
		cache: make(map[int]*User),
	}
}

// Deep nesting - should use early returns
func (s *UserService) GetUser(id int) (*User, error) {
	if id > 0 {
		user, err := s.fetchUser(id)
		if err == nil {
			if user != nil {
				if user.Active {
					return user, nil
				} else {
					return nil, errors.New("user is inactive")
				}
			} else {
				return nil, errors.New("user not found")
			}
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid user id")
	}
}

// Error handling issues - loses context
func (s *UserService) fetchUser(id int) (*User, error) {
	// SQL injection vulnerability!
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", id)
	row := db.QueryRow(query)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Active)
	if err != nil {
		return nil, errors.New("database error") // loses original error
	}
	return &user, nil
}

// Ignored error - critical issue
func (s *UserService) UpdateCache(user *User) {
	_ = s.validateUser(user) // ignored error!
	s.cache[user.ID] = user
}

// Goroutine without cancellation - resource leak
func (s *UserService) StartBackgroundSync() {
	go func() {
		for {
			s.syncCache()
			time.Sleep(time.Minute)
		}
	}()
}

// Magic numbers
func (s *UserService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("too short")
	}
	if len(password) > 72 {
		return errors.New("too long")
	}
	return nil
}

// Unchecked type assertion - can panic
func processValue(v interface{}) string {
	return v.(string) // panics if not string!
}

// Modifying input slice - side effects
func processUsers(users []*User) {
	for i := range users {
		users[i].Name = "processed_" + users[i].Name
	}
}

// Using int for duration - type confusion
func waitForRetry(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// Inefficient string building
func buildQuery(fields []string) string {
	query := "SELECT "
	for i, f := range fields {
		query += f
		if i < len(fields)-1 {
			query += ", "
		}
	}
	return query + " FROM users"
}

// Missing error handling in defer
func (s *UserService) SaveToFile(path string) error {
	// Would need: file, err := os.Create(path)
	// defer file.Close() // error from Close is ignored
	return nil
}

// Getters with "Get" prefix - not idiomatic
func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetID() int {
	return u.ID
}

// Placeholder types
type User struct {
	ID     int
	Name   string
	Active bool
}

func (s *UserService) validateUser(u *User) error {
	if u == nil {
		return errors.New("nil user")
	}
	return nil
}

func (s *UserService) syncCache() {}
