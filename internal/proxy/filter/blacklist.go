package filter

import "sync"

type BlacklistFilter struct {
	mu               sync.RWMutex
	filteredCommands map[string]struct{}
}

func (b *BlacklistFilter) IsValid(cmd string) (string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, blocked := b.filteredCommands[cmd]
	return cmd, !blocked
}

func (b *BlacklistFilter) UpdateCommands(newBlocked []string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.filteredCommands = make(map[string]struct{})
	for _, cmd := range newBlocked {
		b.filteredCommands[cmd] = struct{}{}
	}
}
