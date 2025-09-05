package executor

import (
	"encoding/base64"
	"fmt"
	"os"
	"unicode/utf16"

	"execute_command/utils"
)

// Base64Executor implements CommandExecutor interface for base64 encoded commands
type Base64Executor struct {
	logger         *utils.ModuleLogger
	shellType      ShellType
	defaultCommand string
}

// NewBase64Executor creates a new Base64Executor instance
func NewBase64Executor() CommandExecutor {
	shellType := AutoShell
	// For base64 executor, prefer PowerShell for better base64 support
	if IsWindows() {
		shellType = PowerShellShell
	}
	return &Base64Executor{
		logger:         utils.GetModuleLogger("executor.base64"),
		shellType:      shellType,
		defaultCommand: getDefaultBase64Command(),
	}
}

// NewBase64ExecutorWithShell creates a new Base64Executor instance with specific shell
func NewBase64ExecutorWithShell(shellType ShellType) CommandExecutor {
	return &Base64Executor{
		logger:         utils.GetModuleLogger("executor.base64"),
		shellType:      shellType,
		defaultCommand: getDefaultBase64Command(),
	}
}

// SetShellType sets the shell type for the executor
func (be *Base64Executor) SetShellType(shellType ShellType) {
	be.shellType = shellType
}

// ExecuteCommand executes a base64 encoded command
func (be *Base64Executor) ExecuteCommand(encodedCommand string) error {
	// Use default command if no command provided
	if encodedCommand == "" {
		encodedCommand = be.defaultCommand
		be.logger.Info("Using default base64 command")
	} else {
		be.logger.Info("Executing base64 command")
	}
	be.logger.Debug("Executing base64 command (shell: %s)", be.shellType.String())

	// For PowerShell and Linux shell, we can use native/piped methods
	if be.shellType == PowerShellShell || be.shellType == ShShell {
		return be.executeBase64CommandDirect(encodedCommand)
	}

	// For CMD and auto (on Windows), we need to decode first
	decoded, err := be.DecodeCommand(encodedCommand)
	if err != nil {
		be.logger.Error("Failed to decode base64: %v", err)
		return fmt.Errorf("failed to decode base64: %v", err)
	}

	return be.executeCommand(decoded)
}

// EncodeCommand encodes a command to base64
func (be *Base64Executor) EncodeCommand(command string) string {
	// For PowerShell, we need UTF-16LE encoding
	if be.shellType == PowerShellShell {
		return be.encodeForPowerShell(command)
	}
	// For other shells, use UTF-8 encoding
	return base64.StdEncoding.EncodeToString([]byte(command))
}

// encodeForPowerShell encodes a command to base64 using UTF-16LE for PowerShell
func (be *Base64Executor) encodeForPowerShell(command string) string {
	// Convert string to UTF-16 code points
	utf16CodePoints := utf16.Encode([]rune(command))

	// Convert UTF-16 code points to bytes (little-endian)
	utf16Bytes := make([]byte, len(utf16CodePoints)*2)
	for i, codePoint := range utf16CodePoints {
		utf16Bytes[i*2] = byte(codePoint & 0xFF)          // Low byte
		utf16Bytes[i*2+1] = byte((codePoint >> 8) & 0xFF) // High byte
	}

	// PowerShell expects UTF-16LE encoding for -EncodedCommand
	return base64.StdEncoding.EncodeToString(utf16Bytes)
}

// DecodeCommand decodes a base64 encoded command
func (be *Base64Executor) DecodeCommand(encodedCommand string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedCommand)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}
	return string(decoded), nil
}

// executeCommand is a private method that handles the actual command execution
func (be *Base64Executor) executeCommand(command string) error {
	be.logger.Debug("Executing: %s", command)

	// Get the appropriate shell command
	cmd := GetShellCommand(command, be.shellType)

	// Set output to current process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Execute command
	err := cmd.Run()
	if err != nil {
		be.logger.Error("Command execution failed: %v", err)
		return fmt.Errorf("command execution failed: %v", err)
	}

	be.logger.Info("Command executed successfully")
	return nil
}

// executeBase64CommandDirect is a private method that handles base64 command execution directly
func (be *Base64Executor) executeBase64CommandDirect(encodedCommand string) error {
	be.logger.Debug("Executing base64 directly (shell: %s)", be.shellType.String())

	// Get the appropriate shell command for base64
	cmd := GetShellCommandForBase64(encodedCommand, be.shellType)

	// Set output to current process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Execute command
	err := cmd.Run()
	if err != nil {
		be.logger.Error("Base64 command execution failed: %v", err)
		return fmt.Errorf("base64 command execution failed: %v", err)
	}

	be.logger.Info("Base64 command executed successfully")
	return nil
}

// getDefaultBase64Command returns the default base64 encoded command based on OS
func getDefaultBase64Command() string {
	if IsWindows() {
		// "echo Hello from Windows!" encoded for PowerShell (UTF-16LE)
		// We need to use the same encoding logic as EncodeCommand
		command := "ipconfig"
		utf16CodePoints := utf16.Encode([]rune(command))
		utf16Bytes := make([]byte, len(utf16CodePoints)*2)
		for i, codePoint := range utf16CodePoints {
			utf16Bytes[i*2] = byte(codePoint & 0xFF)
			utf16Bytes[i*2+1] = byte((codePoint >> 8) & 0xFF)
		}
		return base64.StdEncoding.EncodeToString(utf16Bytes)
	}
	// "echo Hello from Linux!" encoded for Linux (UTF-8)
	return base64.StdEncoding.EncodeToString([]byte("ifconfig"))
}
