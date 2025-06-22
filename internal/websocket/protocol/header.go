package protocol

type Header byte

const (
	// Input for normal user input
	Input Header = 0x00

	// Resize called when the terminal is resized
	Resize Header = 0x01

	// ParseError indicates an error in parsing the header
	ParseError Header = 0xff
)
