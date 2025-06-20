package io

import "bytes"

// applyFilters applies whitelist/blacklist rules
func applyFilters(b []byte, bCfg *bridgeConfig) []byte {
	if bCfg.blacklist != nil && bCfg.blacklist.Match(b) {
		return nil
	}
	if bCfg.whitelist != nil && !bCfg.whitelist.Match(b) {
		return nil
	}
	// return copy to avoid holding big buffer slice
	return bytes.Clone(b)
}
