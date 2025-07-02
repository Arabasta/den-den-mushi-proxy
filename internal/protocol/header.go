package protocol

type Header byte

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

	// ParseError indicates an error in parsing the header
	ParseError Header = 0xff
)
