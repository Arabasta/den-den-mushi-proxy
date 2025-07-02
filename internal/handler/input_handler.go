package handler

import (
	"den-den-mushi-Go/internal/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

type InputHandler struct{}

func (h *InputHandler) Handle(pkt protocol.Packet, pty io.Writer, ws *websocket.Conn, claims *token.Claims) (string, error) {
	//if claims.Connection.Purpose == dto.Change {
	_, err := pty.Write(pkt.Data)
	if err != nil {
		return "", err
	}
	return "", err
	//} else if claims.Connection.Purpose == dto.Healthcheck { // dummy implementation
	//	if banned, reason := isBannedControl(pkt.Data); banned {
	//		fmt.Println("Blocked input:", reason)
	//
	//		ws.WriteMessage(websocket.BinaryMessage, protocol.PacketToByte(protocol.Packet{Header: protocol.BlockedControl, Data: []byte(fmt.Sprintf("\r\n[BLOCKED CONTROL CHAR] %s", reason))}))
	//
	//		return fmt.Sprintf("\n[BLOCKED CONTROL CHAR DETECTED] %s", reason), nil
	//	} else {
	//		if bytes.Equal(pkt.Data, Enter) {
	//			if h.isNaiveFilterOk() {
	//				// allowed, write to pty
	//				_, err := pty.Write(pkt.Data)
	//				if err != nil {
	//					return "", err
	//				}
	//				h.buf = nil // clear after enter
	//			} else {
	//				// blocked, send Ctrl-C to pty
	//				fmt.Println("Blocked input:", string(h.buf))
	//				ws.WriteMessage(websocket.BinaryMessage, protocol.PacketToByte(protocol.Packet{Header: protocol.BlockedCommand, Data: []byte(fmt.Sprintf("[BLACKLISTED COMMAND] %s", string(h.buf)))}))
	//				_, _ = pty.Write([]byte{3}) // send Ctrl-C
	//				h.buf = nil                 // clear even if blocked
	//				return fmt.Sprintf("\n[BLOCKED COMMAND DETECTED, SENDING CTRL+C] %s", string(h.buf)), nil
	//			}
	//		} else if bytes.Equal(pkt.Data, Backspace) { // handle Backspace
	//			if len(h.buf) > 0 {
	//				h.buf = h.buf[:len(h.buf)-1]
	//			}
	//			_, err := pty.Write(pkt.Data)
	//			if err != nil {
	//				return "", err
	//			}
	//		} else { // is not Enter, Backspace or control character, add to buf
	//			// todo: only add to buf if is allowed char, eg abc space etc, eg arrow key, don't add it
	//			h.buf = append(h.buf, pkt.Data...)
	//			_, err := pty.Write(pkt.Data)
	//			if err != nil {
	//				return "", err
	//			}
	//		}
	//	}
	//}
	//
	//return "", nil
}

//var (
//	ArrowUp          = []byte{27, 91, 65}
//	ArrowDown        = []byte{27, 91, 66}
//	ArrowRight       = []byte{27, 91, 67}
//	ArrowLeft        = []byte{27, 91, 68}
//	CtrlR            = []byte{18}
//	Enter            = []byte{13}
//	Backspace        = []byte{127}
//	PasteStart       = []byte{27, 91, 50, 48, 48, 126}
//	PasteEnd         = []byte{27, 91, 50, 48, 49, 126}
//	SemiColon        = []byte{59}
//	Ampersand        = []byte{38}
//	Pipe             = []byte{124}
//	LeftParenthesis  = []byte{40}
//	RightParenthesis = []byte{41}
//)
