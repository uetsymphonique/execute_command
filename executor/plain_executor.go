package executor

import (
	"fmt"
	"os"

	"execute_command/utils"
)

// PlainExecutor implements CommandExecutor interface for plaintext commands
type PlainExecutor struct {
	logger         *utils.ModuleLogger
	shellType      ShellType
	defaultCommand string
}

// NewPlainExecutor creates a new PlainExecutor instance
func NewPlainExecutor() CommandExecutor {
	return &PlainExecutor{
		logger:         utils.GetModuleLogger("executor.plain"),
		shellType:      AutoShell,
		defaultCommand: getDefaultPlainCommand(),
	}
}

// NewPlainExecutorWithShell creates a new PlainExecutor instance with specific shell
func NewPlainExecutorWithShell(shellType ShellType) CommandExecutor {
	return &PlainExecutor{
		logger:         utils.GetModuleLogger("executor.plain"),
		shellType:      shellType,
		defaultCommand: getDefaultPlainCommand(),
	}
}

// SetShellType sets the shell type for the executor
func (pe *PlainExecutor) SetShellType(shellType ShellType) {
	pe.shellType = shellType
}

// ExecuteCommand executes a plaintext command directly
func (pe *PlainExecutor) ExecuteCommand(command string) error {
	// Use default command if no command provided
	if command == "" {
		command = pe.defaultCommand
		pe.logger.Info("Using default plaintext command: %s", command)
	} else {
		pe.logger.Info("Executing plaintext command: %s", command)
	}
	return pe.executeCommand(command)
}

// EncodeCommand returns the command as-is (no encoding for plain executor)
func (pe *PlainExecutor) EncodeCommand(command string) string {
	pe.logger.Debug("Plain executor - no encoding needed")
	return command
}

// DecodeCommand returns the command as-is (no decoding for plain executor)
func (pe *PlainExecutor) DecodeCommand(command string) (string, error) {
	pe.logger.Debug("Plain executor - no decoding needed")
	return command, nil
}

// executeCommand executes the command using the appropriate shell
func (pe *PlainExecutor) executeCommand(command string) error {
	pe.logger.Debug("Executing: %s (shell: %s)", command, pe.shellType.String())

	cmd := GetShellCommand(command, pe.shellType)

	// Set output to current process stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Execute the command
	err := cmd.Run()
	if err != nil {
		pe.logger.Error("Command execution failed: %v", err)
		return fmt.Errorf("command execution failed: %v", err)
	}

	pe.logger.Info("Command executed successfully")
	return nil
}

// getDefaultPlainCommand returns the default plaintext command based on OS
func getDefaultPlainCommand() string {
	if IsWindows() {
		return "ipconfig"
	}
	return "ifconfig"
}
