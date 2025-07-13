package filter

import (
	"sync"
)

const (
	MaxBufferLength = 8192
	ErrorMaxBuffer  = "???Error Max Buffer"
)

type LineEditor struct {
	mu     sync.RWMutex
	Buffer []rune
	Cursor int
}

func (e *LineEditor) Insert(r rune) {
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
		e.Buffer = append(e.Buffer[:e.Cursor], append([]rune{r}, e.Buffer[e.Cursor:]...)...)
	}
	e.Cursor++
}

func (e *LineEditor) Backspace() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Cursor > 0 && len(e.Buffer) > 0 {
		e.Buffer = append(e.Buffer[:e.Cursor-1], e.Buffer[e.Cursor:]...)
		e.Cursor--
	}
}

func (e *LineEditor) MoveLeft() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Cursor > 0 {
		e.Cursor--
	}
}

func (e *LineEditor) MoveRight() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Cursor < len(e.Buffer) {
		e.Cursor++
	}
}

func (e *LineEditor) String() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return string(e.Buffer)
}

func (e *LineEditor) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Buffer = nil
	e.Cursor = 0
}
