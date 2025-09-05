package main

import (
	"fmt"
	"os"

	"execute_command/executor"
	"execute_command/parser"
	"execute_command/utils"
)

func main() {
	// Parse command line configuration
	config, err := parser.ParseConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		parser.PrintUsage()
		os.Exit(1)
	}

	// Initialize global logger
	utils.InitGlobalLogger(config.LogLevel)
	logger := utils.GetModuleLogger("main")

	// Show help if requested
	if config.Help {
		parser.PrintUsage()
		return
	}

	// Validate action
	if err := config.ValidateAction(); err != nil {
		logger.Error("%v", err)
		parser.PrintUsage()
		os.Exit(1)
	}

	logger.Info("Starting Command Executor")

	// Create executor factory and get executor with specified type and shell
	factory := executor.NewExecutorFactory()

	// For base64 executor with auto shell on Windows, prefer PowerShell
	shellType := config.ShellType
	if config.ExecutorType == executor.Base64Type && config.ShellType == executor.AutoShell && executor.IsWindows() {
		shellType = executor.PowerShellShell
	}

	cmdExecutor := factory.CreateExecutorWithShell(config.ExecutorType, shellType)

	// Display system information
	sysInfo := executor.GetSystemInfo()
	logger.Info("Running on %s/%s", sysInfo.OS, sysInfo.Arch)
	logger.Info("Using shell: %s", shellType.String())
	logger.Info("Using executor: %s", config.ExecutorType.String())

	// Execute action
	action := config.Action

	switch action {
	case "execute":
		command := config.GetCommand()
		logger.Debug("Executing command: %s", command)
		err := cmdExecutor.ExecuteCommand(command)
		if err != nil {
			logger.Error("Error executing command: %v", err)
			os.Exit(1)
		}

	case "encode":
		command := config.GetCommand()
		logger.Debug("Encoding command: %s", command)
		encoded := cmdExecutor.EncodeCommand(command)
		fmt.Printf("Base64 encoded command: %s\n", encoded)
		logger.Info("Command encoded successfully")

	case "decode":
		encoded := config.GetCommand()
		logger.Debug("Decoding base64 command, length: %d", len(encoded))
		decoded, err := cmdExecutor.DecodeCommand(encoded)
		if err != nil {
			logger.Error("Error decoding: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Decoded command: %s\n", decoded)
		logger.Info("Command decoded successfully")

	case "info":
		logger.Info("Displaying system information")
		printSystemInfo()

	default:
		logger.Warn("Unknown action: %s", action)
		parser.PrintUsage()
	}
}

func printSystemInfo() {
	sysInfo := executor.GetSystemInfo()
	fmt.Printf("System Information:\n")
	fmt.Printf("  OS: %s\n", sysInfo.OS)
	fmt.Printf("  Architecture: %s\n", sysInfo.Arch)
	fmt.Printf("  Is Windows: %t\n", executor.IsWindows())
	fmt.Printf("  Is Linux: %t\n", executor.IsLinux())
	fmt.Printf("  Is Unix-like: %t\n", executor.IsUnix())
}
