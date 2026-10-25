package formatter

// this file is for the automatic semicolon system in KrabbaScript
// it's job is to automatically insert semicolons at the end of lines so you can omit semicolons in your code.

// - Reece Turner (reeceturner1358@gmail.com)


import (
    "strings"
)

// receives lines without semicolons——
// returns lines with semicolons
// for example, this turns 
// var foo: I32 = 0 
// into
// var foo: I32 = 0;
// This should ideally be called before the lexer to standardize input.
func InsertSemicolons(lines []string) []string {
	// the lines with semicolons we'll be returning.
	var newLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// ignore empty lines
		if line == "" {
			continue
		}

		// append lines without changes if they already end with a semicolon or don't need one
		if strings.HasSuffix(line, ";") || 
		strings.HasSuffix(line, "{") || 
		strings.HasSuffix(line, "}") ||
		strings.HasSuffix(line, "(") ||
		strings.HasSuffix(line, "[") ||
		strings.HasSuffix(line, ",") ||
		strings.HasSuffix(line, ":") ||
		strings.HasSuffix(line, "&&") ||
		strings.HasSuffix(line, "||") ||
		strings.HasSuffix(line, "+") ||
		strings.HasSuffix(line, "-") ||
		strings.HasSuffix(line, "*") ||
		strings.HasSuffix(line, "/") ||
		strings.HasSuffix(line, "%") ||
		strings.HasPrefix(line, "#") {
			newLines = append(newLines, line)
		} else {
			// append semicolon
			newLine := line + ";"
			newLines = append(newLines, newLine)
		}
	}

	return newLines
}


