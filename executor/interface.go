package executor

// CommandExecutor defines the interface for command execution operations
type CommandExecutor interface {
	// ExecuteCommand executes a command (behavior depends on executor type)
	ExecuteCommand(command string) error

	// EncodeCommand encodes a command to base64
	EncodeCommand(command string) string

	// DecodeCommand decodes a base64 encoded command
	DecodeCommand(encodedCommand string) (string, error)
}
