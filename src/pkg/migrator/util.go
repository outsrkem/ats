package migrator

import (
	"errors"
	"strings"
)

func compareVersions(a, b string) int {
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")

	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	maxLen := max(len(aParts), len(bParts))
	for i := 0; i < maxLen; i++ {
		aVal := 0
		if i < len(aParts) {
			aVal = atoi(aParts[i])
		}

		bVal := 0
		if i < len(bParts) {
			bVal = atoi(bParts[i])
		}

		if aVal < bVal {
			return -1
		} else if aVal > bVal {
			return 1
		}
	}
	return 0
}

func atoi(s string) int {
	res := 0
	for _, c := range s {
		res = res*10 + int(c-'0')
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Split SQL statements and strictly enforce Linux-style line breaks (\n)
// Windows-style CRLF (\r\n) are not allowed and will cause an error
func splitSQL(sql string) ([]string, error) {
	// Check for Windows-style line breaks
	if strings.Contains(sql, "\r\n") {
		return nil, errors.New("SQL script contains Windows-style line breaks (\\r\\n). Only Linux-style line breaks (\\n) are allowed")
	}

	// Check for standalone carriage returns (likely invalid line breaks)
	if strings.Contains(sql, "\r") {
		return nil, errors.New("SQL script contains invalid \\r characters. Only Linux-style line breaks (\\n) are permitted")
	}

	// Split by semicolon
	parts := strings.Split(sql, ";")
	var statements []string
	for _, part := range parts {
		// Trim surrounding whitespace and empty lines
		stmt := strings.TrimSpace(part)
		if stmt != "" {
			// Append semicolon to maintain statement integrity
			statements = append(statements, stmt+";")
		}
	}
	return statements, nil
}
