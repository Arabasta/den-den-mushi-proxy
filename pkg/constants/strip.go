package constants

func StripPaste(b []byte) []byte {
	if isPaste(b) {
		return b[len(PasteStart) : len(b)-len(PasteEnd)]
	}
	return b
}
