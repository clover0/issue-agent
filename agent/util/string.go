package util

import "strings"

// TruncateLines truncates lines by keeping the first `keepStart` lines and the last `keepEnd` lines,
// replacing the middle section with `placeholder`.
// If the total lines are less than or equal to keepStart + keepEnd, no truncation occurs.
func TruncateLines(text string, keepStart, keepEnd int, placeholder string) string {
	if keepStart < 0 || keepEnd < 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	totalLines := len(lines)

	if totalLines <= keepStart+keepEnd {
		return text
	}

	startLines := lines[:keepStart]
	endLines := lines[totalLines-keepEnd:]

	result := strings.Join(startLines, "\n")
	if keepStart > 0 {
		result += "\n"
	}
	result += placeholder
	if keepEnd > 0 {
		result += "\n" + strings.Join(endLines, "\n")
	}

	return result
}
