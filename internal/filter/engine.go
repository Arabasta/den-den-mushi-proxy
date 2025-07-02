package filter

import (
	"fmt"
	"strings"
)

type CommandFilter struct {
	buf []byte
}

func (f *CommandFilter) Feed(str string) (string, bool) {
	return "", true
}

func isNaiveFilterOk(inputStr string) bool {
	fmt.Println("User entered:", inputStr)

	if strings.Contains(strings.TrimSpace(inputStr), "su -") {
		return false
	} else if strings.HasPrefix(strings.TrimSpace(inputStr), "exec") {
		return false
	}
	return true
}

//// isBannedControl is a naive implementation to filter arrows and other control characters
//func isBannedControl(data []byte) (bool, string) {
//	switch {
//	case bytes.Equal(data, ArrowUp):
//		return true, "Arrow Up"
//	case bytes.Equal(data, ArrowDown):
//		return true, "Arrow Down"
//	case bytes.Equal(data, ArrowRight):
//		return true, "Arrow Right"
//	case bytes.Equal(data, ArrowLeft):
//		return true, "Arrow Left"
//	case bytes.Equal(data, CtrlR):
//		return true, "Ctrl+R"
//	case bytes.HasPrefix(data, PasteStart) && bytes.HasSuffix(data, PasteEnd):
//		return true, "Paste"
//	case bytes.Equal(data, SemiColon):
//		return true, ";"
//	case bytes.Equal(data, Ampersand):
//		return true, "&"
//	case bytes.Equal(data, Pipe):
//		return true, "|"
//	case bytes.Equal(data, LeftParenthesis):
//		return true, "("
//	case bytes.Equal(data, RightParenthesis):
//		return true, ")"
//	default:
//		return false, ""
//	}
//}
