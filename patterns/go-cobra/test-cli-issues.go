// test-cli-issues.go
// Sample CLI code with various Cobra/Viper issues for testing the go-cobra pattern.
// This file intentionally contains anti-patterns and mistakes for review detection.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ISSUE: Complex main.go with business logic and configuration
// Should be minimal with just Execute() call
var (
	cfgFile string
	verbose bool
	port    int
)

// ISSUE: Using global viper singleton instead of dedicated instance
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}

// ISSUE: Business logic mixed with command definition
func processData(data string) error {
	// This should be in internal/app package
	fmt.Printf("Processing: %s\n", data)
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "badcli",
	Short: "A CLI with many issues",
	// ISSUE: Missing Long description
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	// ISSUE: Using Run instead of RunE - errors are swallowed
	Run: func(cmd *cobra.Command, args []string) {
		// ISSUE: Reading directly from flags, bypassing Viper precedence
		p, _ := cmd.Flags().GetInt("port")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		// ISSUE: Ignoring errors
		_ = startServer(p, timeout)
	},
}

var processCmd = &cobra.Command{
	Use: "process",
	// ISSUE: Missing Short description
	// ISSUE: No Args validation
	Run: func(cmd *cobra.Command, args []string) {
		// ISSUE: No error handling
		if len(args) > 0 {
			processData(args[0])
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	// ISSUE: Using Run with error handling - should use RunE
	Run: func(cmd *cobra.Command, args []string) {
		// ISSUE: Hardcoded config path
		cfg, err := loadConfig("/etc/myapp/config.yaml")
		if err != nil {
			// ISSUE: Using fmt.Println for errors instead of returning
			fmt.Println("Error loading config:", err)
			return
		}
		fmt.Printf("Config loaded: %v\n", cfg)
	},
}

// ISSUE: No Args validation, will panic on empty args
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource",
	Run: func(cmd *cobra.Command, args []string) {
		// ISSUE: Will panic if args is empty
		id := args[0]
		fmt.Printf("Deleting: %s\n", id)
	},
}

func init() {
	// ISSUE: Not binding flags to Viper
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")

	// ISSUE: Flag descriptions are too short/unclear
	serveCmd.Flags().IntP("port", "p", 8080, "port")
	serveCmd.Flags().Duration("timeout", 30*time.Second, "timeout")

	// ISSUE: Not using flag groups for related flags
	serveCmd.Flags().String("cert", "", "TLS certificate")
	serveCmd.Flags().String("key", "", "TLS key")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(deleteCmd)
}

// ISSUE: No shell completion command
// ISSUE: No version command

// ISSUE: No error context/wrapping
func startServer(port int, timeout time.Duration) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port")
	}
	fmt.Printf("Starting server on :%d with timeout %s\n", port, timeout)
	return nil
}

// ISSUE: Error not wrapped with context
func loadConfig(path string) (map[string]interface{}, error) {
	// Simulated config loading
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
	return map[string]interface{}{"loaded": true}, nil
}

// ISSUE: No graceful shutdown handling
// ISSUE: No configuration validation
// ISSUE: No secrets management pattern

func main() {
	// ISSUE: Not minimal - should only call Execute()
	initLogging()

	if err := rootCmd.Execute(); err != nil {
		// ISSUE: Using fmt.Println instead of os.Stderr
		fmt.Println(err)
		os.Exit(1)
	}
}

// ISSUE: This should not be in main.go
func initLogging() {
	fmt.Println("Logging initialized")
}
