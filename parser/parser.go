package parser

import (
	"flag"
	"fmt"
	"strings"

	"execute_command/executor"
	"execute_command/utils"
)

// Config holds all parsed configuration
type Config struct {
	LogLevel     utils.LogLevel
	ShellType    executor.ShellType
	ExecutorType executor.ExecutorType
	Help         bool
	Action       string
	Args         []string
}

// ParseConfig parses command line arguments and returns configuration
func ParseConfig() (*Config, error) {
	// Parse command line flags
	var logLevel = flag.String("log-level", "ERROR", "Set logging level (DEBUG, INFO, WARN, ERROR, FATAL)")
	var shell = flag.String("shell", "auto", "Set shell type (auto, cmd, powershell, sh)")
	var executorType = flag.String("executor", "base64", "Set executor type (base64, plain)")
	var help = flag.Bool("help", false, "Show help information")
	flag.Parse()

	// Parse log level
	level := utils.ParseLogLevel(*logLevel)

	// Parse shell type
	shellType := executor.ParseShellType(*shell)

	// Parse executor type
	execType := executor.ParseExecutorType(*executorType)

	// Get remaining arguments after flag parsing
	args := flag.Args()

	// Validate arguments
	if !*help && len(args) < 1 {
		return nil, fmt.Errorf("no action specified")
	}

	var action string
	if len(args) > 0 {
		action = args[0]
	}

	return &Config{
		LogLevel:     level,
		ShellType:    shellType,
		ExecutorType: execType,
		Help:         *help,
		Action:       action,
		Args:         args,
	}, nil
}

// ValidateAction validates the action and its arguments
func (c *Config) ValidateAction() error {
	if c.Help {
		return nil // Help doesn't need validation
	}

	// Validate executor and shell compatibility
	if err := c.ValidateExecutorShellCompatibility(); err != nil {
		return err
	}

	switch c.Action {
	case "execute":
		// execute action can work without command (will use default)
	case "encode":
		if len(c.Args) < 2 {
			return fmt.Errorf("usage: go run main.go encode <command>")
		}
	case "decode":
		if len(c.Args) < 2 {
			return fmt.Errorf("usage: go run main.go decode <base64-encoded-command>")
		}
	case "info":
		// No additional arguments needed
	default:
		return fmt.Errorf("unknown action: %s", c.Action)
	}

	return nil
}

// ValidateExecutorShellCompatibility validates that executor and shell types are compatible
func (c *Config) ValidateExecutorShellCompatibility() error {
	// Base64 executor only works with PowerShell and sh shells
	if c.ExecutorType.String() == "base64" {
		if c.ShellType.String() == "cmd" {
			return fmt.Errorf("base64 executor is not compatible with cmd shell. Use powershell or sh instead")
		}
	}

	// Plain executor works with all shells, but warn about suboptimal combinations
	if c.ExecutorType.String() == "plain" && c.ShellType.String() == "powershell" {
		// This is fine, just informational
	}

	return nil
}

// GetCommand returns the command string from arguments
func (c *Config) GetCommand() string {
	if len(c.Args) < 2 {
		return "" // Empty string will trigger default command
	}
	return strings.Join(c.Args[1:], " ")
}

// PrintUsage prints the usage information
func PrintUsage() {
	fmt.Println("Command Executor")
	fmt.Println("Usage:")
	fmt.Println("  go run main.go [flags] <action> [arguments]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -log-level string    Set logging level (DEBUG, INFO, WARN, ERROR, FATAL) (default \"INFO\")")
	fmt.Println("  -shell string        Set shell type (auto, cmd, powershell, sh) (default \"auto\")")
	fmt.Println("  -executor string     Set executor type (base64, plain) (default \"base64\")")
	fmt.Println("  -help               Show help information")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  execute [command]                 - Execute command using specified executor (uses default if no command)")
	fmt.Println("  encode <command>                  - Encode command to base64")
	fmt.Println("  decode <base64-command>           - Decode base64 command")
	fmt.Println("  info                              - Show system information")
	fmt.Println()
	fmt.Println("Executor Types:")
	fmt.Println("  base64     - Execute base64 encoded command (default)")
	fmt.Println("             - Compatible with: powershell, sh")
	fmt.Println("             - NOT compatible with: cmd")
	fmt.Println("  plain      - Execute plaintext command directly")
	fmt.Println("             - Compatible with: cmd, powershell, sh")
	fmt.Println()
	fmt.Println("Shell Types:")
	fmt.Println("  auto        - Automatically choose based on OS (default)")
	fmt.Println("  cmd         - Windows Command Prompt")
	fmt.Println("  powershell  - Windows PowerShell")
	fmt.Println("  sh          - Linux/Unix Sh")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run main.go -executor plain execute                    # Use default plain command")
	fmt.Println("  go run main.go -executor base64 execute                   # Use default base64 command")
	fmt.Println("  go run main.go -executor plain execute \"echo Hello World\"")
	fmt.Println("  go run main.go -executor base64 execute \"ZWNobyBIZWxsbyBXb3JsZA==\"")
	fmt.Println("  go run main.go encode \"dir\"")
	fmt.Println("  go run main.go decode \"ZGly\"")
	fmt.Println("  go run main.go info")
	fmt.Println("  go run main.go -log-level DEBUG -executor plain execute")
	fmt.Println("  go run main.go -shell powershell -executor plain execute")
	fmt.Println("  go run main.go -shell cmd -executor plain execute")
	fmt.Println()
	fmt.Println("Compatibility Matrix:")
	fmt.Println("  Plain Executor:  cmd, powershell, sh")
	fmt.Println("  Base64 Executor: powershell, sh")
	fmt.Println()
	fmt.Println("Note: Flags must come BEFORE the action, not after the command!")
	fmt.Println("  Correct: go run main.go -log-level DEBUG execute \"whoami\"")
	fmt.Println("  Wrong:   go run main.go execute \"whoami\" --log-level DEBUG")
}
