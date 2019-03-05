package model

import (
	"regexp"
	"strings"
)

var (
	cleanNameRegex        = regexp.MustCompile(`[^#A-Za-z_]+`)
	cleanDisplayNameRegex = regexp.MustCompile(`[^0-9A-Za-z_]+`)
)

// StringCleanDisplayName sanitizes NPCs
func StringCleanDisplayName(in string) (out string) {
	out = strings.Replace(in, " ", "_", -1)
	out = strings.Replace(out, "#", "", -1)
	out = strings.TrimSpace(cleanDisplayNameRegex.ReplaceAllString(out, ""))
	out = strings.Replace(out, "_", " ", -1)
	return
}

// StringCleanName sanitizes NPCs
func StringCleanName(in string) (out string) {
	out = strings.Replace(in, " ", "_", -1)
	out = strings.TrimSpace(cleanNameRegex.ReplaceAllString(out, ""))
	return
}
