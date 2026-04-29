package ui

import (
	"strings"

	"github.com/Gu1llaum-3/sshm/internal/config"
)

// applySourceFileFilter returns only the hosts whose SourceFile equals
// selected. An empty selected string is treated as "no filter" and returns
// the slice unchanged.
func applySourceFileFilter(hosts []config.SSHHost, selected string) []config.SSHHost {
	if selected == "" {
		return hosts
	}
	out := make([]config.SSHHost, 0, len(hosts))
	for _, h := range hosts {
		if h.SourceFile == selected {
			out = append(out, h)
		}
	}
	return out
}

// filterByTag returns hosts that have the given tag in either Tags or InheritedTags.
// The comparison is case-insensitive.
func filterByTag(hosts []config.SSHHost, tag string) []config.SSHHost {
	if tag == "" {
		return hosts
	}

	tag = strings.ToLower(tag)
	var result []config.SSHHost

	for _, host := range hosts {
		// Use AllTags() to check both own and inherited tags
		for _, t := range host.AllTags() {
			if strings.ToLower(t) == tag {
				result = append(result, host)
				break
			}
		}
	}

	return result
}
