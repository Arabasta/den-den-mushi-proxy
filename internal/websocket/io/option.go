package io

import "regexp"

type Option func(*bridgeConfig)

type bridgeConfig struct {
	whitelist *regexp.Regexp
	blacklist *regexp.Regexp
}

// WithWhitelist keeps only frames/lines that match the pattern.
func WithWhitelist(rx string) Option {
	return func(c *bridgeConfig) { c.whitelist = regexp.MustCompile(rx) }
}

// WithBlacklist discards frames/lines that match the pattern.
func WithBlacklist(rx string) Option {
	return func(c *bridgeConfig) { c.blacklist = regexp.MustCompile(rx) }
}
