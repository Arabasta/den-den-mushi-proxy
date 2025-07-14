package constants

import "bytes"

func isPaste(b []byte) bool {
	return bytes.HasPrefix(b, PasteStart) && bytes.HasSuffix(b, PasteEnd)
}
