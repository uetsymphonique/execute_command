# Command Executor

A Golang program that can execute commands using different executor types (plain and base64) on both Windows and Linux. Designed with modular architecture, comprehensive logging, and smart shell compatibility.

## Features

- **Two Executor Types**: Plain (direct execution) and Base64 (encoded execution)
- **Smart Shell Selection**: Automatic shell selection based on executor type
- **Default Commands**: Built-in default commands for both executor types
- **Cross-platform**: Works on both Windows and Linux
- **Modular Architecture**: Separated logic into distinct modules (executor, parser, utils)
- **Comprehensive Logging**: Multi-level logging with module names and colored output
- **Shell Compatibility**: Enforced compatibility between executor and shell types
- **Command Line Flags**: Flexible configuration through command line arguments

## Installation

1. Ensure you have Go installed (version 1.21 or higher)
2. Clone or download this project
3. Run the command:

```bash
go mod tidy
```

## Usage

### Basic Usage

```bash
# Show help
go run main.go -help

# Execute with default commands
go run main.go -executor plain execute                    # Use default plain command
go run main.go -executor base64 execute                   # Use default base64 command

# Execute with custom commands
go run main.go -executor plain execute "whoami"
go run main.go -executor base64 execute "d2hvYW1p"        # Base64 encoded "whoami"

# Encode/Decode commands
go run main.go encode "whoami"                            # Encode command
go run main.go decode "d2hvYW1p"                          # Decode command

# Show system information
go run main.go info
```

### Command Line Flags

| Flag         | Description                                         | Example                                             |
| ------------ | --------------------------------------------------- | --------------------------------------------------- |
| `-log-level` | Set logging level (DEBUG, INFO, WARN, ERROR, FATAL) | `go run main.go -log-level DEBUG execute "whoami"`  |
| `-shell`     | Set shell type (auto, cmd, powershell, sh)          | `go run main.go -shell powershell execute "whoami"` |
| `-executor`  | Set executor type (base64, plain)                   | `go run main.go -executor plain execute "whoami"`   |
| `-help`      | Show help information                               | `go run main.go -help`                              |

### Executor Types

| Executor Type | Description                         | Compatible Shells   | Default Command                                       |
| ------------- | ----------------------------------- | ------------------- | ----------------------------------------------------- |
| `plain`       | Execute plaintext commands directly | cmd, powershell, sh | `echo Hello from Windows!` / `echo Hello from Linux!` |
| `base64`      | Execute base64 encoded commands     | powershell, sh      | Base64 encoded version of default commands            |

### Shell Types

| Shell Type   | Description                                | Platform   | Base64 Support     |
| ------------ | ------------------------------------------ | ---------- | ------------------ |
| `auto`       | Automatically choose based on OS (default) | All        | Platform dependent |
| `cmd`        | Windows Command Prompt                     | Windows    | Not supported      |
| `powershell` | Windows PowerShell                         | Windows    | Native support     |
| `sh`         | Linux/Unix Sh                              | Linux/Unix | Pipe method        |

### Compatibility Matrix

| Executor Type | cmd | powershell | sh  |
| ------------- | --- | ---------- | --- |
| Plain         | ✓   | ✓          | ✓   |
| Base64        | ✗   | ✓          | ✓   |

## Real-world Examples

### Windows Examples

```bash
# Plain executor with CMD
go run main.go -executor plain -shell cmd execute "dir /b"
go run main.go -executor plain -shell cmd execute                    # Use default

# Plain executor with PowerShell
go run main.go -executor plain -shell powershell execute "Get-Process"
go run main.go -executor plain -shell powershell execute             # Use default

# Base64 executor with PowerShell (CMD not supported)
go run main.go -executor base64 -shell powershell execute "d2hvYW1p"  # Base64 "whoami"
go run main.go -executor base64 -shell powershell execute             # Use default

# Encode commands for base64 execution
go run main.go encode "whoami"                                       # Get base64 string
go run main.go -executor base64 execute "d2hvYW1p"                   # Execute base64
```

### Linux Examples

```bash
# Plain executor with sh
go run main.go -executor plain -shell sh execute "ls -la"
go run main.go -executor plain -shell sh execute                     # Use default

# Base64 executor with sh
go run main.go -executor base64 -shell sh execute "bHMgLWxh"         # Base64 "ls -la"
go run main.go -executor base64 -shell sh execute                     # Use default
```

## File Structure

```
execute_command/
├── main.go                    # Main entry point with CLI interface
├── go.mod                     # Go module file
├── parser/                    # Command line parsing module
│   └── parser.go             # Argument parsing and validation
├── executor/                  # Executor module
│   ├── interface.go          # CommandExecutor interface
│   ├── base64_executor.go    # Base64Executor implementation
│   ├── plain_executor.go     # PlainExecutor implementation
│   └── executor.go           # Factory and utility functions
└── utils/                     # Utilities module
    └── logger.go             # Logging utilities with module names
```

### Module Architecture

- **`parser/parser.go`**: Handles command line argument parsing and validation
- **`executor/interface.go`**: Defines the common `CommandExecutor` interface
- **`executor/base64_executor.go`**: Base64 executor with UTF-16LE encoding for PowerShell
- **`executor/plain_executor.go`**: Plain text executor for direct command execution
- **`executor/executor.go`**: Factory pattern and utility functions
- **`utils/logger.go`**: Comprehensive logging system with module names and colored output
- **`main.go`**: CLI interface using the parser and executor modules

## Default Commands

The program includes built-in default commands for both executor types:

### Plain Executor Defaults

- **Windows**: `echo Hello from Windows!`
- **Linux**: `echo Hello from Linux!`

### Base64 Executor Defaults

- **Windows**: Base64 encoded `echo Hello from Windows!` (UTF-16LE for PowerShell)
- **Linux**: Base64 encoded `echo Hello from Linux!` (UTF-8)

## Base64 Command Execution Methods

### PowerShell (Windows)

- Uses native `-EncodedCommand` parameter
- Commands are encoded using UTF-16LE (Little Endian)
- Example: `powershell -EncodedCommand "base64string"`

### Linux Shell (sh)

- Uses pipe method: `echo "base64string" | base64 -d | sh`
- Commands are encoded using UTF-8
- Example: `echo "base64string" | base64 -d | sh`

### CMD (Windows)

- **Not supported** for base64 executor
- Would require decoding first, which defeats the purpose

## Security Notice

⚠️ **Warning**: This program can execute any command. Only use with trusted commands and in a secure environment.

## Build and Run

### Build for Different Platforms

The program can be built for various operating systems and architectures using Go's cross-compilation features:

#### Windows (64-bit)

```bash
GOOS=windows GOARCH=amd64 go build -o execute_command.exe main.go
```

#### Windows (32-bit)

```bash
GOOS=windows GOARCH=386 go build -o execute_command.exe main.go
```

#### Linux (64-bit)

```bash
GOOS=linux GOARCH=amd64 go build -o execute_command main.go
```

#### Linux (32-bit)

```bash
GOOS=linux GOARCH=386 go build -o execute_command main.go
```

#### Linux (ARM64)

```bash
GOOS=linux GOARCH=arm64 go build -o execute_command main.go
```

#### macOS (64-bit)

```bash
GOOS=darwin GOARCH=amd64 go build -o execute_command main.go
```

#### macOS (Apple Silicon)

```bash
GOOS=darwin GOARCH=arm64 go build -o execute_command main.go
```

#### FreeBSD (64-bit)

```bash
GOOS=freebsd GOARCH=amd64 go build -o execute_command main.go
```

#### OpenBSD (64-bit)

```bash
GOOS=openbsd GOARCH=amd64 go build -o execute_command main.go
```

### Build Scripts

You can create build scripts for easier compilation:

#### build.bat (Windows)

```batch
@echo off
echo Building for Windows 64-bit...
GOOS=windows GOARCH=amd64 go build -o execute_command.exe main.go
echo Building for Linux 64-bit...
GOOS=linux GOARCH=amd64 go build -o execute_command_linux main.go
echo Building for macOS 64-bit...
GOOS=darwin GOARCH=amd64 go build -o execute_command_macos main.go
echo Build complete!
```

#### build.sh (Linux/macOS)

```bash
#!/bin/bash
echo "Building for multiple platforms..."

echo "Building for Windows 64-bit..."
GOOS=windows GOARCH=amd64 go build -o execute_command.exe main.go

echo "Building for Linux 64-bit..."
GOOS=linux GOARCH=amd64 go build -o execute_command_linux main.go

echo "Building for macOS 64-bit..."
GOOS=darwin GOARCH=amd64 go build -o execute_command_macos main.go

echo "Building for macOS Apple Silicon..."
GOOS=darwin GOARCH=arm64 go build -o execute_command_macos_arm64 main.go

echo "Build complete!"
```

### Run the Built Files

#### Windows

```bash
# Run the Windows executable
./execute_command.exe -executor plain execute "whoami"
./execute_command.exe -executor base64 execute "d2hvYW1p"
```

#### Linux

```bash
# Make executable and run
chmod +x execute_command
./execute_command -executor plain execute "whoami"
./execute_command -executor base64 execute "d2hvYW1p"
```

#### macOS

```bash
# Make executable and run
chmod +x execute_command
./execute_command -executor plain execute "whoami"
./execute_command -executor base64 execute "d2hvYW1p"
```

### Supported Platforms

| OS      | Architecture  | GOOS    | GOARCH | Status       |
| ------- | ------------- | ------- | ------ | ------------ |
| Windows | 64-bit        | windows | amd64  | ✅ Supported |
| Windows | 32-bit        | windows | 386    | ✅ Supported |
| Linux   | 64-bit        | linux   | amd64  | ✅ Supported |
| Linux   | 32-bit        | linux   | 386    | ✅ Supported |
| Linux   | ARM64         | linux   | arm64  | ✅ Supported |
| macOS   | 64-bit        | darwin  | amd64  | ✅ Supported |
| macOS   | Apple Silicon | darwin  | arm64  | ✅ Supported |
| FreeBSD | 64-bit        | freebsd | amd64  | ✅ Supported |
| OpenBSD | 64-bit        | openbsd | amd64  | ✅ Supported |

### Build Notes

- **Cross-compilation**: Go supports cross-compilation out of the box
- **CGO**: If you use CGO, you may need to set `CGO_ENABLED=0` for cross-compilation
- **File extensions**: Windows executables should have `.exe` extension
- **Permissions**: Linux/macOS executables need execute permissions (`chmod +x`)
- **Dependencies**: All dependencies are statically linked, so no external libraries are required

## Available Commands

| Command             | Description                              | Example                                           |
| ------------------- | ---------------------------------------- | ------------------------------------------------- |
| `execute [command]` | Execute command using specified executor | `go run main.go -executor plain execute "whoami"` |
| `encode <command>`  | Encode command to base64                 | `go run main.go encode "whoami"`                  |
| `decode <base64>`   | Decode base64 to command                 | `go run main.go decode "d2hvYW1p"`                |
| `info`              | Show system information                  | `go run main.go info`                             |
