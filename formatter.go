package formatter

// this file is for the automatic semicolon system in KrabbaScript
// it's job is to automatically insert semicolons at the end of lines so you can omit semicolons in your code.

// - Reece Turner (reeceturner1358@gmail.com)


import (
    "strings"
	"regexp"
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
		if strings.HasSuffix(line, ";") || strings.HasSuffix(line, "{") || strings.HasSuffix(line, "}") {
			newLines = append(newLines, line)
		} else {
			// append semicolon
			newLine := line + ";"
			newLines = append(newLines, newLine)
		}
	}

	return newLines
}

// receives lines without types and returns lines with types
// this allows you to create variables without specifying a type, granted you provide a value at creation.
// Therefore, creating a variable like:
// val foo
// will not work — you have to either specify a type or provide a value.
// Here is an example of what will work:
// val year = 1704
// This will become:
// val year: I32 = 1704
// this should also be called before running the lexer
func InferTypes(lines []string) []string {
	// the lines with types we'll be returning.
	var newLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// ignore empty lines
		if line == "" {
			continue
		}

		// if the line already has a type then skip it
		if strings.Contains(line, "Str") || strings.Contains(line, "I32") || strings.Contains(line, "Bool") || strings.Contains(line, "F32") {
			newLines = append(newLines, line)
			continue
		}

		// the regex patterns will be recompiled for every function call——
		// but it's not that bad since this should only be called once per file

		// value is a number (any)
		intRegex := regexp.MustCompile(`val (\w+) = (-?\d+)`)

		matchedIntRegex := intRegex.MatchString(line)
		if matchedIntRegex {
			// the type is inferred as an I32 — it's unspecific so don't use type inferrence with your F64s and whatnot.
			match := intRegex.FindStringSubmatch(line)

			// sandwich the type between the name and value
			newLine := "val " + match[1] + ": I32 = " + match[2]
			newLines = append(newLines, newLine) // DRY this up
		}


		// value is a string
		strRegex := regexp.MustCompile(`val (\w+) = "([^"]*)"`)

		matchedStrRegex := strRegex.MatchString(line)
		if matchedStrRegex {
			match := strRegex.FindStringSubmatch(line)

			// sandwich the type between the name and value
			newLine := "val " + match[1] + ": Str = \"" + match[2] + "\""
			newLines = append(newLines, newLine)
		}


		// value is a boolean
		boolRegex := regexp.MustCompile(`val (\w+) = (true|false)`)

		matchedBoolRegex := boolRegex.MatchString(line)
		if matchedBoolRegex {
			match := boolRegex.FindStringSubmatch(line)
			
			// sandwich the type between the name and value
			newLine := "val " + match[1] + ": Bool = " + match[2]
			newLines = append(newLines, newLine)
		}
	}

	return newLines
}
