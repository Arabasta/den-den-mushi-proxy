package protocol

import "fmt"

type Header byte

func (h Header) String() string {
	switch h {
	case Input:
		return "Input"
	case Output:
		return "Output"
	case Error:
		return "Error"
	case BlockedControl:
		return "BlockedControl"
	case BlockedCommand:
		return "BlockedCommand"
	case Warn:
		return "Warn"
	case Broadcast:
		return "Broadcast"
	case PtySessionEvent:
		return "PtySessionEvent"
	case Resize:
		return "Resize"
	case Sudo:
		return "Sudo"
	case ClientClose:
		return "ClientClose"
	case PtyNormalClose:
		return "PtyNormalClose"
	case PtyErrorClose:
		return "PtyErrorClose"
	case ParseError:
		return "ParseError"
	default:
		return fmt.Sprintf("Unknown(0x%02x)", byte(h))
	}
}

const (
	// Input for normal client input from websocket to pty
	Input Header = 0x00

	// Output from the pty to the websocket
	Output Header = 0x01

	// Error for errors from the proxy server
	Error Header = 0x02

	// BlockedControl for messages from the proxy server, e.g., "Arrow keys not allowed"
	BlockedControl Header = 0x03

	// BlockedCommand for messages from the proxy server, e.g., "You are not allowed to run su -"
	BlockedCommand Header = 0x04

	// Warn for warnings from the proxy server, e.g., "10 minutes left on your session"
	Warn Header = 0x05

	// Broadcast for broadcasting messages to all connected clients, e.g., "Server is going down for maintenance"
	Broadcast Header = 0x06

	// PtySessionEvent for stuff like "Kei has joined your session"
	PtySessionEvent Header = 0x07

	Reserved4 Header = 0x08

	Reserved5 Header = 0x09

	// Resize called when the terminal is resized from the client
	Resize Header = 0x10

	// Sudo is called when the client wants to switch users. Required since "su" is blocked.
	Sudo Header = 0x11

	// ClientClose for closing the websocket connection
	ClientClose Header = 0x12

	// PtyNormalClose is used when the pty session ends normally
	PtyNormalClose Header = 0x13

	// PtyErrorClose is used when the pty session ends with an error
	PtyErrorClose Header = 0x14

	// ParseError indicates an error in parsing the header
	ParseError Header = 0xff
)
