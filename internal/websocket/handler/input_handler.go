package handler

import (
	"bytes"
	"den-den-mushi-Go/internal/websocket/protocol"
	"den-den-mushi-Go/pkg/token"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"strings"
)

type InputHandler struct {
	buf []byte
}

func (h *InputHandler) Handle(pkt protocol.Packet, pty io.Writer, _ *websocket.Conn, claims *token.Claims) error {
	fmt.Println(pkt.Data)

	// allow all
	if claims.Connection.Purpose == "change" {
		_, err := pty.Write(pkt.Data)
		if err != nil {
			return err
		}
	} else if claims.Connection.Purpose == "health" { // dummy implementation
		if isBannedControl(pkt.Data) {
			fmt.Println("Blocked input:", pkt.Data)
		} else {
			if bytes.Equal(pkt.Data, Enter) {
				if h.isNaiveFilterOk() {
					_, err := pty.Write(pkt.Data)
					if err != nil {
						return err
					}
					h.buf = nil // clear after enter
				} else {
					fmt.Println("Blocked input:", string(h.buf))
					_, _ = pty.Write([]byte{3}) // send Ctrl-C
					h.buf = nil                 // clear even if blocked
				}
			} else if bytes.Equal(pkt.Data, Backspace) {
				if len(h.buf) > 0 {
					h.buf = h.buf[:len(h.buf)-1]
				}
				_, err := pty.Write(pkt.Data)
				if err != nil {
					return err
				}
			} else { // is Enter, check buf
				h.buf = append(h.buf, pkt.Data...)
				_, err := pty.Write(pkt.Data)
				if err != nil {
					return err
				}
			}

		}

	}

	return nil
}

func (h *InputHandler) isNaiveFilterOk() bool {
	inputStr := string(h.buf)
	fmt.Println("User entered:", inputStr)

	if strings.TrimSpace(inputStr) == "su -" {
		return false
	}
	return true
}

// isBannedControl is a naive implementation to filter arrows and other control characters
func isBannedControl(data []byte) bool {
	switch {
	case bytes.Equal(data, ArrowUp):
		return true
	case bytes.Equal(data, ArrowDown):
		return true
	case bytes.Equal(data, ArrowLeft):
		return true
	case bytes.Equal(data, ArrowRight):
		return true
	case bytes.Equal(data, CtrlR):
		return true
	default:
		return false
	}
}

var (
	ArrowUp    = []byte{27, 91, 65}
	ArrowDown  = []byte{27, 91, 66}
	ArrowRight = []byte{27, 91, 67}
	ArrowLeft  = []byte{27, 91, 68}
	CtrlR      = []byte{18}

	Enter     = []byte{13}
	Backspace = []byte{127}
)

//var (
//	ArrowUp    = []byte{27, 91, 65}
//	ArrowDown  = []byte{27, 91, 66}
//	ArrowRight = []byte{27, 91, 67}
//	ArrowLeft  = []byte{27, 91, 68}
//
//	Enter     = []byte{13}
//	Backspace = []byte{127}
//
//	CtrlA            = []byte{1}
//	CtrlB            = []byte{2}
//	CtrlC            = []byte{3}
//	CtrlD            = []byte{4}
//	CtrlE            = []byte{5}
//	CtrlF            = []byte{6}
//	CtrlG            = []byte{7}
//	CtrlH            = []byte{8}
//	CtrlI            = []byte{9}
//	CtrlJ            = []byte{10}
//	CtrlK            = []byte{11}
//	CtrlL            = []byte{12}
//	CtrlM            = []byte{13}
//	CtrlN            = []byte{14}
//	CtrlO            = []byte{15}
//	CtrlP            = []byte{16}
//	CtrlQ            = []byte{17}
//	CtrlR            = []byte{18}
//	CtrlS            = []byte{19}
//	CtrlT            = []byte{20}
//	CtrlU            = []byte{21}
//	CtrlV            = []byte{22}
//	CtrlW            = []byte{23}
//	CtrlX            = []byte{24}
//	CtrlY            = []byte{25}
//	CtrlZ            = []byte{26}
//	CtrlBackslash    = []byte{28}
//	CtrlRightBracket = []byte{29}
//	CtrlCaret        = []byte{30}
//	CtrlUnderscore   = []byte{31}
//
//	Space           = []byte{32}
//	ExclamationMark = []byte{33} // !
//
//	PasteStart = []byte{27, 91, 50, 48, 48, 126}
//	PasteEnd   = []byte{27, 91, 50, 48, 49, 126}
//)
//
//func isArrowKey(data []byte) []byte {
//	switch {
//	case bytes.Equal(data, ArrowUp):
//		return nil
//	case bytes.Equal(data, ArrowDown):
//		return nil
//	case bytes.Equal(data, ArrowLeft):
//		return nil
//	case bytes.Equal(data, ArrowRight):
//		return nil
//	default:
//		return data
//	}
//}
