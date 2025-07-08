package filter

import "sync"

type WhitelistFilter struct {
	mu               sync.RWMutex
	filteredCommands map[string]struct{}
}

func (w *WhitelistFilter) IsValid(cmd string) (string, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	_, blocked := w.filteredCommands[cmd]
	return cmd, blocked
}

func (w *WhitelistFilter) UpdateCommands(newAllowed []string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.filteredCommands = make(map[string]struct{})
	for _, cmd := range newAllowed {
		w.filteredCommands[cmd] = struct{}{}
	}
}
