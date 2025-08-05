package util

import (
	"regexp"
)

// WhitelistEntry represents a host with allowed paths.
type WhitelistEntry struct {
	Host  string
	Paths []string //If empty, all paths are allowed for that host
	Exts  []string //If empty, all extensions are allowed
}

// Checks if the given host and path are whitelisted.
func IsWhitelisted(host, path string, whitelists ...WhitelistEntry) bool {
	for _, entry := range whitelists {
		//If the host matches, then allow further processing
		hostMatched, err := regexp.MatchString(entry.Host, host)
		if err != nil || !hostMatched {
			continue
		}

		//Allow all paths if the array is empty (after extension regex check)
		if len(entry.Paths) == 0 {
			if extAllowed(path, entry.Exts) {
				return true
			}
			continue
		}

		//If the current path matches any current regexp, then allow it (extension regex check)
		for _, pathPattern := range entry.Paths {
			matched, err := regexp.MatchString(pathPattern, path)
			if err == nil && matched {
				if extAllowed(path, entry.Exts) {
					return true
				}
			}
		}
	}

	//Options exhausted; deny by default
	return false
}

// Checks if the given extension is allowed (using regexp patterns).
func extAllowed(path string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}
	for _, extPattern := range exts {
		matched, err := regexp.MatchString(extPattern, path)
		if err == nil && matched {
			return true
		}
	}
	return false
}
