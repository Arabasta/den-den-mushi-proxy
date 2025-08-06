package line

import (
	"sync"
)

const (
	MaxBufferLength = 8192
	ErrorMaxBuffer  = "???Error Max Buffer"
)

type Editor struct {
	mu     sync.RWMutex
	Buffer []rune
	Cursor int
}

func (e *Editor) Insert(r rune) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.Buffer) >= MaxBufferLength {
		e.Buffer = []rune(ErrorMaxBuffer)
		e.Cursor = len(e.Buffer)
		return
	}

	if e.Cursor >= len(e.Buffer) {
		e.Buffer = append(e.Buffer, r)
	} else {
		e.Buffer = append(e.Buffer, 0)
		copy(e.Buffer[e.Cursor+1:], e.Buffer[e.Cursor:])
		e.Buffer[e.Cursor] = r
	}
	e.Cursor++
}

func (e *Editor) Backspace() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Cursor > 0 && len(e.Buffer) > 0 {
		copy(e.Buffer[e.Cursor-1:], e.Buffer[e.Cursor:])
		e.Buffer = e.Buffer[:len(e.Buffer)-1]
		e.Cursor--
	}
}

func (e *Editor) MoveLeft() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Cursor > 0 {
		e.Cursor--
	}
}

func (e *Editor) MoveRight() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Cursor < len(e.Buffer) {
		e.Cursor++
	}
}
func (e *Editor) MoveStart() {
	e.mu.Lock()
	e.Cursor = 0
	e.mu.Unlock()
}

func (e *Editor) MoveEnd() {
	e.mu.Lock()
	e.Cursor = len(e.Buffer)
	e.mu.Unlock()
}

func (e *Editor) String() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return string(e.Buffer)
}

func (e *Editor) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Buffer = e.Buffer[:0]
	e.Cursor = 0
}
