package formatter

// this file is for the automatic semicolon system in KrabbaScript
// its job is to automatically insert semicolons at the end of lines so you can omit semicolons in your code.

// there are two ways the code in this file can be reached:
// 1. the format command (`format`) — seen in cmd/format.go
// 2. the format flag in the build command (--format) — seen in cmd/build.go

// - Reece Turner (reeceturner1358@gmail.com)

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"regexp"
)

// Nothing really special, but it returns an array of strings and not a stupid array of bytes.
func read(path string) ([]string, error) {
	// read file
	lines, readError := os.ReadFile(path)

	if readError != nil {
		return nil, readError
	}

	// convert []byte to []string
	stringLines := strings.Split(strings.TrimSpace(string(lines)), "\n")

	return stringLines, nil
}


// for use in the insertSemicolons function
// an extremely asinine and convoluted way to remove trailing comments from a code.
var trailingComment = regexp.MustCompile(`^((?:[^"#']*|"[^"]*"|'[^']*')*)\s*(?:#.*)?$`)
func stripTrailingComment(line string) string {
	match := trailingComment.FindStringSubmatch(line)
	// no match was found — why the vittu doesn't it return false?
	if match == nil {
		return line // didn't match regex for a reason such as unbalanced quotes, etc. — leave the line alone
	}
	return strings.TrimRight(match[1], " \t")
}

// receives lines without semicolons;
// returns lines with semicolons
// for example, this turns
// var foo: I32 = 0
// into
// var foo: I32 = 0;
func insertSemicolons(lines []string) []string {
	// the lines with semicolons we'll be returning.
	var newLines []string

	symbolsToIgnore := []string{"{", "}", "(", "[", ",", ":", "&&", "||", "+", "-", "*", "/", "%", ";"}

	outer:
		for _, line := range lines {
			// keep the leading tabspace, only strip trailing whitespace
			line = strings.TrimRight(line, " \t\r")

			// ignore empty lines
			if strings.TrimSpace(line) == "" {
				newLines = append(newLines, line)
				continue
			}

			// string has a comment — this is mainly for inline comments
			if strings.Contains(line, "#") {
				// remove comment
				// this is done to avoid the case where a line like this:
				// print("I love Finland") # I really do!
				// being formatted like:
				// print("I love Finland") # I really do!;

				// in the future we'll make an effort to preserve the comment but this is fine for now.
				// We're making the mother of all omelettes here, Jack. Can't fret over every egg.
				line = stripTrailingComment(line)
			}

			// if the line ends with any of these symbols, 
			// then we don't need to add a semicolon.
			// I'll admit that guy was right — this is much cleaner than a thousand || conditions.
			for _, symbol := range symbolsToIgnore {
				if strings.HasSuffix(line, symbol) {
					newLines = append(newLines, line)
					continue outer
				}
			}

			// otherwise, append semicolon
			newLines = append(newLines, line+";")
		}

	return newLines
}

// format a file — read the file into lines, add semicolons to the lines, and overwrite the file with the new formatted lines.
func FormatFile(file string) error {
	lines, err := read(file)

	if err != nil {
		return err
	}

	formattedLines := insertSemicolons(lines)

	// replace old lines with formatted lines
	// the line noise is because I have to convert []strings into []bytes
	writeError := os.WriteFile(file, []byte(strings.Join(formattedLines, "\n")), 0644)

	if writeError != nil {
		return writeError
	}

	return nil
}

// batch format an entire directory
// once again, I use the fancy word format, but I really just mean insert semicolons.
// recursively walks a directory looking for krabbascript (.ks) files and formats each file
func FormatDirectory(dir string) error {
	// extract only the KrabbaScript (.ks) files
	var filePaths []string

	// I just copied this from https://zetcode.com/golang/find-file/
	// Don't hate the player, hate the game.
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		// only append if path is a file, and is a KrabbaScript file
		if !d.IsDir() && filepath.Ext(path) == ".ks" {
			filePaths = append(filePaths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// now format each KrabbaScript file
	for _, file := range filePaths {
		err := FormatFile(file)

		if err != nil {
			return err
		}
	}

	// no error
	return nil
}
