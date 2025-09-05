package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"execute_command/utils"
)

// ExecutorType represents different types of executors
type ExecutorType int

const (
	Base64Type ExecutorType = iota
	PlainType
	// Future executor types can be added here
	// EncryptedType
)

// String returns the string representation of ExecutorType
func (et ExecutorType) String() string {
	switch et {
	case Base64Type:
		return "base64"
	case PlainType:
		return "plain"
	default:
		return "base64"
	}
}

// ExecutorFactory creates executors based on type
type ExecutorFactory struct {
	logger *utils.ModuleLogger
}

// NewExecutorFactory creates a new ExecutorFactory
func NewExecutorFactory() *ExecutorFactory {
	return &ExecutorFactory{
		logger: utils.GetModuleLogger("executor.factory"),
	}
}

// CreateExecutor creates an executor based on the specified type
func (ef *ExecutorFactory) CreateExecutor(executorType ExecutorType) CommandExecutor {
	ef.logger.Debug("Creating %s executor", executorType.String())
	switch executorType {
	case Base64Type:
		return NewBase64Executor()
	case PlainType:
		return NewPlainExecutor()
	default:
		return NewBase64Executor() // Default to base64
	}
}

// CreateExecutorWithShell creates an executor with specific shell type
func (ef *ExecutorFactory) CreateExecutorWithShell(executorType ExecutorType, shellType ShellType) CommandExecutor {
	ef.logger.Debug("Creating %s executor (shell: %s)", executorType.String(), shellType.String())
	switch executorType {
	case Base64Type:
		return NewBase64ExecutorWithShell(shellType)
	case PlainType:
		return NewPlainExecutorWithShell(shellType)
	default:
		return NewBase64ExecutorWithShell(shellType) // Default to base64
	}
}

// GetDefaultExecutor creates a default executor (Base64)
func (ef *ExecutorFactory) GetDefaultExecutor() CommandExecutor {
	return ef.CreateExecutor(Base64Type)
}

// GetDefaultExecutorWithShell creates a default executor with specific shell
func (ef *ExecutorFactory) GetDefaultExecutorWithShell(shellType ShellType) CommandExecutor {
	return ef.CreateExecutorWithShell(Base64Type, shellType)
}

// SystemInfo provides information about the current system
type SystemInfo struct {
	OS   string
	Arch string
}

// GetSystemInfo returns current system information
func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// IsWindows checks if the current OS is Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux checks if the current OS is Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsUnix checks if the current OS is Unix-like
func IsUnix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd"
}

// ShellType represents the type of shell to use
type ShellType int

const (
	AutoShell       ShellType = iota // Automatically choose based on OS
	CMDShell                         // Windows CMD
	PowerShellShell                  // Windows PowerShell
	ShShell                          // Linux/Unix Sh
)

// String returns the string representation of ShellType
func (st ShellType) String() string {
	switch st {
	case AutoShell:
		return "auto"
	case CMDShell:
		return "cmd"
	case PowerShellShell:
		return "powershell"
	case ShShell:
		return "sh"
	default:
		return "auto"
	}
}

// ParseShellType parses a string to ShellType
func ParseShellType(shell string) ShellType {
	switch strings.ToLower(shell) {
	case "cmd":
		return CMDShell
	case "powershell", "ps", "ps1":
		return PowerShellShell
	case "sh":
		return ShShell
	case "auto":
		return AutoShell
	default:
		return AutoShell
	}
}

// ParseExecutorType parses a string to ExecutorType
func ParseExecutorType(executor string) ExecutorType {
	switch strings.ToLower(executor) {
	case "base64":
		return Base64Type
	case "plain":
		return PlainType
	default:
		return Base64Type // Default to base64
	}
}

// GetShellCommand returns the appropriate shell command for the current OS and shell type
func GetShellCommand(command string, shellType ShellType) *exec.Cmd {
	// If auto, determine based on OS
	if shellType == AutoShell {
		if IsWindows() {
			shellType = CMDShell
		} else {
			shellType = ShShell
		}
	}

	switch shellType {
	case CMDShell:
		return exec.Command("cmd", "/C", command)
	case PowerShellShell:
		return exec.Command("powershell", "-Command", command)
	case ShShell:
		return exec.Command("sh", "-c", command)
	default:
		// Fallback to auto behavior
		if IsWindows() {
			return exec.Command("cmd", "/C", command)
		} else {
			return exec.Command("sh", "-c", command)
		}
	}
}

// GetShellCommandForBase64 returns the appropriate shell command for base64 encoded commands
func GetShellCommandForBase64(encodedCommand string, shellType ShellType) *exec.Cmd {
	// If auto, determine based on OS
	if shellType == AutoShell {
		if IsWindows() {
			shellType = CMDShell
		} else {
			shellType = ShShell
		}
	}

	switch shellType {
	case CMDShell:
		// For CMD, we need to decode the base64 first since CMD can't handle encoded commands
		// This will be handled by the caller
		return exec.Command("cmd", "/C", encodedCommand)
	case PowerShellShell:
		// PowerShell can handle base64 encoded commands directly with -EncodedCommand
		return exec.Command("powershell", "-EncodedCommand", encodedCommand)
	case ShShell:
		// Sh can handle base64 through pipe: echo "base64" | base64 -d | sh
		command := fmt.Sprintf("echo '%s' | base64 -d | sh", encodedCommand)
		return exec.Command("sh", "-c", command)
	default:
		// Fallback to auto behavior
		if IsWindows() {
			return exec.Command("cmd", "/C", encodedCommand)
		} else {
			// For Linux, use pipe method
			command := fmt.Sprintf("echo '%s' | base64 -d | sh", encodedCommand)
			return exec.Command("sh", "-c", command)
		}
	}
}

// GetDefaultShellCommand returns the default shell command (backward compatibility)
func GetDefaultShellCommand(command string) *exec.Cmd {
	return GetShellCommand(command, AutoShell)
}

// ValidateCommand checks if a command is safe to execute (basic validation)
func ValidateCommand(command string) error {
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Add more validation rules here if needed
	// For example, check for dangerous commands, etc.

	return nil
}

// ExecuteCommandWithValidation executes a command with basic validation
func ExecuteCommandWithValidation(command string, shellType ShellType) error {
	if err := ValidateCommand(command); err != nil {
		return err
	}

	cmd := GetShellCommand(command, shellType)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Executing command: %s\n", command)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}

	return nil
}

// ExecuteCommandWithValidationDefault executes a command with default shell (backward compatibility)
func ExecuteCommandWithValidationDefault(command string) error {
	return ExecuteCommandWithValidation(command, AutoShell)
}
